package schemas

type House struct {
	NameEN         string  `json:"name_en" form:"name_en" binding:"required"`
	NameKZ         string  `json:"name_kz" form:"name_kz"`
	NameRU         string  `json:"name_ru" form:"name_ru"`
	Slug           string  `json:"slug" form:"slug"`
	Price          int     `json:"price" form:"price" binding:"required"`
	RoomsQty       int     `json:"rooms_qty" form:"rooms_qty"`
	GuestQty       int     `json:"guest_qty" form:"guest_qty"`
	BedroomQty     int     `json:"bedroom_qty" form:"bedroom_qty"`
	BathQty        *int    `json:"bath_qty" form:"bath_qty"`
	DescriptionEN  string  `json:"description_en" form:"description_en"`
	DescriptionKZ  string  `json:"description_kz" form:"description_kz"`
	DescriptionRU  string  `json:"description_ru" form:"description_ru"`
	AddressEN      string  `json:"address_en" form:"address_en"`
	AddressKZ      string  `json:"address_kz" form:"address_kz"`
	AddressRU      string  `json:"address_ru" form:"address_ru"`
	Lng            string  `json:"lng" form:"lng"`
	Lat            string  `json:"lat" form:"lat"`
	IsActive       bool    `json:"is_active" form:"is_active"`
	Priority       string  `json:"priority" form:"priority"`
	OwnerID        int     `json:"owner_id" form:"owner_id" binding:"required"`
	TypeID         int     `json:"type_id" form:"type_id"`
	CityID         *int    `json:"city_id" form:"city_id"`
	CountryID      *int    `json:"country_id" form:"country_id"`
	GuestsWithPets bool    `json:"guests_with_pets" form:"guests_with_pets"`
	BestHouse      bool    `json:"best_house" form:"best_house"`
	Promotion      bool    `json:"promotion" form:"promotion"`
	DistrictEN     string  `json:"district_en" form:"district_en"`
	DistrictKZ     string  `json:"district_kz" form:"district_kz"`
	DistrictRU     string  `json:"district_ru" form:"district_ru"`
	PhoneNumber    string  `json:"phone_number" form:"phone_number"`
	Images         []Image `json:"images" form:"images"`
}

type HouseCreateRequest struct {
	NameEN         string `json:"name_en" form:"name_en" binding:"required"`
	NameKZ         string `json:"name_kz" form:"name_kz"`
	NameRU         string `json:"name_ru" form:"name_ru"`
	Slug           string `json:"slug" form:"slug"`
	Price          int    `json:"price" form:"price" binding:"required"`
	RoomsQty       int    `json:"rooms_qty" form:"rooms_qty"`
	GuestQty       int    `json:"guest_qty" form:"guest_qty"`
	BedroomQty     int    `json:"bedroom_qty" form:"bedroom_qty"`
	BathQty        *int   `json:"bath_qty" form:"bath_qty"`
	DescriptionEN  string `json:"description_en" form:"description_en"`
	DescriptionKZ  string `json:"description_kz" form:"description_kz"`
	DescriptionRU  string `json:"description_ru" form:"description_ru"`
	AddressEN      string `json:"address_en" form:"address_en"`
	AddressKZ      string `json:"address_kz" form:"address_kz"`
	AddressRU      string `json:"address_ru" form:"address_ru"`
	Lng            string `json:"lng" form:"lng"`
	Lat            string `json:"lat" form:"lat"`
	IsActive       bool   `json:"is_active" form:"is_active"`
	Priority       string `json:"priority" form:"priority"`
	OwnerID        int    `json:"owner_id" form:"owner_id" binding:"required"`
	TypeID         int    `json:"type_id" form:"type_id"`
	CityID         *int   `json:"city_id" form:"city_id"`
	CountryID      *int   `json:"country_id" form:"country_id"`
	GuestsWithPets bool   `json:"guests_with_pets" form:"guests_with_pets"`
	BestHouse      bool   `json:"best_house" form:"best_house"`
	Promotion      bool   `json:"promotion" form:"promotion"`
	DistrictEN     string `json:"district_en" form:"district_en"`
	DistrictKZ     string `json:"district_kz" form:"district_kz"`
	DistrictRU     string `json:"district_ru" form:"district_ru"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
}

type HouseUpdateRequest struct {
	NameEN         *string `json:"name_en,omitempty" form:"name_en"`
	NameKZ         *string `json:"name_kz,omitempty" form:"name_kz"`
	NameRU         *string `json:"name_ru,omitempty" form:"name_ru"`
	Slug           *string `json:"slug" form:"slug"`
	Price          *int    `json:"price,omitempty" form:"price"`
	RoomsQty       *int    `json:"rooms_qty,omitempty" form:"rooms_qty"`
	GuestQty       *int    `json:"guest_qty,omitempty" form:"guest_qty"`
	BedroomQty     *int    `json:"bedroom_qty,omitempty" form:"bedroom_qty"`
	BathQty        *int    `json:"bath_qty,omitempty" form:"bath_qty"`
	DescriptionEN  *string `json:"description_en,omitempty" form:"description_en"`
	DescriptionKZ  *string `json:"description_kz,omitempty" form:"description_kz"`
	DescriptionRU  *string `json:"description_ru,omitempty" form:"description_ru"`
	AddressEN      *string `json:"address_en,omitempty" form:"address_en"`
	AddressKZ      *string `json:"address_kz,omitempty" form:"address_kz"`
	AddressRU      *string `json:"address_ru,omitempty" form:"address_ru"`
	Lng            *string `json:"lng,omitempty" form:"lng"`
	Lat            *string `json:"lat,omitempty" form:"lat"`
	IsActive       *bool   `json:"is_active,omitempty" form:"is_active"`
	Priority       *string `json:"priority,omitempty" form:"priority"`
	TypeID         *int    `json:"type_id,omitempty" form:"type_id"`
	CityID         *int    `json:"city_id,omitempty" form:"city_id"`
	CountryID      *int    `json:"country_id,omitempty" form:"country_id"`
	GuestsWithPets *bool   `json:"guests_with_pets,omitempty" form:"guests_with_pets"`
	BestHouse      *bool   `json:"best_house,omitempty" form:"best_house"`
	Promotion      *bool   `json:"promotion,omitempty" form:"promotion"`
	DistrictEN     *string `json:"district_en,omitempty" form:"district_en"`
	DistrictKZ     *string `json:"district_kz,omitempty" form:"district_kz"`
	DistrictRU     *string `json:"district_ru,omitempty" form:"district_ru"`
	PhoneNumber    *string `json:"phone_number,omitempty" form:"phone_number"`
}

type HouseListItem struct {
	ID     int    `json:"id"`
	NameEN string `json:"name_en"`
	NameKZ string `json:"name_kz"`
	NameRU string `json:"name_ru"`
	Slug   string `json:"slug" form:"slug"`

	Price int `json:"price"`

	AddressEN string `json:"address_en"`
	AddressKZ string `json:"address_kz"`
	AddressRU string `json:"address_ru"`

	Priority       string `json:"priority"`
	GuestsWithPets bool   `json:"guests_with_pets"`
	BestHouse      bool   `json:"best_house"`
	Promotion      bool   `json:"promotion"`

	CountryCityNameKZ string `json:"country_city_name_kz"`
	CountryCityNameRU string `json:"country_city_name_ru"`
	CountryCityNameEN string `json:"country_city_name_en"`

	OwnerFullName string `json:"owner_full_name"`

	Images []Image `json:"images" form:"images"`
}

type SlugCheckResponse struct {
	Available bool   `json:"available"`
	Slug      string `json:"slug"`
}
