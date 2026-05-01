package model

import "time"

type HouseLike struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	HouseID   int       `json:"house_id"`
	CreatedAt time.Time `json:"created_at"`
}
