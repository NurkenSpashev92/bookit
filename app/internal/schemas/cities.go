package schemas

type City struct {
	ID          int     `json:"id"`
	NameKZ      string  `json:"name_kz"`
	NameEN      string  `json:"name_en"`
	NameRU      string  `json:"name_ru"`
	PostallCode string  `json:"postall_code,omitempty"`
	Country     Country `json:"country"`
}

type CityCreateRequest struct {
	NameKZ      string `json:"name_kz" form:"name_kz"`
	NameEN      string `json:"name_en" form:"name_en"`
	NameRU      string `json:"name_ru" form:"name_ru"`
	PostallCode string `json:"postall_code" form:"postall_code"`
	CountryID   int    `json:"country_id" form:"country_id"`
}

type CityUpdateRequest struct {
	NameKZ      *string `json:"name_kz,omitempty" form:"name_kz"`
	NameEN      *string `json:"name_en,omitempty" form:"name_en"`
	NameRU      *string `json:"name_ru,omitempty" form:"name_ru"`
	PostallCode *string `json:"postall_code,omitempty" form:"postall_code"`
	CountryID   *int    `json:"country_id,omitempty" form:"country_id"`
}
