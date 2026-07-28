package service

import (
	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

var (
	ErrEmailAlreadyExists = model.ErrEmailExists
	ErrPhoneAlreadyExists = model.ErrPhoneExists
	ErrUserNotFound       = model.ErrUserNotFound

	ErrInvalidCredentials = shared.Unauthorized("invalid email or password")
	ErrWrongPassword      = shared.Unauthorized("old password is incorrect")
	ErrInvalidToken       = shared.Unauthorized("invalid or expired token")
	ErrAccountDisabled    = shared.Unauthorized("account is disabled")
	ErrAvatarNotFound     = shared.NotFound("avatar not found")
)
