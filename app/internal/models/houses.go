package models

import "time"

type House struct {
	ID             int       `json:"id"`
	NameEN         string    `json:"name_en"`
	NameKZ         string    `json:"name_kz"`
	NameRU         string    `json:"name_ru"`
	Slug           string    `json:"slug"`
	Price          int       `json:"price"`
	RoomsQty       int       `json:"rooms_qty"`
	GuestQty       int       `json:"guest_qty"`
	BedroomQty     int       `json:"bedroom_qty"`
	BathQty        *int      `json:"bath_qty,omitempty"`
	DescriptionEN  string    `json:"description_en"`
	DescriptionKZ  string    `json:"description_kz"`
	DescriptionRU  string    `json:"description_ru"`
	AddressEN      string    `json:"address_en"`
	AddressKZ      string    `json:"address_kz"`
	AddressRU      string    `json:"address_ru"`
	Lng            string    `json:"lng,omitempty"`
	Lat            string    `json:"lat,omitempty"`
	IsActive       bool      `json:"is_active"`
	Priority       string    `json:"priority"`
	LikeCount      int       `json:"like_count"`
	CommentsRU     *string   `json:"comments_ru,omitempty"`
	CommentsEN     *string   `json:"comments_en,omitempty"`
	CommentsKZ     *string   `json:"comments_kz,omitempty"`
	OwnerID        int       `json:"owner_id"`
	TypeID         int       `json:"type_id"`
	CityID         *int      `json:"city_id,omitempty"`
	CountryID      *int      `json:"country_id,omitempty"`
	GuestsWithPets bool      `json:"guests_with_pets"`
	BestHouse      bool      `json:"best_house"`
	Promotion      bool      `json:"promotion"`
	DistrictEN     string    `json:"district_en,omitempty"`
	DistrictKZ     string    `json:"district_kz,omitempty"`
	DistrictRU     string    `json:"district_ru,omitempty"`
	PhoneNumber    string    `json:"phone_number,omitempty"`
	Images         []Image   `json:"images,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
