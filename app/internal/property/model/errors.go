package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrHouseNotFound       = shared.NotFound("house not found")
	ErrHouseForbidden      = shared.Forbidden("not allowed to modify this house")
	ErrSlugExists          = shared.Conflict("slug already exists")
	ErrCategoryNotFound    = shared.NotFound("category not found")
	ErrCategoryRefInvalid  = shared.Invalid("one of category_ids does not exist")
	ErrTypeNotFound        = shared.NotFound("type not found")
	ErrConvenienceNotFound = shared.NotFound("convenience not found")

	ErrCategorySlugExists    = shared.Conflict("category slug already exists")
	ErrTypeSlugExists        = shared.Conflict("type slug already exists")
	ErrConvenienceSlugExists = shared.Conflict("convenience slug already exists")
)
