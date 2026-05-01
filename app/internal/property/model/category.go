package model

import "time"

type Category struct {
	ID        int       `json:"id"`
	NameKz    string    `json:"name_kz,omitempty"`
	NameRu    string    `json:"name_ru,omitempty"`
	NameEn    string    `json:"name_en,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	Icon      *string   `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HouseCategory struct {
	ID         int `json:"id"`
	HouseID    int `json:"house_id"`
	CategoryID int `json:"category_id"`
}
