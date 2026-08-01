package model

import "time"

type Type struct {
	ID        int       `json:"id"`
	NameKz    string    `json:"name_kz"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
