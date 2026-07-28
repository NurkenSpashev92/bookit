package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func RegisterRoutes(api fiber.Router, service BookingService, guards shared.Guards) {
	handler := NewBookingHandler(service)

	bookings := api.Group("/bookings")

	bookings.Post("/", guards.Required, handler.Create)
	bookings.Get("/", guards.Required, handler.GetMyBookings)
	bookings.Get("/owner", guards.Required, handler.GetOwnerBookings)
	bookings.Get("/:id", guards.Required, handler.GetByID)
	bookings.Patch("/:id/status", guards.Required, handler.UpdateStatus)
}
