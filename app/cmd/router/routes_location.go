package router

import (
	"github.com/gofiber/fiber/v3"

	locationh "github.com/nurkenspashev92/bookit/internal/location/handler"
)

func registerLocationRoutes(api fiber.Router, svc *Services, access guards) {
	registerReferenceRoutes(api.Group("/countries"), access, locationh.NewCountryHandler(svc.Country))
	registerReferenceRoutes(api.Group("/cities"), access, locationh.NewCityHandler(svc.City))
}
