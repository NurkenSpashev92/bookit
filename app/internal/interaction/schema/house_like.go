package schema

import "time"

// HouseLikeResponse like status response
type HouseLikeResponse struct {
	Liked     bool `json:"liked" example:"true"`
	LikeCount int  `json:"like_count" example:"5"`
}

// HouseLikeItem liked house item for user's liked houses list
type HouseLikeItem struct {
	ID        int       `json:"id" example:"1"`
	NameEN    string    `json:"name_en" example:"Beach House"`
	NameKZ    string    `json:"name_kz" example:"Жағажай үйі"`
	NameRU    string    `json:"name_ru" example:"Пляжный дом"`
	Slug      string    `json:"slug" example:"beach-house"`
	Price     int       `json:"price" example:"50000"`
	AddressEN string    `json:"address_en" example:"123 Beach Rd"`
	AddressKZ string    `json:"address_kz"`
	AddressRU string    `json:"address_ru"`
	LikedAt   time.Time `json:"liked_at" example:"2026-03-30T12:00:00Z" format:"date-time"`
}
