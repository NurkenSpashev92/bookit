package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/booking/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type BookingService interface {
	Create(ctx context.Context, userID int, req schema.BookingCreateRequest) (schema.BookingResponse, error)
	GetMyBookings(ctx context.Context, userID int) ([]schema.BookingResponse, error)
	GetOwnerBookings(ctx context.Context, ownerID int) ([]schema.BookingResponse, error)
	GetByID(ctx context.Context, id, userID int) (schema.BookingResponse, error)
	UpdateStatus(ctx context.Context, bookingID, userID int, status string) error
}

type BookingHandler struct {
	bookingService BookingService
}

func NewBookingHandler(bookingService BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// Create godoc
// @Summary      Create a booking
// @Description  Book a house for specific dates
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Param        booking body schema.BookingCreateRequest true "Booking data"
// @Success      201 {object} schema.BookingResponse
// @Failure      400 {object} shared.ErrorResponse
// @Failure      401 {object} shared.ErrorResponse
// @Failure      409 {object} shared.ErrorResponse
// @Failure      500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /bookings [post]
func (h *BookingHandler) Create(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.BookingCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	booking, err := h.bookingService.Create(c.Context(), user.ID, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, booking)
}

// GetMyBookings godoc
// @Summary      Get my bookings
// @Description  Returns bookings made by the authenticated user
// @Tags         Bookings
// @Produce      json
// @Success      200 {array} schema.BookingResponse
// @Failure      401 {object} shared.ErrorResponse
// @Failure      500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /bookings [get]
func (h *BookingHandler) GetMyBookings(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	bookings, err := h.bookingService.GetMyBookings(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.List(c, bookings)
}

// GetOwnerBookings godoc
// @Summary      Get bookings for my houses
// @Description  Returns all bookings for houses owned by the authenticated user
// @Tags         Bookings
// @Produce      json
// @Success      200 {array} schema.BookingResponse
// @Failure      401 {object} shared.ErrorResponse
// @Failure      500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /bookings/owner [get]
func (h *BookingHandler) GetOwnerBookings(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	bookings, err := h.bookingService.GetOwnerBookings(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.List(c, bookings)
}

// GetByID godoc
// @Summary      Get booking by ID
// @Description  Returns a booking (visible to guest and house owner)
// @Tags         Bookings
// @Produce      json
// @Param        id path int true "Booking ID"
// @Success      200 {object} schema.BookingResponse
// @Failure      401 {object} shared.ErrorResponse
// @Failure      404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /bookings/{id} [get]
func (h *BookingHandler) GetByID(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	booking, err := h.bookingService.GetByID(c.Context(), id, user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(booking)
}

// UpdateStatus godoc
// @Summary      Update booking status
// @Description  Owner can confirm/reject. Guest can cancel.
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Param        id path int true "Booking ID"
// @Param        body body schema.BookingUpdateStatusRequest true "New status"
// @Success      200 {object} shared.MessageResponse
// @Failure      400 {object} shared.ErrorResponse
// @Failure      401 {object} shared.ErrorResponse
// @Failure      403 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /bookings/{id}/status [patch]
func (h *BookingHandler) UpdateStatus(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.BookingUpdateStatusRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.bookingService.UpdateStatus(c.Context(), id, user.ID, request.Status); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "status updated")
}
