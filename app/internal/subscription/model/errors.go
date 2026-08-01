package model

import "github.com/nurkenspashev92/bookit/internal/shared"

var (
	ErrSubscriptionNotFound     = shared.NotFound("subscription not found")
	ErrUserRefInvalid           = shared.Invalid("user_id does not exist")
	ErrActiveSubscriptionExists = shared.Conflict("user already has an active subscription")
)
