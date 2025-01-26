package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serviceCatalog/internal/models"
	"serviceCatalog/internal/validation"
)

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
