package model

import "time"

type Image struct {
	ID        int       `json:"id"`
	Original  string    `json:"original,omitempty"`
	Thumbnail string    `json:"thumbnail,omitempty"`
	Width     *int      `json:"width,omitempty"`
	Height    *int      `json:"height,omitempty"`
	MimeType  string    `json:"mimetype,omitempty"`
	Size      *int      `json:"size,omitempty"`
	IsLabel   *bool     `json:"is_label,omitempty"`
	HouseID   *int      `json:"house_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
