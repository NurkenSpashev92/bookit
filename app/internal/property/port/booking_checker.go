package port

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
)

type BookingChecker interface {
	GetUserActiveBooking(ctx context.Context, houseID, userID int) (*schema.HouseBooking, error)
}
