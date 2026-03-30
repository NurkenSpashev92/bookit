package models

import "time"

type Type struct {
	ID        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	Icon      *string   `json:"icon,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
