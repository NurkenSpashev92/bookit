package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrHouseNotFound      = shared.NotFound("house not found")
	ErrBookingNotFound    = shared.NotFound("booking not found")
	ErrBookingOverlap     = shared.Conflict("these dates are already booked")
	ErrInvalidDateRange   = shared.Invalid("end_date must be after start_date")
	ErrStartDateInPast    = shared.Invalid("start_date cannot be in the past")
	ErrCancelNotAllowed   = shared.Forbidden("only the guest can cancel a booking")
	ErrDecisionNotAllowed = shared.Forbidden("only the owner can confirm or reject a booking")
)
