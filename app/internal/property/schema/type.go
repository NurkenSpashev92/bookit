package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

type TypeResponse struct {
	ID       int    `json:"id" example:"1"`
	NameKz   string `json:"name_kz" example:"Пәтер"`
	NameRu   string `json:"name_ru" example:"Квартира"`
	NameEn   string `json:"name_en" example:"Apartment"`
	IsActive bool   `json:"is_active" example:"true"`
}

// TypeCreateRequest create type
// @Description Request body for creating a type
type TypeCreateRequest struct {
	NameKz   string `json:"name_kz" maxLength:"255" example:"Пәтер" validate:"required"`
	NameRu   string `json:"name_ru" maxLength:"255" example:"Квартира" validate:"required"`
	NameEn   string `json:"name_en" maxLength:"255" example:"Apartment" validate:"required"`
	IsActive *bool  `json:"is_active" example:"true"`
}

func (r TypeCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("name_kz", r.NameKz)
	v.MaxLen("name_kz", r.NameKz, 255)
	v.Required("name_ru", r.NameRu)
	v.MaxLen("name_ru", r.NameRu, 255)
	v.Required("name_en", r.NameEn)
	v.MaxLen("name_en", r.NameEn, 255)
	return v.Result()
}

// TypeUpdateRequest partial update type
// @Description Request body for updating a type (all fields optional)
type TypeUpdateRequest struct {
	NameKz   *string `json:"name_kz,omitempty" maxLength:"255" example:"Пәтер"`
	NameRu   *string `json:"name_ru,omitempty" maxLength:"255" example:"Квартира"`
	NameEn   *string `json:"name_en,omitempty" maxLength:"255" example:"Apartment"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

func (r TypeUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("name_kz", r.NameKz, 255)
	v.MaxLenPtr("name_ru", r.NameRu, 255)
	v.MaxLenPtr("name_en", r.NameEn, 255)
	return v.Result()
}
