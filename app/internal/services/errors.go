package services

import "errors"

var (
	ErrSlugExists        = errors.New("slug already exists")
	ErrMaxImagesExceeded = errors.New("maximum 15 images allowed")
	ErrImageNotFound     = errors.New("image not found")
)
