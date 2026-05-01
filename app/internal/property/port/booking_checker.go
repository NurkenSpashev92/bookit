package port

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
)

// BookingChecker is the port through which the property domain queries an
// active booking for a (house, user) pair. The booking domain provides an
// implementation. This inversion keeps property → booking out of the import graph.
type BookingChecker interface {
	GetUserActiveBooking(ctx context.Context, houseID, userID int) (*schema.HouseBooking, error)
}
