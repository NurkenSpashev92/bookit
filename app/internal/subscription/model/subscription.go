package model

import "time"

type Type string

const (
	TypeBasic Type = "basic"
	TypePro   Type = "pro"
	TypeMax   Type = "max"
)

func (t Type) Valid() bool {
	switch t {
	case TypeBasic, TypePro, TypeMax:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "in_active"
)

func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusInactive:
		return true
	default:
		return false
	}
}

type Subscription struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	UserFullName string     `json:"user_full_name,omitempty"`
	Type         Type       `json:"type"`
	Status       Status     `json:"status"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
