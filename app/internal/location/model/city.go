package model

import "time"

type City struct {
	ID          int       `json:"id"`
	NameRU      string    `json:"name_ru"`
	NameEN      string    `json:"name_en"`
	NameKZ      string    `json:"name_kz"`
	PostallCode string    `json:"postall_code,omitempty"`
	CountryID   int       `json:"country_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
