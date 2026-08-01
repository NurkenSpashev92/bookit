package port

import "context"

// UserSubscriptionSetter lets the subscription domain keep the denormalized
// users.subscription_type column in sync with a user's current plan, without
// depending on the identity domain directly.
type UserSubscriptionSetter interface {
	SetSubscriptionType(ctx context.Context, userID int, subType string) error
}
