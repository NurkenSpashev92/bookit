package schemas

// Country response DTO
type Country struct {
	ID     int    `json:"id" example:"1"`
	NameKZ string `json:"name_kz" form:"name_kz" example:"Қазақстан"`
	NameEN string `json:"name_en" form:"name_en" example:"Kazakhstan"`
	NameRU string `json:"name_ru" form:"name_ru" example:"Казахстан"`
	Code   string `json:"code" form:"code" example:"KZ"`
}

// CountryCreateRequest create country request
// @Description Request body for creating a country
type CountryCreateRequest struct {
	NameKZ string `json:"name_kz" form:"name_kz" maxLength:"255" example:"Қазақстан" validate:"required"`
	NameEN string `json:"name_en" form:"name_en" maxLength:"255" example:"Kazakhstan" validate:"required"`
	NameRU string `json:"name_ru" form:"name_ru" maxLength:"255" example:"Казахстан" validate:"required"`
	Code   string `json:"code" form:"code" maxLength:"10" example:"KZ"`
}

func (r CountryCreateRequest) Validate() error {
	v := newValidator()
	v.required("name_kz", r.NameKZ)
	v.maxLen("name_kz", r.NameKZ, 255)
	v.required("name_en", r.NameEN)
	v.maxLen("name_en", r.NameEN, 255)
	v.required("name_ru", r.NameRU)
	v.maxLen("name_ru", r.NameRU, 255)
	v.maxLen("code", r.Code, 10)
	return v.result()
}

// CountryUpdateRequest partial update country request
// @Description Request body for updating a country (all fields optional)
type CountryUpdateRequest struct {
	NameKZ *string `json:"name_kz,omitempty" form:"name_kz" maxLength:"255" example:"Қазақстан"`
	NameEN *string `json:"name_en,omitempty" form:"name_en" maxLength:"255" example:"Kazakhstan"`
	NameRU *string `json:"name_ru,omitempty" form:"name_ru" maxLength:"255" example:"Казахстан"`
	Code   *string `json:"code,omitempty" form:"code" maxLength:"10" example:"KZ"`
}

func (r CountryUpdateRequest) Validate() error {
	v := newValidator()
	v.maxLenPtr("name_kz", r.NameKZ, 255)
	v.maxLenPtr("name_en", r.NameEN, 255)
	v.maxLenPtr("name_ru", r.NameRU, 255)
	v.maxLenPtr("code", r.Code, 10)
	return v.result()
}
