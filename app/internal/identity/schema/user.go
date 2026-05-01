package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

// User response DTO
type User struct {
	ID          int    `json:"id" example:"1"`
	Email       string `json:"email" example:"user@example.com"`
	FirstName   string `json:"first_name,omitempty" example:"John"`
	LastName    string `json:"last_name,omitempty" example:"Doe"`
	MiddleName  string `json:"middle_name,omitempty" example:"M"`
	Avatar      string `json:"avatar,omitempty" example:"https://s3.amazonaws.com/avatar.jpg"`
	IsSuperuser bool   `json:"is_superuser" example:"false"`
	IsActive    bool   `json:"is_active" example:"true"`
}

// AuthUser authenticated user DTO
type AuthUser struct {
	ID          int    `json:"id" example:"1"`
	Email       string `json:"email" example:"user@example.com"`
	FirstName   string `json:"first_name,omitempty" example:"John"`
	LastName    string `json:"last_name,omitempty" example:"Doe"`
	MiddleName  string `json:"middle_name,omitempty" example:"M"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+77001234567"`
	DateOfBirth string `json:"date_of_birth,omitempty" example:"1992-09-12"`
	Avatar      string `json:"avatar,omitempty"`
}

// UserCreateRequest registration request
// @Description Request body for user registration
type UserCreateRequest struct {
	Email       string `json:"email" example:"user@example.com" format:"email" maxLength:"255" validate:"required"`
	Password    string `json:"password" example:"secret123" minLength:"6" maxLength:"255" validate:"required"`
	FirstName   string `json:"first_name,omitempty" example:"John" maxLength:"255"`
	LastName    string `json:"last_name,omitempty" example:"Doe" maxLength:"255"`
	MiddleName  string `json:"middle_name,omitempty" example:"M" maxLength:"255"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+77001234567" maxLength:"128"`
	DateOfBirth string `json:"date_of_birth,omitempty" example:"1990-01-15" format:"date"`
}

func (r UserCreateRequest) Validate() error {
	v := shared.NewValidator()
	if r.Email == "" && r.PhoneNumber == "" {
		v.Append("email or phone_number is required")
	}
	if r.Email != "" {
		v.Email("email", r.Email)
		v.MaxLen("email", r.Email, 255)
	}
	v.Required("password", r.Password)
	v.MinLen("password", r.Password, 6)
	v.MaxLen("password", r.Password, 255)
	v.MaxLen("first_name", r.FirstName, 255)
	v.MaxLen("last_name", r.LastName, 255)
	v.MaxLen("middle_name", r.MiddleName, 255)
	v.MaxLen("phone_number", r.PhoneNumber, 128)
	v.Date("date_of_birth", r.DateOfBirth)
	return v.Result()
}

// UserUpdateRequest profile update request
// @Description Request body for updating user profile
type UserUpdateRequest struct {
	FirstName   *string `json:"first_name,omitempty" maxLength:"255"`
	LastName    *string `json:"last_name,omitempty" maxLength:"255"`
	MiddleName  *string `json:"middle_name,omitempty" maxLength:"255"`
	PhoneNumber *string `json:"phone_number,omitempty" maxLength:"128"`
	DateOfBirth *string `json:"date_of_birth,omitempty" format:"date"`
}

func (r UserUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("first_name", r.FirstName, 255)
	v.MaxLenPtr("last_name", r.LastName, 255)
	v.MaxLenPtr("middle_name", r.MiddleName, 255)
	v.MaxLenPtr("phone_number", r.PhoneNumber, 128)
	if r.DateOfBirth != nil {
		v.Date("date_of_birth", *r.DateOfBirth)
	}
	return v.Result()
}

// ChangePasswordRequest password change request
// @Description Request body for changing user password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"secret123" validate:"required"`
	NewPassword string `json:"new_password" example:"newsecret456" minLength:"6" maxLength:"255" validate:"required"`
}

func (r ChangePasswordRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("old_password", r.OldPassword)
	v.Required("new_password", r.NewPassword)
	v.MinLen("new_password", r.NewPassword, 6)
	v.MaxLen("new_password", r.NewPassword, 255)
	return v.Result()
}

// UserLoginRequest login request
// @Description Request body for user login
type UserLoginRequest struct {
	Email       string `json:"email,omitempty" example:"user@example.com" format:"email"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+77001234567"`
	Password    string `json:"password" example:"secret123" validate:"required"`
}

func (r UserLoginRequest) Validate() error {
	v := shared.NewValidator()
	if r.Email == "" && r.PhoneNumber == "" {
		v.Append("email or phone_number is required")
	}
	if r.Email != "" {
		v.Email("email", r.Email)
	}
	v.Required("password", r.Password)
	return v.Result()
}

// AuthResponse auth success response
type AuthResponse struct {
	User         AuthUser `json:"user"`
	AccessToken  string   `json:"access_token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string   `json:"refresh_token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// RefreshRequest refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
