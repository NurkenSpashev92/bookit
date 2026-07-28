package router

import (
	"github.com/gofiber/fiber/v3"

	analyticsh "github.com/nurkenspashev92/bookit/internal/analytics/handler"
)

func registerAnalyticsRoutes(api fiber.Router, svc *Services, access guards) {
	handler := analyticsh.NewStatsHandler(svc.Stats)

	stats := api.Group("/stats")

	stats.Get("/dashboard", access.required, handler.Dashboard)
	stats.Get("/houses", access.required, handler.HouseStats)
	stats.Get("/houses/:slug", access.required, handler.HouseDetail)
	stats.Get("/charts", access.required, handler.Charts)
}
