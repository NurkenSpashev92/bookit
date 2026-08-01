package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func RegisterRoutes(api fiber.Router, service SubscriptionService, guards shared.Guards) {
	handler := NewSubscriptionHandler(service)

	group := api.Group("/subscriptions")
	group.Get("/me", guards.Required, handler.GetMy)
	group.Post("/activate", guards.Required, handler.Activate)
	group.Post("/cancel", guards.Required, handler.Cancel)
	group.Get("/user/:userId", guards.Admin, handler.ByUser)
	group.Get("/", guards.Admin, handler.List)
	group.Post("/", guards.Admin, handler.Create)
	group.Get("/:id", guards.Admin, handler.GetByID)
	group.Patch("/:id", guards.Admin, handler.Update)
	group.Delete("/:id", guards.Admin, handler.Delete)
}
