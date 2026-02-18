package schemas

type Category struct {
	Id       string  `json:"id"`
	NameKz   string  `json:"name_kz"`
	NameRu   string  `json:"name_ru"`
	NameEn   string  `json:"name_en"`
	Icon     *string `json:"icon"`
	IsActive bool    `json:"is_active"`
}

type CategoryCreateRequest struct {
	NameKz   string  `json:"name_kz"`
	NameRu   string  `json:"name_ru"`
	NameEn   string  `json:"name_en"`
	Icon     *string `json:"icon"`
	IsActive bool    `json:"is_active"`
}

type CategoryUpdateRequest struct {
	NameKz   *string `form:"name_kz"`
	NameRu   *string `form:"name_ru"`
	NameEn   *string `form:"name_en"`
	IsActive *bool   `form:"is_active"`
}

type CategoryPaginate struct {
	Id       string  `json:"id"`
	NameKz   string  `json:"name_kz"`
	NameRu   string  `json:"name_ru"`
	NameEn   string  `json:"name_en"`
	Icon     *string `json:"icon"`
	IsActive bool    `json:"is_active"`
}
