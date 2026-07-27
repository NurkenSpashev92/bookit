package service

import (
	"errors"

	"github.com/nurkenspashev92/bookit/internal/property/model"
)

var (
	ErrSlugExists       = model.ErrSlugExists
	ErrCategoryNotFound = model.ErrCategoryNotFound
	ErrHouseNotFound    = model.ErrHouseNotFound

	ErrMaxImagesExceeded = errors.New("maximum 15 images allowed")
	ErrImageTooLarge     = errors.New("image must be at most 5 MB")
	ErrImageNotFound     = errors.New("image not found")
)
