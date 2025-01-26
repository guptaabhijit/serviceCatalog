// Package handlers implements HTTP handlers for service catalog CRUD operations.
package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"serviceCatalog/internal/constants"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

// GetService handles GET /services/:id endpoint.
//
// Retrieves a single service by ID with its version count.
// Supports including soft-deleted services via showDeleted parameter.
//
// URL Parameters:
//   - id (uint): Service ID
//
// Query Parameters:
//   - showDeleted (bool): Include soft-deleted service if true
//
// Returns:
//
//	200: ServiceResponse with service details and version count
//	400: Invalid service ID
//	404: Service not found
//	500: Database or server errors
//
// Example:
//
//	GET /services/1?showDeleted=true
func (h *Handler) GetService(c *gin.Context) {
	serviceID, validationErr := validation.ValidateServiceID(c)
	if validationErr != nil {
		c.JSON(validationErr.Status, validationErr)
		return
	}

	var service models.Service
	var result *gorm.DB

	// Add deleted filter if provided
	if showDeleted := c.Query(constants.ShowDeleted); showDeleted == constants.True {
		result = h.db.Unscoped().First(&service, serviceID)
	} else {
		result = h.db.First(&service, serviceID)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, &constants.ServiceError{
				Status:  constants.StatusNotFound,
				Message: constants.ErrServiceNotFound,
				Details: result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &constants.ServiceError{
			Status:  constants.StatusInternalServerError,
			Message: constants.ErrServiceFetchFailed,
			Details: result.Error.Error(),
		})

		return
	}

	// Get version count
	var versionCount int64
	h.db.Model(&models.Version{}).Where("service_id = ?", service.ID).Count(&versionCount)

	response := service.ToResponse(int(versionCount))
	c.JSON(http.StatusOK, response)
}
