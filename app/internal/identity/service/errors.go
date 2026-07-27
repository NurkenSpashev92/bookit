package service

import (
	"errors"

	"github.com/nurkenspashev92/bookit/internal/identity/model"
)

var (
	ErrEmailAlreadyExists = model.ErrEmailExists
	ErrPhoneAlreadyExists = model.ErrPhoneExists

	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrAvatarNotFound     = errors.New("avatar not found")
)
