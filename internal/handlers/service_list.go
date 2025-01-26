package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"serviceCatalog/internal/models"
	"time"
)

type serviceWithVersion struct {
	models.Service
	VersionCount int64 `gorm:"column:version_count"`
}

func (h *Handler) ListServices(c *gin.Context) {
	var params QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build base query
	query := h.db.Model(&models.Service{})

	services, totalCount := h.fetchListServices(c, params, query)

	// Convert to response format without additional DB calls
	serviceResponses := make([]models.ServiceResponse, 0, len(services))
	for _, service := range services {
		serviceResponses = append(serviceResponses, models.ServiceResponse{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Versions:    int(service.VersionCount),
		})
	}

	response := ListServicesResponse{
		Services:    serviceResponses,
		TotalCount:  totalCount,
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) fetchListServices(c *gin.Context, params QueryParams, query *gorm.DB) ([]serviceWithVersion, int64) {
	// Use gin.Context's timeout or default to 5 seconds
	timeoutDuration := 5 * time.Second
	if deadline, ok := c.Request.Context().Deadline(); ok {
		timeoutDuration = time.Until(deadline)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), timeoutDuration)
	defer cancel()

	query = query.WithContext(ctx)

	// Apply search filter if provided
	if params.Search != "" {
		query = query.Where("services.name ILIKE ? OR services.description ILIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Apply sorting
	var sortColumn string
	switch params.SortBy {
	case "id":
		sortColumn = "services.id"
	case "name":
		sortColumn = "services.name"
	case "description":
		sortColumn = "services.description"
	default:
		sortColumn = "services.id"
	}

	direction := "asc"
	if params.SortDir == "desc" {
		direction = "desc"
	}

	// Calculate total count with a subquery
	var totalCount int64
	countSubQuery := query.Session(&gorm.Session{}).Select("COUNT(*)")
	if err := countSubQuery.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count services"})
		return nil, 0
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize

	var services []serviceWithVersion

	// Build and execute final query with versions count in a single call
	result := query.
		Select("services.*, COALESCE(COUNT(versions.id), 0) as version_count").
		Joins("LEFT JOIN versions ON versions.service_id = services.id").
		Group("services.id").
		Order(sortColumn + " " + direction).
		Offset(offset).
		Limit(params.PageSize).
		Find(&services)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return nil, 0
	}

	return services, totalCount
}
