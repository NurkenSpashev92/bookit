package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func RegisterRoutes(api fiber.Router, service StatsService, guards shared.Guards) {
	handler := NewStatsHandler(service)

	stats := api.Group("/stats")

	stats.Get("/dashboard", guards.Required, handler.Dashboard)
	stats.Get("/houses", guards.Required, handler.HouseStats)
	stats.Get("/houses/:slug", guards.Required, handler.HouseDetail)
	stats.Get("/charts", guards.Required, handler.Charts)
}
