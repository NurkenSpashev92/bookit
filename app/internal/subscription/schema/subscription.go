package schema

import (
	"time"

	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/internal/subscription/model"
)

const (
	typeEnumMsg   = "type must be one of: basic, pro, max"
	statusEnumMsg = "status must be one of: active, in_active"
)

type SubscriptionResponse struct {
	ID        int        `json:"id" example:"1"`
	UserID    int        `json:"user_id" example:"1"`
	Type      string     `json:"type" example:"basic"`
	Status    string     `json:"status" example:"active"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// SubscriptionActivateRequest activate/change the current user's own plan.
// @Description Request body for activating a subscription plan
type SubscriptionActivateRequest struct {
	Type         string `json:"type" example:"pro" validate:"required"`
	DurationDays int    `json:"duration_days,omitempty" example:"30"`
}

func (r SubscriptionActivateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("type", r.Type)
	if r.Type != "" && !model.Type(r.Type).Valid() {
		v.Append(typeEnumMsg)
	}
	if r.DurationDays < 0 {
		v.Append("duration_days must be positive")
	}
	return v.Result()
}

type SubscriptionActivationResponse struct {
	Message       string               `json:"message" example:"subscription plan activated"`
	AlreadyActive bool                 `json:"already_active" example:"false"`
	Subscription  SubscriptionResponse `json:"subscription"`
}

// SubscriptionCreateRequest create a subscription for a user (admin).
// @Description Request body for creating a subscription
type SubscriptionCreateRequest struct {
	UserID  int        `json:"user_id" example:"1" validate:"required"`
	Type    string     `json:"type,omitempty" example:"basic"`
	Status  string     `json:"status,omitempty" example:"active"`
	EndDate *time.Time `json:"end_date,omitempty"`
}

func (r SubscriptionCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.RequiredInt("user_id", r.UserID)
	if r.Type != "" && !model.Type(r.Type).Valid() {
		v.Append(typeEnumMsg)
	}
	if r.Status != "" && !model.Status(r.Status).Valid() {
		v.Append(statusEnumMsg)
	}
	return v.Result()
}

// SubscriptionUpdateRequest partial update of a subscription (admin).
// @Description Request body for updating a subscription (all fields optional)
type SubscriptionUpdateRequest struct {
	Type    *string    `json:"type,omitempty" example:"pro"`
	Status  *string    `json:"status,omitempty" example:"in_active"`
	EndDate *time.Time `json:"end_date,omitempty"`
}

func (r SubscriptionUpdateRequest) Validate() error {
	v := shared.NewValidator()
	if r.Type != nil && !model.Type(*r.Type).Valid() {
		v.Append(typeEnumMsg)
	}
	if r.Status != nil && !model.Status(*r.Status).Valid() {
		v.Append(statusEnumMsg)
	}
	return v.Result()
}
