// Package handlers implements HTTP handlers for service catalog CRUD operations.
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serviceCatalog/internal/constants"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

// GetServiceVersions handles GET /services/:id/versions endpoint.
//
// Retrieves all versions associated with a service.
// Returns versions in chronological order.
//
// URL Parameters:
//   - id (uint): Service ID
//
// Returns:
//
//	200: []Version - List of service versions
//	400: Invalid service ID
//	500: Database error
//
// Example:
//
//	GET /services/1/versions
func (h *Handler) GetServiceVersions(c *gin.Context) {
	serviceID, validationErr := validation.ValidateServiceID(c)
	if validationErr != nil {
		c.JSON(validationErr.Status, validationErr)
		return
	}

	var versions []models.Version
	result := h.db.Where("service_id = ?", serviceID).Find(&versions)

	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, &constants.ServiceError{
			Status:  constants.StatusInternalServerError,
			Message: constants.ErrVersionFetchFailed,
			Details: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, versions)
}
