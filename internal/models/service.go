package models

import (
	"time"
)

type Service struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Versions    []Version `json:"versions,omitempty" gorm:"foreignKey:ServiceID"`
}

// ServiceResponse is the API response structure
type ServiceResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Versions    int    `json:"versions"`
}

// ToResponse converts the Service model to a ServiceResponse
func (s *Service) ToResponse(versionCount int) ServiceResponse {
	return ServiceResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Versions:    versionCount,
	}
}
