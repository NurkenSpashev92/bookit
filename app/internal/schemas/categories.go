package schemas

// Category response DTO
type Category struct {
	Id       int     `json:"id" example:"1"`
	NameKz   string  `json:"name_kz" example:"Апартаменты"`
	NameRu   string  `json:"name_ru" example:"Апартаменты"`
	NameEn   string  `json:"name_en" example:"Apartments"`
	Icon     *string `json:"icon" example:"https://s3.amazonaws.com/icon.png"`
	IsActive bool    `json:"is_active" example:"true"`
}

// CategoryCreateRequest create category (multipart/form-data)
// @Description Request body for creating a category
type CategoryCreateRequest struct {
	NameKz   string  `json:"name_kz" maxLength:"255" example:"Апартаменты" validate:"required"`
	NameRu   string  `json:"name_ru" maxLength:"255" example:"Апартаменты" validate:"required"`
	NameEn   string  `json:"name_en" maxLength:"255" example:"Apartments" validate:"required"`
	Icon     *string `json:"icon"`
	IsActive bool    `json:"is_active" example:"true"`
}

func (r CategoryCreateRequest) Validate() error {
	v := newValidator()
	v.required("name_kz", r.NameKz)
	v.maxLen("name_kz", r.NameKz, 255)
	v.required("name_ru", r.NameRu)
	v.maxLen("name_ru", r.NameRu, 255)
	v.required("name_en", r.NameEn)
	v.maxLen("name_en", r.NameEn, 255)
	return v.result()
}

// CategoryUpdateRequest partial update category (multipart/form-data)
// @Description Request body for updating a category (all fields optional)
type CategoryUpdateRequest struct {
	NameKz   *string `form:"name_kz" maxLength:"255" example:"Апартаменты"`
	NameRu   *string `form:"name_ru" maxLength:"255" example:"Апартаменты"`
	NameEn   *string `form:"name_en" maxLength:"255" example:"Apartments"`
	IsActive *bool   `form:"is_active" example:"true"`
}

func (r CategoryUpdateRequest) Validate() error {
	v := newValidator()
	v.maxLenPtr("name_kz", r.NameKz, 255)
	v.maxLenPtr("name_ru", r.NameRu, 255)
	v.maxLenPtr("name_en", r.NameEn, 255)
	return v.result()
}

// CategoryPaginate category list item
type CategoryPaginate struct {
	Id       int     `json:"id" example:"1"`
	NameKz   string  `json:"name_kz" example:"Апартаменты"`
	NameRu   string  `json:"name_ru" example:"Апартаменты"`
	NameEn   string  `json:"name_en" example:"Apartments"`
	Icon     *string `json:"icon"`
	IsActive bool    `json:"is_active" example:"true"`
}
