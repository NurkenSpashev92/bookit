package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrPhoneExists  = errors.New("phone number already exists")
)
