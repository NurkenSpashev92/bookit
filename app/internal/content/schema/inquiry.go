package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

// Inquiry response DTO
type Inquiry struct {
	ID          int    `json:"id" example:"1"`
	Email       string `json:"email" example:"guest@example.com" format:"email"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+77001234567"`
	Text        string `json:"text" example:"I have a question about booking"`
	IsApproved  bool   `json:"is_approved" example:"false"`
}

// InquiryCreateRequest create inquiry request
// @Description Request body for creating an inquiry
type InquiryCreateRequest struct {
	Email       string `json:"email" format:"email" maxLength:"255" example:"guest@example.com" validate:"required"`
	PhoneNumber string `json:"phone_number,omitempty" maxLength:"20" example:"+77001234567"`
	Text        string `json:"text" example:"I have a question about booking" validate:"required"`
	IsApproved  bool   `json:"is_approved" example:"false"`
}

func (r InquiryCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("email", r.Email)
	v.Email("email", r.Email)
	v.MaxLen("email", r.Email, 255)
	v.MaxLen("phone_number", r.PhoneNumber, 20)
	v.Required("text", r.Text)
	return v.Result()
}

// InquiryUpdateRequest partial update inquiry
// @Description Request body for updating an inquiry (all fields optional)
type InquiryUpdateRequest struct {
	Email       *string `json:"email,omitempty" format:"email" maxLength:"255" example:"guest@example.com"`
	PhoneNumber *string `json:"phone_number,omitempty" maxLength:"20"`
	Text        *string `json:"text,omitempty"`
	IsApproved  *bool   `json:"is_approved,omitempty" example:"true"`
}

func (r InquiryUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.EmailPtr("email", r.Email)
	v.MaxLenPtr("email", r.Email, 255)
	v.MaxLenPtr("phone_number", r.PhoneNumber, 20)
	return v.Result()
}
