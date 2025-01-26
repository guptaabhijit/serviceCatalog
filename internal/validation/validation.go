package validation

import (
	"github.com/gin-gonic/gin"
	"serviceCatalog/internal/constants"
)

type ServiceIDParam struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

func ValidateServiceID(c *gin.Context) (uint64, *constants.ServiceError) {
	var param ServiceIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		return 0, &constants.ServiceError{
			Status:  constants.StatusBadRequest,
			Message: constants.ErrInvalidServiceID,
			Details: err.Error(),
		}
	}
	return param.ID, nil
}
