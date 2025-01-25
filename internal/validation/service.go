package validation

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type ServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func ValidateServiceID(c *gin.Context) (uint64, *ServiceError) {
	id := c.Param("id")
	if id == "" {
		return 0, &ServiceError{
			Status:  400,
			Message: "service ID is required",
		}
	}

	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, &ServiceError{
			Status:  400,
			Message: "Invalid service ID: must be a positive integer",
			Details: err.Error(),
		}
	}
	return serviceID, nil
}
