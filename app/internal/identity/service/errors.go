package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone number already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrAvatarNotFound     = errors.New("avatar not found")
)
