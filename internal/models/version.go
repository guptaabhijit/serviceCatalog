package models

import (
	"time"
)

type Version struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ServiceID uint      `json:"service_id"`
	Number    string    `json:"number" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
