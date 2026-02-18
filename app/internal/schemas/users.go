package schemas

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	MiddleName  string `json:"middle_name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	IsSuperuser bool   `json:"is_superuser"`
	IsActive    bool   `json:"is_active"`
}

type AuthUser struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
}

type UserCreateRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	MiddleName  string `json:"middle_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  AuthUser `json:"user"`
	Token string   `json:"token"`
}
