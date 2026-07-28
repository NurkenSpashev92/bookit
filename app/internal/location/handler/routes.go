package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

type Deps struct {
	Country CountryService
	City    CityService
}

func RegisterRoutes(api fiber.Router, deps Deps, guards shared.Guards) {
	shared.RegisterCRUD(api.Group("/countries"), guards, NewCountryHandler(deps.Country))
	shared.RegisterCRUD(api.Group("/cities"), guards, NewCityHandler(deps.City))
}
