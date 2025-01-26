// Package handlers implements HTTP handlers for service catalog CRUD operations.
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serviceCatalog/internal/constants"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

// DeleteService handles DELETE /services/:id endpoint.
//
// Performs soft deletion of a service and its associated versions.
// Uses GORM's soft delete mechanism to maintain data history.
//
// URL Parameters:
//   - id (uint): Service ID to delete
//
// Returns:
//
//	204: Service successfully deleted
//	400: Invalid service ID
//	500: Deletion failed
//
// Notes:
//   - Service is soft-deleted by default
//   - Associated versions are also soft-deleted via foreign key
//
// Example:
//
//	DELETE /services/1
func (h *Handler) DeleteService(c *gin.Context) {
	serviceID, validationErr := validation.ValidateServiceID(c)
	if validationErr != nil {
		c.JSON(validationErr.Status, validationErr)
		return
	}

	result := h.db.Delete(&models.Service{}, serviceID)
	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, &constants.ServiceError{
			Status:  constants.StatusInternalServerError,
			Message: constants.ErrServiceDeleteFailed,
			Details: result.Error.Error(),
		})

		return
	}

	c.Status(http.StatusNoContent)
}
