package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/booking/schema"
	"github.com/nurkenspashev92/bookit/internal/booking/service"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
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
	user := c.Locals("user").(identitymodel.User)

	var req schema.BookingCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}
	if err := req.Validate(); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	booking, err := h.bookingService.Create(c.Context(), user.ID, req)
	if err != nil {
		if errors.Is(err, service.ErrBookingOverlap) {
			return shared.Fail(c, http.StatusConflict, err)
		}
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	return c.Status(http.StatusCreated).JSON(booking)
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
	user := c.Locals("user").(identitymodel.User)

	bookings, err := h.bookingService.GetMyBookings(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, http.StatusInternalServerError, err)
	}
	if bookings == nil {
		bookings = []schema.BookingResponse{}
	}
	return c.JSON(bookings)
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
	user := c.Locals("user").(identitymodel.User)

	bookings, err := h.bookingService.GetOwnerBookings(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, http.StatusInternalServerError, err)
	}
	if bookings == nil {
		bookings = []schema.BookingResponse{}
	}
	return c.JSON(bookings)
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
	user := c.Locals("user").(identitymodel.User)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid id")
	}

	booking, err := h.bookingService.GetByID(c.Context(), id, user.ID)
	if err != nil {
		return shared.Fail(c, http.StatusNotFound, err)
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
	user := c.Locals("user").(identitymodel.User)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid id")
	}

	var req schema.BookingUpdateStatusRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}
	if err := req.Validate(); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	if err := h.bookingService.UpdateStatus(c.Context(), id, user.ID, req.Status); err != nil {
		return shared.Fail(c, http.StatusForbidden, err)
	}

	return c.JSON(shared.MessageResponse{Message: "status updated"})
}
