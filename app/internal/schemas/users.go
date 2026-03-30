package schemas

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
	ID         int    `json:"id" example:"1"`
	Email      string `json:"email" example:"user@example.com"`
	FirstName  string `json:"first_name,omitempty" example:"John"`
	LastName   string `json:"last_name,omitempty" example:"Doe"`
	MiddleName string `json:"middle_name,omitempty" example:"M"`
	Avatar     string `json:"avatar,omitempty"`
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
	v := newValidator()
	v.required("email", r.Email)
	v.email("email", r.Email)
	v.maxLen("email", r.Email, 255)
	v.required("password", r.Password)
	v.minLen("password", r.Password, 6)
	v.maxLen("password", r.Password, 255)
	v.maxLen("first_name", r.FirstName, 255)
	v.maxLen("last_name", r.LastName, 255)
	v.maxLen("middle_name", r.MiddleName, 255)
	v.maxLen("phone_number", r.PhoneNumber, 128)
	v.date("date_of_birth", r.DateOfBirth)
	return v.result()
}

// UserLoginRequest login request
// @Description Request body for user login
type UserLoginRequest struct {
	Email    string `json:"email" example:"user@example.com" format:"email" validate:"required"`
	Password string `json:"password" example:"secret123" validate:"required"`
}

func (r UserLoginRequest) Validate() error {
	v := newValidator()
	v.required("email", r.Email)
	v.email("email", r.Email)
	v.required("password", r.Password)
	return v.result()
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
