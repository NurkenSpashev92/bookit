package models

import "time"

type User struct {
	ID          int        `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name,omitempty"`
	LastName    string     `json:"last_name,omitempty"`
	MiddleName  string     `json:"middle_name,omitempty"`
	Password    string     `json:"password,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Avatar      string     `json:"avatar,omitempty"`
	IsSuperuser bool       `json:"is_superuser"`
	IsActive    bool       `json:"is_active"`
	DateJoined  time.Time  `json:"date_joined"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
