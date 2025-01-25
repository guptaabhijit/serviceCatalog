package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
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

func (h *Handler) GetService(c *gin.Context) {
	serviceID, validationErr := validation.ValidateServiceID(c)
	if validationErr != nil {
		c.JSON(validationErr.Status, validationErr)
		return
	}

	var service models.Service
	result := h.db.First(&service, serviceID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service"})
		return
	}

	// Get version count
	var versionCount int64
	h.db.Model(&models.Version{}).Where("service_id = ?", service.ID).Count(&versionCount)

	response := service.ToResponse(int(versionCount))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetServiceVersions(c *gin.Context) {
	serviceID, validationErr := validation.ValidateServiceID(c)
	if validationErr != nil {
		c.JSON(validationErr.Status, validationErr)
		return
	}

	var versions []models.Version
	result := h.db.Where("service_id = ?", serviceID).Find(&versions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

type serviceWithVersion struct {
	models.Service
	VersionCount int64 `gorm:"column:version_count"`
}

func (h *Handler) fetchListServices(c *gin.Context, params QueryParams, query *gorm.DB) ([]serviceWithVersion, int64) {
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
