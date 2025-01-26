// Package handlers implements HTTP endpoints for the service catalog API.
// It provides CRUD operations for services and their versions with support
// for pagination, searching, sorting and soft deletion.
package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"serviceCatalog/internal/constants"
	"serviceCatalog/internal/models"
	"time"
)

// ServiceWithVersion represents a service entity with its version count.
// This structure is used internally to optimize database queries by avoiding
// N+1 problems when fetching version counts.
//
// Fields:
//   - Service: Base service model containing core service information
//   - VersionCount: Total number of versions associated with this service,
//     computed via JOIN in the same query
type serviceWithVersion struct {
	models.Service
	VersionCount int64 `gorm:"column:version_count"`
}

// ListServices returns a paginated list of services.
//
// This endpoint supports:
//   - Pagination via page and pageSize parameters
//   - Case-insensitive search in name and description fields
//   - Sorting by multiple columns with direction control
//   - Inclusion of soft-deleted records
//   - Efficient version counting via JOIN
//
// Query Parameters:
//   - page (int): Page number, starting from 1
//   - pageSize (int): Number of items per page (default: 10, max: 100)
//   - search (string): Search term for filtering by name or description
//   - sortBy (string): Column to sort by ("id", "name", "description")
//   - sortDir (string): Sort direction ("asc", "desc")
//   - showDeleted (bool): Whether to include soft-deleted records
//
// Returns:
//   200 OK: ListServicesResponse{
//     services: []ServiceResponse - List of services with version counts
//     totalCount: int - Total number of matching records
//     currentPage: int - Current page number
//     pageSize: int - Items per page
//   }
//   400 Bad Request: Invalid query parameters
//   500 Internal Server Error: Database or server errors
//
// Example Usage:
//   GET /services?page=1&pageSize=10&search=auth&sortBy=name&sortDir=asc

func (h *Handler) ListServices(c *gin.Context) {
	// Bind and validate query parameters
	var params QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, &constants.ServiceError{
			Status:  constants.StatusBadRequest,
			Message: constants.ErrBadRequest,
			Details: err.Error(),
		})

		return
	}

	// Initialize base query
	query := h.db.Model(&models.Service{})

	// Handle soft deletion filter
	// If showDeleted=true, include soft-deleted records
	if showDeleted := c.Query(constants.ShowDeleted); showDeleted == constants.True {
		query = query.Unscoped()
	}

	// Fetch services with optimized version counting
	services, totalCount := h.fetchListServices(c, params, query)

	// Transform database models to response DTOs
	// Pre-allocate slice capacity for better performance
	serviceResponses := make([]models.ServiceResponse, 0, len(services))
	for _, service := range services {
		serviceResponses = append(serviceResponses, models.ServiceResponse{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Versions:    int(service.VersionCount),
		})
	}

	// Construct final response with pagination metadata
	response := ListServicesResponse{
		Services:    serviceResponses,
		TotalCount:  totalCount,
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
	}

	c.JSON(http.StatusOK, response)
}

// fetchListServices retrieves and paginates services from the database.
//
// This function handles complex query building including:
//   - Timeout management via context
//   - Search filtering with ILIKE for PostgreSQL
//   - Dynamic column sorting with validation
//   - Efficient pagination with total count
//   - Version counting via LEFT JOIN
//
// Parameters:
//   - c *gin.Context: Request context for timeout and cancellation
//   - params QueryParams: Validated query parameters for filtering and pagination
//   - query *gorm.DB: Base query to build upon, may include initial filters
//
// Returns:
//   - []ServiceWithVersion: Slice of services with their version counts
//   - int64: Total count of matching records before pagination
//   - error: Any errors encountered during query execution
//
// Query Performance:
//   - Uses single query with JOIN for version counting
//   - Implements pagination before JOIN for better performance
//   - Uses indexed columns for sorting and filtering
//   - Handles NULL cases with COALESCE
func (h *Handler) fetchListServices(c *gin.Context, params QueryParams, query *gorm.DB) ([]serviceWithVersion, int64) {

	// Setup query timeout using context deadline or default 5s
	timeoutDuration := 5 * time.Second
	if deadline, ok := c.Request.Context().Deadline(); ok {
		timeoutDuration = time.Until(deadline)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), timeoutDuration)
	defer cancel()

	query = query.WithContext(ctx)

	// Apply case-insensitive search on name and description
	// ILIKE is PostgreSQL specific, provides better performance than LOWER()
	if params.Search != "" {
		query = query.Where("services.name ILIKE ? OR services.description ILIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Determine sort column with input validation
	var sortColumn string
	switch params.SortBy {
	case constants.DefaultSortField:
		sortColumn = "services.id"
	case constants.Name:
		sortColumn = "services.name"
	case constants.Description:
		sortColumn = "services.description"
	default:
		sortColumn = "services.id" // Fallback to ID sorting
	}

	// Apply sort direction
	direction := constants.DefaultSortOrder
	if params.SortDir == constants.DescSortOrder {
		direction = constants.DescSortOrder
	}

	// Get total count before pagination for metadata
	var totalCount int64
	countSubQuery := query.Session(&gorm.Session{}).Select("COUNT(*)")
	if err := countSubQuery.Count(&totalCount).Error; err != nil {

		c.JSON(http.StatusInternalServerError, &constants.ServiceError{
			Status:  constants.StatusInternalServerError,
			Message: constants.ErrServiceCountFailed,
			Details: err.Error(),
		})

		return nil, 0
	}

	// Calculate pagination offset
	offset := (params.Page - 1) * params.PageSize

	var services []serviceWithVersion

	// Build and execute final query with versions count in a single call
	// Execute final query combining:
	// - Base filters from the input query
	// - Version counting using LEFT JOIN
	// - Grouping to handle the aggregate
	// - Sorting and pagination
	result := query.
		Select("services.*, COALESCE(COUNT(versions.id), 0) as version_count").
		Joins("LEFT JOIN versions ON versions.service_id = services.id").
		Group("services.id").
		Order(sortColumn + " " + direction).
		Offset(offset).
		Limit(params.PageSize).
		Find(&services)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &constants.ServiceError{
			Status:  constants.StatusInternalServerError,
			Message: constants.ErrServicesFetchFailed,
			Details: result.Error.Error(),
		})

		return nil, 0
	}

	return services, totalCount
}
