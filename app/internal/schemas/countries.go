package schemas

type Country struct {
	ID     int    `json:"id"`
	NameKZ string `json:"name_kz" form:"name_kz"`
	NameEN string `json:"name_en" form:"name_en"`
	NameRU string `json:"name_ru" form:"name_ru"`
	Code   string `json:"code" form:"code"`
}

type CountryCreateRequest struct {
	NameKZ string `json:"name_kz" form:"name_kz"`
	NameEN string `json:"name_en" form:"name_en"`
	NameRU string `json:"name_ru" form:"name_ru"`
	Code   string `json:"code" form:"code"`
}

type CountryUpdateRequest struct {
	NameKZ *string `json:"name_kz,omitempty" form:"name_kz"`
	NameEN *string `json:"name_en,omitempty" form:"name_en"`
	NameRU *string `json:"name_ru,omitempty" form:"name_ru"`
	Code   *string `json:"code,omitempty" form:"code"`
}
