package schemas

// House full response
type House struct {
	NameEN         string       `json:"name_en" form:"name_en" swaggertype:"string" example:"Beach House"`
	NameKZ         string       `json:"name_kz" form:"name_kz" swaggertype:"string" example:"Жағажай үйі"`
	NameRU         string       `json:"name_ru" form:"name_ru" swaggertype:"string" example:"Пляжный дом"`
	Slug           string       `json:"slug" form:"slug" swaggertype:"string" example:"beach-house"`
	Price          FlexInt      `json:"price" form:"price" swaggertype:"integer" example:"50000" minimum:"0"`
	RoomsQty       FlexInt      `json:"rooms_qty" form:"rooms_qty" swaggertype:"integer" example:"3" minimum:"0"`
	GuestQty       FlexInt      `json:"guest_qty" form:"guest_qty" swaggertype:"integer" example:"6" minimum:"0"`
	BedroomQty     FlexInt      `json:"bedroom_qty" form:"bedroom_qty" swaggertype:"integer" example:"2" minimum:"0"`
	BathQty        *FlexInt     `json:"bath_qty" form:"bath_qty" swaggertype:"integer" example:"1" minimum:"0"`
	DescriptionEN  string       `json:"description_en" form:"description_en" swaggertype:"string" example:"A beautiful beach house"`
	DescriptionKZ  string       `json:"description_kz" form:"description_kz" swaggertype:"string"`
	DescriptionRU  string       `json:"description_ru" form:"description_ru" swaggertype:"string"`
	AddressEN      string       `json:"address_en" form:"address_en" swaggertype:"string" maxLength:"255" example:"123 Beach Rd"`
	AddressKZ      string       `json:"address_kz" form:"address_kz" swaggertype:"string" maxLength:"255"`
	AddressRU      string       `json:"address_ru" form:"address_ru" swaggertype:"string" maxLength:"255"`
	Lng            *FlexFloat64 `json:"lng" form:"lng" swaggertype:"number" example:"51.1694"`
	Lat            *FlexFloat64 `json:"lat" form:"lat" swaggertype:"number" example:"71.4491"`
	IsActive       bool         `json:"is_active" form:"is_active" example:"true"`
	Priority       FlexInt      `json:"priority" form:"priority" swaggertype:"integer" example:"0" minimum:"0"`
	OwnerID        FlexInt      `json:"owner_id" form:"owner_id" swaggertype:"integer" example:"1"`
	TypeID         FlexInt      `json:"type_id" form:"type_id" swaggertype:"integer" example:"1"`
	CityID         *FlexInt     `json:"city_id" form:"city_id" swaggertype:"integer" example:"1"`
	CountryID      *FlexInt     `json:"country_id" form:"country_id" swaggertype:"integer" example:"1"`
	GuestsWithPets bool         `json:"guests_with_pets" form:"guests_with_pets" example:"false"`
	BestHouse      bool         `json:"best_house" form:"best_house" example:"false"`
	Promotion      bool         `json:"promotion" form:"promotion" example:"false"`
	DistrictEN     string       `json:"district_en" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     string       `json:"district_kz" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     string       `json:"district_ru" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    string       `json:"phone_number" form:"phone_number" swaggertype:"string" maxLength:"20" example:"+77001234567"`
	Images         []Image      `json:"images" form:"images"`
}

// HouseCreateRequest create house request body
// @Description Request body for creating a new house
type HouseCreateRequest struct {
	NameEN         string       `json:"name_en" form:"name_en" swaggertype:"string" maxLength:"255" example:"Beach House" validate:"required"`
	NameKZ         string       `json:"name_kz" form:"name_kz" swaggertype:"string" maxLength:"255" example:"Жағажай үйі" validate:"required"`
	NameRU         string       `json:"name_ru" form:"name_ru" swaggertype:"string" maxLength:"255" example:"Пляжный дом" validate:"required"`
	Slug           string       `json:"slug" form:"slug" swaggertype:"string" maxLength:"255" example:"beach-house"`
	Price          FlexInt      `json:"price" form:"price" swaggertype:"integer" minimum:"0" example:"50000" validate:"required"`
	RoomsQty       FlexInt      `json:"rooms_qty" form:"rooms_qty" swaggertype:"integer" minimum:"0" example:"3"`
	GuestQty       FlexInt      `json:"guest_qty" form:"guest_qty" swaggertype:"integer" minimum:"0" example:"6"`
	BedroomQty     FlexInt      `json:"bedroom_qty" form:"bedroom_qty" swaggertype:"integer" minimum:"0" example:"2"`
	BathQty        *FlexInt     `json:"bath_qty" form:"bath_qty" swaggertype:"integer" minimum:"0" example:"1"`
	DescriptionEN  string       `json:"description_en" form:"description_en" swaggertype:"string" example:"A beautiful beach house" validate:"required"`
	DescriptionKZ  string       `json:"description_kz" form:"description_kz" swaggertype:"string" validate:"required"`
	DescriptionRU  string       `json:"description_ru" form:"description_ru" swaggertype:"string" validate:"required"`
	AddressEN      string       `json:"address_en" form:"address_en" swaggertype:"string" maxLength:"255" example:"123 Beach Rd" validate:"required"`
	AddressKZ      string       `json:"address_kz" form:"address_kz" swaggertype:"string" maxLength:"255" validate:"required"`
	AddressRU      string       `json:"address_ru" form:"address_ru" swaggertype:"string" maxLength:"255" validate:"required"`
	Lng            *FlexFloat64 `json:"lng" form:"lng" swaggertype:"number" example:"51.1694"`
	Lat            *FlexFloat64 `json:"lat" form:"lat" swaggertype:"number" example:"71.4491"`
	IsActive       bool         `json:"is_active" form:"is_active" example:"true"`
	Priority       FlexInt      `json:"priority" form:"priority" swaggertype:"integer" minimum:"0" example:"0"`
	OwnerID        int          `json:"owner_id" form:"owner_id" swaggerignore:"true"`
	TypeID         FlexInt      `json:"type_id" form:"type_id" swaggertype:"integer" example:"1" validate:"required"`
	CityID         *FlexInt     `json:"city_id" form:"city_id" swaggertype:"integer" example:"1"`
	CountryID      *FlexInt     `json:"country_id" form:"country_id" swaggertype:"integer" example:"1"`
	GuestsWithPets bool         `json:"guests_with_pets" form:"guests_with_pets" example:"false"`
	BestHouse      bool         `json:"best_house" form:"best_house" example:"false"`
	Promotion      bool         `json:"promotion" form:"promotion" example:"false"`
	DistrictEN     string       `json:"district_en" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     string       `json:"district_kz" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     string       `json:"district_ru" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    string       `json:"phone_number" form:"phone_number" swaggertype:"string" maxLength:"20" example:"+77001234567"`
}

func (r HouseCreateRequest) Validate() error {
	v := newValidator()
	v.required("name_en", r.NameEN)
	v.maxLen("name_en", r.NameEN, 255)
	v.required("name_kz", r.NameKZ)
	v.maxLen("name_kz", r.NameKZ, 255)
	v.required("name_ru", r.NameRU)
	v.maxLen("name_ru", r.NameRU, 255)
	v.maxLen("slug", r.Slug, 255)
	v.minInt("price", r.Price.Int(), 0)
	v.minInt("rooms_qty", r.RoomsQty.Int(), 0)
	v.minInt("guest_qty", r.GuestQty.Int(), 0)
	v.minInt("bedroom_qty", r.BedroomQty.Int(), 0)
	if r.BathQty != nil {
		v.minInt("bath_qty", r.BathQty.Int(), 0)
	}
	v.required("description_en", r.DescriptionEN)
	v.required("description_kz", r.DescriptionKZ)
	v.required("description_ru", r.DescriptionRU)
	v.required("address_en", r.AddressEN)
	v.maxLen("address_en", r.AddressEN, 255)
	v.required("address_kz", r.AddressKZ)
	v.maxLen("address_kz", r.AddressKZ, 255)
	v.required("address_ru", r.AddressRU)
	v.maxLen("address_ru", r.AddressRU, 255)
	v.requiredFlexInt("type_id", r.TypeID)
	v.maxLen("district_en", r.DistrictEN, 255)
	v.maxLen("district_kz", r.DistrictKZ, 255)
	v.maxLen("district_ru", r.DistrictRU, 255)
	v.maxLen("phone_number", r.PhoneNumber, 20)
	return v.result()
}

// HouseUpdateRequest partial update house request body
// @Description Request body for updating a house (all fields optional)
type HouseUpdateRequest struct {
	NameEN         *string      `json:"name_en,omitempty" form:"name_en" swaggertype:"string" maxLength:"255" example:"Beach House"`
	NameKZ         *string      `json:"name_kz,omitempty" form:"name_kz" swaggertype:"string" maxLength:"255"`
	NameRU         *string      `json:"name_ru,omitempty" form:"name_ru" swaggertype:"string" maxLength:"255"`
	Slug           *string      `json:"slug" form:"slug" swaggertype:"string" maxLength:"255"`
	Price          *FlexInt     `json:"price,omitempty" form:"price" swaggertype:"integer" minimum:"0" example:"50000"`
	RoomsQty       *FlexInt     `json:"rooms_qty,omitempty" form:"rooms_qty" swaggertype:"integer" minimum:"0"`
	GuestQty       *FlexInt     `json:"guest_qty,omitempty" form:"guest_qty" swaggertype:"integer" minimum:"0"`
	BedroomQty     *FlexInt     `json:"bedroom_qty,omitempty" form:"bedroom_qty" swaggertype:"integer" minimum:"0"`
	BathQty        *FlexInt     `json:"bath_qty,omitempty" form:"bath_qty" swaggertype:"integer" minimum:"0"`
	DescriptionEN  *string      `json:"description_en,omitempty" form:"description_en" swaggertype:"string"`
	DescriptionKZ  *string      `json:"description_kz,omitempty" form:"description_kz" swaggertype:"string"`
	DescriptionRU  *string      `json:"description_ru,omitempty" form:"description_ru" swaggertype:"string"`
	AddressEN      *string      `json:"address_en,omitempty" form:"address_en" swaggertype:"string" maxLength:"255"`
	AddressKZ      *string      `json:"address_kz,omitempty" form:"address_kz" swaggertype:"string" maxLength:"255"`
	AddressRU      *string      `json:"address_ru,omitempty" form:"address_ru" swaggertype:"string" maxLength:"255"`
	Lng            *FlexFloat64 `json:"lng,omitempty" form:"lng" swaggertype:"number" example:"51.1694"`
	Lat            *FlexFloat64 `json:"lat,omitempty" form:"lat" swaggertype:"number" example:"71.4491"`
	IsActive       *bool        `json:"is_active,omitempty" form:"is_active" example:"true"`
	Priority       *FlexInt     `json:"priority,omitempty" form:"priority" swaggertype:"integer" minimum:"0"`
	TypeID         *FlexInt     `json:"type_id,omitempty" form:"type_id" swaggertype:"integer" example:"1"`
	CityID         *FlexInt     `json:"city_id,omitempty" form:"city_id" swaggertype:"integer"`
	CountryID      *FlexInt     `json:"country_id,omitempty" form:"country_id" swaggertype:"integer"`
	GuestsWithPets *bool        `json:"guests_with_pets,omitempty" form:"guests_with_pets" example:"false"`
	BestHouse      *bool        `json:"best_house,omitempty" form:"best_house" example:"false"`
	Promotion      *bool        `json:"promotion,omitempty" form:"promotion" example:"false"`
	DistrictEN     *string      `json:"district_en,omitempty" form:"district_en" swaggertype:"string" maxLength:"255"`
	DistrictKZ     *string      `json:"district_kz,omitempty" form:"district_kz" swaggertype:"string" maxLength:"255"`
	DistrictRU     *string      `json:"district_ru,omitempty" form:"district_ru" swaggertype:"string" maxLength:"255"`
	PhoneNumber    *string      `json:"phone_number,omitempty" form:"phone_number" swaggertype:"string" maxLength:"20"`
}

func (r HouseUpdateRequest) Validate() error {
	v := newValidator()
	v.maxLenPtr("name_en", r.NameEN, 255)
	v.maxLenPtr("name_kz", r.NameKZ, 255)
	v.maxLenPtr("name_ru", r.NameRU, 255)
	v.maxLenPtr("slug", r.Slug, 255)
	if r.Price != nil {
		v.minInt("price", r.Price.Int(), 0)
	}
	if r.RoomsQty != nil {
		v.minInt("rooms_qty", r.RoomsQty.Int(), 0)
	}
	if r.GuestQty != nil {
		v.minInt("guest_qty", r.GuestQty.Int(), 0)
	}
	if r.BedroomQty != nil {
		v.minInt("bedroom_qty", r.BedroomQty.Int(), 0)
	}
	if r.BathQty != nil {
		v.minInt("bath_qty", r.BathQty.Int(), 0)
	}
	v.maxLenPtr("address_en", r.AddressEN, 255)
	v.maxLenPtr("address_kz", r.AddressKZ, 255)
	v.maxLenPtr("address_ru", r.AddressRU, 255)
	v.maxLenPtr("district_en", r.DistrictEN, 255)
	v.maxLenPtr("district_kz", r.DistrictKZ, 255)
	v.maxLenPtr("district_ru", r.DistrictRU, 255)
	v.maxLenPtr("phone_number", r.PhoneNumber, 20)
	return v.result()
}

// HouseListItem house list item for GET /houses
type HouseListItem struct {
	ID             int     `json:"id" example:"1"`
	NameEN         string  `json:"name_en" example:"Beach House"`
	NameKZ         string  `json:"name_kz" example:"Жағажай үйі"`
	NameRU         string  `json:"name_ru" example:"Пляжный дом"`
	Slug           string  `json:"slug" example:"beach-house"`
	Price          int     `json:"price" example:"50000"`
	AddressEN      string  `json:"address_en" example:"123 Beach Rd"`
	AddressKZ      string  `json:"address_kz"`
	AddressRU      string  `json:"address_ru"`
	Priority       int     `json:"priority" example:"0"`
	GuestsWithPets bool    `json:"guests_with_pets" example:"false"`
	BestHouse      bool    `json:"best_house" example:"false"`
	Promotion      bool    `json:"promotion" example:"false"`
	CountryCityNameKZ string `json:"country_city_name_kz" example:"Қазақстан, Астана"`
	CountryCityNameRU string `json:"country_city_name_ru" example:"Казахстан, Астана"`
	CountryCityNameEN string `json:"country_city_name_en" example:"Kazakhstan, Astana"`
	OwnerFullName  string  `json:"owner_full_name" example:"John Doe"`
	LikeCount      int     `json:"like_count" example:"5"`
	IsLiked     bool    `json:"is_liked" example:"false"`
	Images         []Image `json:"images"`
}

// SlugCheckResponse slug availability check response
type SlugCheckResponse struct {
	Available bool   `json:"available" example:"true"`
	Slug      string `json:"slug" example:"beach-house"`
}
