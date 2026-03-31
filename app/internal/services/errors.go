package services

import "errors"

var (
	ErrSlugExists        = errors.New("slug already exists")
	ErrMaxImagesExceeded = errors.New("maximum 15 images allowed")
	ErrImageNotFound     = errors.New("image not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone number already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrAvatarNotFound     = errors.New("avatar not found")
	ErrBookingOverlap     = errors.New("these dates are already booked")
)
