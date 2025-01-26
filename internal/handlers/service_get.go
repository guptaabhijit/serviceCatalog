package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

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
