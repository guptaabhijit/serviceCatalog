package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceToResponse(t *testing.T) {
	service := Service{
		ID:          1,
		Name:        "Test Service",
		Description: "Test Description",
	}

	response := service.ToResponse(3)

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test Service", response.Name)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, 3, response.Versions)
}
