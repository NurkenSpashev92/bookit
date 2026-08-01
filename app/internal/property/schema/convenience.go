package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

type Convenience struct {
	Id       int    `json:"id" example:"1"`
	Name     string `json:"name" example:"Wi-Fi"`
	Slug     string `json:"slug" example:"wi-fi"`
	IsActive bool   `json:"is_active" example:"true"`
}

// ConvenienceCreateRequest create convenience
// @Description Request body for creating a convenience
type ConvenienceCreateRequest struct {
	Name     string `json:"name" maxLength:"255" example:"Wi-Fi" validate:"required"`
	Slug     string `json:"slug,omitempty" maxLength:"255" example:"wi-fi"`
	IsActive *bool  `json:"is_active" example:"true"`
}

func (r ConvenienceCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("name", r.Name)
	v.MaxLen("name", r.Name, 255)
	v.MaxLen("slug", r.Slug, 255)
	return v.Result()
}

// ConvenienceUpdateRequest partial update convenience
// @Description Request body for updating a convenience (all fields optional)
type ConvenienceUpdateRequest struct {
	Name     *string `json:"name,omitempty" maxLength:"255" example:"Wi-Fi"`
	Slug     *string `json:"slug,omitempty" maxLength:"255" example:"wi-fi"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

func (r ConvenienceUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("name", r.Name, 255)
	v.MaxLenPtr("slug", r.Slug, 255)
	return v.Result()
}

type ConveniencePaginate struct {
	Id       int    `json:"id" example:"1"`
	Name     string `json:"name" example:"Wi-Fi"`
	Slug     string `json:"slug" example:"wi-fi"`
	IsActive bool   `json:"is_active" example:"true"`
}
