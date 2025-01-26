package validation

import (
	"github.com/gin-gonic/gin"
)

type ServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ServiceIDParam struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

func ValidateServiceID(c *gin.Context) (uint64, *ServiceError) {
	var param ServiceIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		return 0, &ServiceError{
			Status:  400,
			Message: "Invalid service ID: must be a positive integer",
			Details: err.Error(),
		}
	}
	return param.ID, nil
}
