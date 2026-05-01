package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

// City response DTO
type City struct {
	ID          int     `json:"id" example:"1"`
	NameKZ      string  `json:"name_kz" example:"Астана"`
	NameEN      string  `json:"name_en" example:"Astana"`
	NameRU      string  `json:"name_ru" example:"Астана"`
	PostallCode string  `json:"postall_code,omitempty" example:"010000"`
	Country     Country `json:"country"`
}

// CityCreateRequest create city request
// @Description Request body for creating a city
type CityCreateRequest struct {
	NameKZ      string `json:"name_kz" form:"name_kz" maxLength:"255" example:"Астана" validate:"required"`
	NameEN      string `json:"name_en" form:"name_en" maxLength:"255" example:"Astana" validate:"required"`
	NameRU      string `json:"name_ru" form:"name_ru" maxLength:"255" example:"Астана" validate:"required"`
	PostallCode string `json:"postall_code" form:"postall_code" maxLength:"20" example:"010000"`
	CountryID   int    `json:"country_id" form:"country_id" example:"1" validate:"required"`
}

func (r CityCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("name_kz", r.NameKZ)
	v.MaxLen("name_kz", r.NameKZ, 255)
	v.Required("name_en", r.NameEN)
	v.MaxLen("name_en", r.NameEN, 255)
	v.Required("name_ru", r.NameRU)
	v.MaxLen("name_ru", r.NameRU, 255)
	v.MaxLen("postall_code", r.PostallCode, 20)
	v.RequiredInt("country_id", r.CountryID)
	return v.Result()
}

// CityUpdateRequest partial update city request
// @Description Request body for updating a city (all fields optional)
type CityUpdateRequest struct {
	NameKZ      *string `json:"name_kz,omitempty" form:"name_kz" maxLength:"255" example:"Астана"`
	NameEN      *string `json:"name_en,omitempty" form:"name_en" maxLength:"255" example:"Astana"`
	NameRU      *string `json:"name_ru,omitempty" form:"name_ru" maxLength:"255" example:"Астана"`
	PostallCode *string `json:"postall_code,omitempty" form:"postall_code" maxLength:"20"`
	CountryID   *int    `json:"country_id,omitempty" form:"country_id" example:"1"`
}

func (r CityUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("name_kz", r.NameKZ, 255)
	v.MaxLenPtr("name_en", r.NameEN, 255)
	v.MaxLenPtr("name_ru", r.NameRU, 255)
	v.MaxLenPtr("postall_code", r.PostallCode, 20)
	return v.Result()
}
