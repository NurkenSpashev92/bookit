package port

import "context"

// LikeChecker is the port through which the property domain queries a user's
// like status for a single house, plus the IDs of houses they have liked.
// The interaction domain provides an implementation.
type LikeChecker interface {
	StatusWithCount(ctx context.Context, userID int, slug string) (bool, int, error)
	GetUserLikedHouseIDs(ctx context.Context, userID int) ([]int, error)
}
