package service

import (
	"errors"

	"github.com/nurkenspashev92/bookit/internal/property/model"
)

var (
	ErrSlugExists = model.ErrSlugExists

	ErrMaxImagesExceeded = errors.New("maximum 15 images allowed")
	ErrImageNotFound     = errors.New("image not found")
)
