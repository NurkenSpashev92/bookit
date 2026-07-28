package router

import (
	"github.com/gofiber/fiber/v3"

	bookingh "github.com/nurkenspashev92/bookit/internal/booking/handler"
)

func registerBookingRoutes(api fiber.Router, svc *Services, access guards) {
	handler := bookingh.NewBookingHandler(svc.Booking)

	bookings := api.Group("/bookings")

	bookings.Post("/", access.required, handler.Create)
	bookings.Get("/", access.required, handler.GetMyBookings)
	bookings.Get("/owner", access.required, handler.GetOwnerBookings)
	bookings.Get("/:id", access.required, handler.GetByID)
	bookings.Patch("/:id/status", access.required, handler.UpdateStatus)
}
