package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrCountryNotFound = shared.NotFound("country not found")
	ErrCityNotFound    = shared.NotFound("city not found")
)
