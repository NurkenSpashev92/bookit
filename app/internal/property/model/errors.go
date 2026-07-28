package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrHouseNotFound      = shared.NotFound("house not found")
	ErrSlugExists         = shared.Conflict("slug already exists")
	ErrCategoryNotFound   = shared.NotFound("category not found")
	ErrCategoryRefInvalid = shared.Invalid("one of category_ids does not exist")
	ErrTypeNotFound       = shared.NotFound("type not found")
)
