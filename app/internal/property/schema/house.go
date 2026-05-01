package schema

import "github.com/nurkenspashev92/bookit/internal/shared"

// House full response
type House struct {
	NameEN         string              `json:"name_en" form:"name_en" swaggertype:"string" example:"Beach House"`
	NameKZ         string              `json:"name_kz" form:"name_kz" swaggertype:"string" example:"Жағажай үйі"`
	NameRU         string              `json:"name_ru" form:"name_ru" swaggertype:"string" example:"Пляжный дом"`
	Slug           string              `json:"slug" form:"slug" swaggertype:"string" example:"beach-house"`
	Price          shared.FlexInt      `json:"price" form:"price" swaggertype:"integer" example:"50000" minimum:"0"`
	RoomsQty       shared.FlexInt      `json:"rooms_qty" form:"rooms_qty" swaggertype:"integer" example:"3" minimum:"0"`
	GuestQty       shared.FlexInt      `json:"guest_qty" form:"guest_qty" swaggertype:"integer" example:"6" minimum:"0"`
	BedroomQty     shared.FlexInt      `json:"bedroom_qty" form:"bedroom_qty" swaggertype:"integer" example:"2" minimum:"0"`
	BathQty        *shared.FlexInt     `json:"bath_qty" form:"bath_qty" swaggertype:"integer" example:"1" minimum:"0"`
	DescriptionEN  string              `json:"description_en" form:"description_en" swaggertype:"string"`
	DescriptionKZ  string              `json:"description_kz" form:"description_kz" swaggertype:"string"`
	DescriptionRU  string              `json:"description_ru" form:"description_ru" swaggertype:"string"`
	AddressEN      string              `json:"address_en" form:"address_en" swaggertype:"string" maxLength:"255"`
	AddressKZ      string              `json:"address_kz" form:"address_kz" swaggertype:"string" maxLength:"255"`
	AddressRU      string              `json:"address_ru" form:"address_ru" swaggertype:"string" maxLength:"255"`
	Lng            *shared.FlexFloat64 `json:"lng" form:"lng" swaggertype:"number"`
	Lat            *shared.FlexFloat64 `json:"lat" form:"lat" swaggertype:"number"`
	IsActive       bool                `json:"is_active" form:"is_active"`
	Priority       shared.FlexInt      `json:"priority" form:"priority" swaggertype:"integer" minimum:"0"`
	OwnerID        shared.FlexInt      `json:"owner_id" form:"owner_id" swaggertype:"integer"`
	TypeID         shared.FlexInt      `json:"type_id" form:"type_id" swaggertype:"integer"`
	CityID         *shared.FlexInt     `json:"city_id" form:"city_id" swaggertype:"integer"`
	CountryID      *shared.FlexInt     `json:"country_id" form:"country_id" swaggertype:"integer"`
	GuestsWithPets bool                `json:"guests_with_pets" form:"guests_with_pets"`
	BestHouse      bool                `json:"best_house" form:"best_house"`
	Promotion      bool                `json:"promotion" form:"promotion"`
	DistrictEN     string              `json:"district_en" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     string              `json:"district_kz" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     string              `json:"district_ru" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    string              `json:"phone_number" form:"phone_number" swaggertype:"string" maxLength:"20"`
	Images         []Image             `json:"images" form:"images"`
}

// HouseCreateRequest create house request body
// @Description Request body for creating a new house
type HouseCreateRequest struct {
	NameEN         string              `json:"name_en" form:"name_en" swaggertype:"string" maxLength:"255" validate:"required"`
	NameKZ         string              `json:"name_kz" form:"name_kz" swaggertype:"string" maxLength:"255" validate:"required"`
	NameRU         string              `json:"name_ru" form:"name_ru" swaggertype:"string" maxLength:"255" validate:"required"`
	Slug           string              `json:"slug" form:"slug" swaggertype:"string" maxLength:"255"`
	Price          shared.FlexInt      `json:"price" form:"price" swaggertype:"integer" minimum:"0" validate:"required"`
	RoomsQty       shared.FlexInt      `json:"rooms_qty" form:"rooms_qty" swaggertype:"integer" minimum:"0"`
	GuestQty       shared.FlexInt      `json:"guest_qty" form:"guest_qty" swaggertype:"integer" minimum:"0"`
	BedroomQty     shared.FlexInt      `json:"bedroom_qty" form:"bedroom_qty" swaggertype:"integer" minimum:"0"`
	BathQty        *shared.FlexInt     `json:"bath_qty" form:"bath_qty" swaggertype:"integer" minimum:"0"`
	DescriptionEN  string              `json:"description_en" form:"description_en" swaggertype:"string" validate:"required"`
	DescriptionKZ  string              `json:"description_kz" form:"description_kz" swaggertype:"string" validate:"required"`
	DescriptionRU  string              `json:"description_ru" form:"description_ru" swaggertype:"string" validate:"required"`
	AddressEN      string              `json:"address_en" form:"address_en" swaggertype:"string" maxLength:"255" validate:"required"`
	AddressKZ      string              `json:"address_kz" form:"address_kz" swaggertype:"string" maxLength:"255" validate:"required"`
	AddressRU      string              `json:"address_ru" form:"address_ru" swaggertype:"string" maxLength:"255" validate:"required"`
	Lng            *shared.FlexFloat64 `json:"lng" form:"lng" swaggertype:"number"`
	Lat            *shared.FlexFloat64 `json:"lat" form:"lat" swaggertype:"number"`
	IsActive       bool                `json:"is_active" form:"is_active"`
	Priority       shared.FlexInt      `json:"priority" form:"priority" swaggertype:"integer" minimum:"0"`
	OwnerID        int                 `json:"owner_id" form:"owner_id" swaggerignore:"true"`
	TypeID         shared.FlexInt      `json:"type_id" form:"type_id" swaggertype:"integer" validate:"required"`
	CityID         *shared.FlexInt     `json:"city_id" form:"city_id" swaggertype:"integer"`
	CountryID      *shared.FlexInt     `json:"country_id" form:"country_id" swaggertype:"integer"`
	GuestsWithPets bool                `json:"guests_with_pets" form:"guests_with_pets"`
	BestHouse      bool                `json:"best_house" form:"best_house"`
	Promotion      bool                `json:"promotion" form:"promotion"`
	DistrictEN     string              `json:"district_en" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     string              `json:"district_kz" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     string              `json:"district_ru" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    string              `json:"phone_number" form:"phone_number" swaggertype:"string" maxLength:"20"`
}

func (r HouseCreateRequest) Validate() error {
	v := shared.NewValidator()
	v.Required("name_en", r.NameEN)
	v.MaxLen("name_en", r.NameEN, 255)
	v.Required("name_kz", r.NameKZ)
	v.MaxLen("name_kz", r.NameKZ, 255)
	v.Required("name_ru", r.NameRU)
	v.MaxLen("name_ru", r.NameRU, 255)
	v.MaxLen("slug", r.Slug, 255)
	v.MinInt("price", r.Price.Int(), 0)
	v.MinInt("rooms_qty", r.RoomsQty.Int(), 0)
	v.MinInt("guest_qty", r.GuestQty.Int(), 0)
	v.MinInt("bedroom_qty", r.BedroomQty.Int(), 0)
	if r.BathQty != nil {
		v.MinInt("bath_qty", r.BathQty.Int(), 0)
	}
	v.Required("description_en", r.DescriptionEN)
	v.Required("description_kz", r.DescriptionKZ)
	v.Required("description_ru", r.DescriptionRU)
	v.Required("address_en", r.AddressEN)
	v.MaxLen("address_en", r.AddressEN, 255)
	v.Required("address_kz", r.AddressKZ)
	v.MaxLen("address_kz", r.AddressKZ, 255)
	v.Required("address_ru", r.AddressRU)
	v.MaxLen("address_ru", r.AddressRU, 255)
	v.RequiredFlexInt("type_id", r.TypeID)
	v.MaxLen("district_en", r.DistrictEN, 255)
	v.MaxLen("district_kz", r.DistrictKZ, 255)
	v.MaxLen("district_ru", r.DistrictRU, 255)
	v.MaxLen("phone_number", r.PhoneNumber, 20)
	return v.Result()
}

// HouseUpdateRequest partial update house request body
// @Description Request body for updating a house (all fields optional)
type HouseUpdateRequest struct {
	NameEN         *string             `json:"name_en,omitempty" form:"name_en" swaggertype:"string" maxLength:"255"`
	NameKZ         *string             `json:"name_kz,omitempty" form:"name_kz" swaggertype:"string" maxLength:"255"`
	NameRU         *string             `json:"name_ru,omitempty" form:"name_ru" swaggertype:"string" maxLength:"255"`
	Slug           *string             `json:"slug" form:"slug" swaggertype:"string" maxLength:"255"`
	Price          *shared.FlexInt     `json:"price,omitempty" form:"price" swaggertype:"integer" minimum:"0"`
	RoomsQty       *shared.FlexInt     `json:"rooms_qty,omitempty" form:"rooms_qty" swaggertype:"integer" minimum:"0"`
	GuestQty       *shared.FlexInt     `json:"guest_qty,omitempty" form:"guest_qty" swaggertype:"integer" minimum:"0"`
	BedroomQty     *shared.FlexInt     `json:"bedroom_qty,omitempty" form:"bedroom_qty" swaggertype:"integer" minimum:"0"`
	BathQty        *shared.FlexInt     `json:"bath_qty,omitempty" form:"bath_qty" swaggertype:"integer" minimum:"0"`
	DescriptionEN  *string             `json:"description_en,omitempty" form:"description_en" swaggertype:"string"`
	DescriptionKZ  *string             `json:"description_kz,omitempty" form:"description_kz" swaggertype:"string"`
	DescriptionRU  *string             `json:"description_ru,omitempty" form:"description_ru" swaggertype:"string"`
	AddressEN      *string             `json:"address_en,omitempty" form:"address_en" swaggertype:"string" maxLength:"255"`
	AddressKZ      *string             `json:"address_kz,omitempty" form:"address_kz" swaggertype:"string" maxLength:"255"`
	AddressRU      *string             `json:"address_ru,omitempty" form:"address_ru" swaggertype:"string" maxLength:"255"`
	Lng            *shared.FlexFloat64 `json:"lng,omitempty" form:"lng" swaggertype:"number"`
	Lat            *shared.FlexFloat64 `json:"lat,omitempty" form:"lat" swaggertype:"number"`
	IsActive       *bool               `json:"is_active,omitempty" form:"is_active"`
	Priority       *shared.FlexInt     `json:"priority,omitempty" form:"priority" swaggertype:"integer" minimum:"0"`
	TypeID         *shared.FlexInt     `json:"type_id,omitempty" form:"type_id" swaggertype:"integer"`
	CityID         *shared.FlexInt     `json:"city_id,omitempty" form:"city_id" swaggertype:"integer"`
	CountryID      *shared.FlexInt     `json:"country_id,omitempty" form:"country_id" swaggertype:"integer"`
	GuestsWithPets *bool               `json:"guests_with_pets,omitempty" form:"guests_with_pets"`
	BestHouse      *bool               `json:"best_house,omitempty" form:"best_house"`
	Promotion      *bool               `json:"promotion,omitempty" form:"promotion"`
	DistrictEN     *string             `json:"district_en,omitempty" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     *string             `json:"district_kz,omitempty" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     *string             `json:"district_ru,omitempty" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    *string             `json:"phone_number,omitempty" form:"phone_number" swaggertype:"string" maxLength:"20"`
}

func (r HouseUpdateRequest) Validate() error {
	v := shared.NewValidator()
	v.MaxLenPtr("name_en", r.NameEN, 255)
	v.MaxLenPtr("name_kz", r.NameKZ, 255)
	v.MaxLenPtr("name_ru", r.NameRU, 255)
	v.MaxLenPtr("slug", r.Slug, 255)
	if r.Price != nil {
		v.MinInt("price", r.Price.Int(), 0)
	}
	if r.RoomsQty != nil {
		v.MinInt("rooms_qty", r.RoomsQty.Int(), 0)
	}
	if r.GuestQty != nil {
		v.MinInt("guest_qty", r.GuestQty.Int(), 0)
	}
	if r.BedroomQty != nil {
		v.MinInt("bedroom_qty", r.BedroomQty.Int(), 0)
	}
	if r.BathQty != nil {
		v.MinInt("bath_qty", r.BathQty.Int(), 0)
	}
	v.MaxLenPtr("address_en", r.AddressEN, 255)
	v.MaxLenPtr("address_kz", r.AddressKZ, 255)
	v.MaxLenPtr("address_ru", r.AddressRU, 255)
	v.MaxLenPtr("district_en", r.DistrictEN, 255)
	v.MaxLenPtr("district_kz", r.DistrictKZ, 255)
	v.MaxLenPtr("district_ru", r.DistrictRU, 255)
	v.MaxLenPtr("phone_number", r.PhoneNumber, 20)
	return v.Result()
}

// HouseListItem house list item for GET /houses
type HouseListItem struct {
	ID                int     `json:"id" example:"1"`
	NameEN            string  `json:"name_en"`
	NameKZ            string  `json:"name_kz"`
	NameRU            string  `json:"name_ru"`
	Slug              string  `json:"slug"`
	Price             int     `json:"price"`
	AddressEN         string  `json:"address_en"`
	AddressKZ         string  `json:"address_kz"`
	AddressRU         string  `json:"address_ru"`
	Priority          int     `json:"priority"`
	GuestsWithPets    bool    `json:"guests_with_pets"`
	BestHouse         bool    `json:"best_house"`
	Promotion         bool    `json:"promotion"`
	CountryCityNameKZ string  `json:"country_city_name_kz"`
	CountryCityNameRU string  `json:"country_city_name_ru"`
	CountryCityNameEN string  `json:"country_city_name_en"`
	OwnerFullName     string  `json:"owner_full_name"`
	LikeCount         int     `json:"like_count"`
	IsLiked           bool    `json:"is_liked"`
	Images            []Image `json:"images"`
}

// HouseDetailResponse full house detail for GET /houses/{slug}
type HouseDetailResponse struct {
	ID                int           `json:"id"`
	NameEN            string        `json:"name_en"`
	NameKZ            string        `json:"name_kz"`
	NameRU            string        `json:"name_ru"`
	Slug              string        `json:"slug"`
	Price             int           `json:"price"`
	RoomsQty          int           `json:"rooms_qty"`
	GuestQty          int           `json:"guest_qty"`
	BedroomQty        int           `json:"bedroom_qty"`
	BathQty           *int          `json:"bath_qty,omitempty"`
	DescriptionEN     string        `json:"description_en"`
	DescriptionKZ     string        `json:"description_kz"`
	DescriptionRU     string        `json:"description_ru"`
	AddressEN         string        `json:"address_en"`
	AddressKZ         string        `json:"address_kz"`
	AddressRU         string        `json:"address_ru"`
	Lng               *float64      `json:"lng,omitempty"`
	Lat               *float64      `json:"lat,omitempty"`
	IsActive          bool          `json:"is_active"`
	Priority          int           `json:"priority"`
	CommentsRU        *string       `json:"comments_ru,omitempty"`
	CommentsEN        *string       `json:"comments_en,omitempty"`
	CommentsKZ        *string       `json:"comments_kz,omitempty"`
	OwnerID           int           `json:"owner_id"`
	TypeID            int           `json:"type_id"`
	CityID            *int          `json:"city_id,omitempty"`
	CountryID         *int          `json:"country_id,omitempty"`
	GuestsWithPets    bool          `json:"guests_with_pets"`
	BestHouse         bool          `json:"best_house"`
	Promotion         bool          `json:"promotion"`
	DistrictEN        string        `json:"district_en,omitempty"`
	DistrictKZ        string        `json:"district_kz,omitempty"`
	DistrictRU        string        `json:"district_ru,omitempty"`
	PhoneNumber       string        `json:"phone_number,omitempty"`
	LikeCount         int           `json:"like_count"`
	IsLiked           bool          `json:"is_liked"`
	CountryCityNameKZ string        `json:"country_city_name_kz,omitempty"`
	CountryCityNameRU string        `json:"country_city_name_ru,omitempty"`
	CountryCityNameEN string        `json:"country_city_name_en,omitempty"`
	OwnerFullName     string        `json:"owner_full_name,omitempty"`
	IsBooked          bool          `json:"is_booked"`
	MyBooking         *HouseBooking `json:"my_booking,omitempty"`
	Images            []Image       `json:"images"`
	CreatedAt         string        `json:"created_at"`
	UpdatedAt         string        `json:"updated_at"`
}

// HouseBooking current user's active booking for this house
type HouseBooking struct {
	ID         int    `json:"id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Status     string `json:"status"`
	GuestCount int    `json:"guest_count"`
	TotalPrice int    `json:"total_price"`
}

// SlugCheckResponse slug availability check response
type SlugCheckResponse struct {
	Available bool   `json:"available"`
	Slug      string `json:"slug"`
}
