package port

import "context"

type UserSubscriptionSetter interface {
	SetSubscriptionType(ctx context.Context, userID int, subType string) error
}
