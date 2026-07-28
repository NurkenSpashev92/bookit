package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrUserNotFound = shared.NotFound("user not found")
	ErrEmailExists  = shared.Conflict("email already exists")
	ErrPhoneExists  = shared.Conflict("phone number already exists")
)
