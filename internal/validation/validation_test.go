package validation

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateServiceID(t *testing.T) {
	tests := []struct {
		name          string
		paramID       string
		expectedID    uint64
		expectedError bool
	}{
		{
			name:          "Valid ID",
			paramID:       "123",
			expectedID:    123,
			expectedError: false,
		},
		{
			name:          "Invalid ID - String",
			paramID:       "abc",
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "Invalid ID - Negative",
			paramID:       "-1",
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(nil)
			c.Params = []gin.Param{{Key: "id", Value: tt.paramID}}

			id, err := ValidateServiceID(c)

			if tt.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
