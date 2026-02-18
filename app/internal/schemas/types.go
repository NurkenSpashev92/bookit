package schemas

type TypeResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	IsActive bool    `json:"is_active"`
}

type TypeCreateRequest struct {
	Name     string `json:"name" form:"name"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

type TypeUpdateRequest struct {
	Name     *string `json:"name,omitempty" form:"name"`
	IsActive *bool   `json:"is_active,omitempty" form:"is_active"`
}
