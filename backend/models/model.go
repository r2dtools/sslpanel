package models

import "time"

// Model contains common fields for all models
type Model struct {
	ID        uint      `gorm:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
