package models

import "time"

type Country struct {
	ID        int       `json:"id"`
	NameKZ    string    `json:"name_kz"`
	NameEN    string    `json:"name_en"`
	NameRU    string    `json:"name_ru"`
	Code      string    `json:"code,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
