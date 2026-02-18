package models

import "time"

type Convenience struct {
	ID        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	Icon      string    `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HouseConvenience struct {
	ID            int `json:"id"`
	HouseID       int `json:"house_id"`
	ConvenienceID int `json:"convenience_id"`
}
