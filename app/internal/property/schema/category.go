package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

type Category struct {
	Id       int    `json:"id" example:"1"`
	NameKz   string `json:"name_kz" example:"Апартаменты"`
	NameRu   string `json:"name_ru" example:"Апартаменты"`
	NameEn   string `json:"name_en" example:"Apartments"`
	IsActive bool   `json:"is_active" example:"true"`
}

// CategoryCreateRequest create category
// @Description Request body for creating a category
type CategoryCreateRequest struct {
	NameKz   string `json:"name_kz" maxLength:"255" example:"Апартаменты" validate:"required"`
	NameRu   string `json:"name_ru" maxLength:"255" example:"Апартаменты" validate:"required"`
	NameEn   string `json:"name_en" maxLength:"255" example:"Apartments" validate:"required"`
	IsActive *bool  `json:"is_active" example:"true"`
}

func (r CategoryCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("name_kz", r.NameKz)
	v.MaxLen("name_kz", r.NameKz, 255)
	v.Required("name_ru", r.NameRu)
	v.MaxLen("name_ru", r.NameRu, 255)
	v.Required("name_en", r.NameEn)
	v.MaxLen("name_en", r.NameEn, 255)
	return v.Result()
}

// CategoryUpdateRequest partial update category
// @Description Request body for updating a category (all fields optional)
type CategoryUpdateRequest struct {
	NameKz   *string `json:"name_kz,omitempty" maxLength:"255" example:"Апартаменты"`
	NameRu   *string `json:"name_ru,omitempty" maxLength:"255" example:"Апартаменты"`
	NameEn   *string `json:"name_en,omitempty" maxLength:"255" example:"Apartments"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

func (r CategoryUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("name_kz", r.NameKz, 255)
	v.MaxLenPtr("name_ru", r.NameRu, 255)
	v.MaxLenPtr("name_en", r.NameEn, 255)
	return v.Result()
}

type CategoryPaginate struct {
	Id       int    `json:"id" example:"1"`
	NameKz   string `json:"name_kz" example:"Апартаменты"`
	NameRu   string `json:"name_ru" example:"Апартаменты"`
	NameEn   string `json:"name_en" example:"Apartments"`
	IsActive bool   `json:"is_active" example:"true"`
}
