package model

import "errors"

var (
	ErrHouseNotFound    = errors.New("house not found")
	ErrSlugExists       = errors.New("slug already exists")
	ErrCategoryNotFound = errors.New("one of category_ids does not exist")
)
