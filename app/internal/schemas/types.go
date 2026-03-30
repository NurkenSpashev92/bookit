package schemas

// TypeResponse type response DTO
type TypeResponse struct {
	ID       int     `json:"id" example:"1"`
	Name     string  `json:"name,omitempty" example:"Квартира" maxLength:"255"`
	Icon     *string `json:"icon,omitempty" example:"https://s3.amazonaws.com/icon.png"`
	IsActive bool    `json:"is_active" example:"true"`
}

// TypeCreateRequest create type (multipart/form-data)
// @Description Request body for creating a type
type TypeCreateRequest struct {
	Name     string `json:"name" form:"name" maxLength:"255" example:"Квартира" validate:"required"`
	IsActive *bool  `json:"is_active" form:"is_active" example:"true"`
}

func (r TypeCreateRequest) Validate() error {
	v := newValidator()
	v.required("name", r.Name)
	v.maxLen("name", r.Name, 255)
	return v.result()
}

// TypeUpdateRequest partial update type (multipart/form-data)
// @Description Request body for updating a type (all fields optional)
type TypeUpdateRequest struct {
	Name     *string `json:"name,omitempty" form:"name" maxLength:"255" example:"Квартира"`
	IsActive *bool   `json:"is_active,omitempty" form:"is_active" example:"true"`
}

func (r TypeUpdateRequest) Validate() error {
	v := newValidator()
	v.maxLenPtr("name", r.Name, 255)
	return v.result()
}
