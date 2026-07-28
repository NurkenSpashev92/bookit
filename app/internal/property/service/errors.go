package service

import (
	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

var (
	ErrSlugExists         = model.ErrSlugExists
	ErrCategoryRefInvalid = model.ErrCategoryRefInvalid
	ErrHouseNotFound      = model.ErrHouseNotFound

	ErrMaxImagesExceeded = shared.Invalid("maximum 15 images allowed")
	ErrImageTooLarge     = shared.Invalid("image must be at most 5 MB")
	ErrImageNotFound     = shared.NotFound("image not found")
)
