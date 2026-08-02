package port

import "context"

type LikeChecker interface {
	StatusWithCountByID(ctx context.Context, userID, houseID int) (bool, int, error)
	GetUserLikedHouseIDs(ctx context.Context, userID int, houseIDs []int) ([]int, error)
}
