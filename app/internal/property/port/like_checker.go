package port

import "context"

type LikeChecker interface {
	StatusWithCount(ctx context.Context, userID int, slug string) (bool, int, error)
	GetUserLikedHouseIDs(ctx context.Context, userID int) ([]int, error)
}
