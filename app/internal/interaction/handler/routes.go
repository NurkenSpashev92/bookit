package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func RegisterRoutes(api fiber.Router, service HouseLikeService, guards shared.Guards) {
	handler := NewHouseLikeHandler(service)

	houses := api.Group("/houses")

	houses.Get("/liked", guards.Required, handler.UserLikedHouses)
	houses.Post("/:slug/like", guards.Required, handler.Like)
	houses.Delete("/:slug/like", guards.Required, handler.Unlike)
	houses.Get("/:slug/like", guards.Required, handler.Status)
}
