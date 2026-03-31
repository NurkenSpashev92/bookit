package schemas

// BookingCreateRequest create booking request
// @Description Request body for creating a booking
type BookingCreateRequest struct {
	HouseSlug  string `json:"house_slug" example:"beach-house" validate:"required"`
	StartDate  string `json:"start_date" example:"2026-04-10" format:"date" validate:"required"`
	EndDate    string `json:"end_date" example:"2026-04-15" format:"date" validate:"required"`
	GuestCount int    `json:"guest_count" example:"4" minimum:"1"`
	Message    string `json:"message,omitempty" example:"Looking forward to staying!"`
}

func (r BookingCreateRequest) Validate() error {
	v := newValidator()
	v.required("house_slug", r.HouseSlug)
	v.required("start_date", r.StartDate)
	v.date("start_date", r.StartDate)
	v.required("end_date", r.EndDate)
	v.date("end_date", r.EndDate)
	if r.GuestCount < 1 {
		v.errs = append(v.errs, "guest_count must be at least 1")
	}
	return v.result()
}

// BookingResponse booking response DTO
type BookingResponse struct {
	ID             int    `json:"id" example:"1"`
	HouseID        int    `json:"house_id" example:"1"`
	HouseSlug      string `json:"house_slug" example:"beach-house"`
	HouseNameEN    string `json:"house_name_en" example:"Beach House"`
	HouseNameKZ    string `json:"house_name_kz"`
	HouseNameRU    string `json:"house_name_ru"`
	OwnerID        int    `json:"owner_id" example:"1"`
	OwnerFullName  string `json:"owner_full_name" example:"John Doe"`
	OwnerEmail     string `json:"owner_email" example:"owner@mail.com"`
	OwnerPhone     string `json:"owner_phone,omitempty" example:"+77001234567"`
	GuestID        int    `json:"guest_id" example:"3"`
	GuestFullName  string `json:"guest_full_name" example:"Jane Smith"`
	GuestEmail     string `json:"guest_email" example:"guest@mail.com"`
	GuestPhone     string `json:"guest_phone,omitempty"`
	StartDate      string `json:"start_date" example:"2026-04-10"`
	EndDate        string `json:"end_date" example:"2026-04-15"`
	GuestCount     int    `json:"guest_count" example:"4"`
	Status         string `json:"status" example:"pending"`
	TotalPrice     int    `json:"total_price" example:"250000"`
	Message        string `json:"message,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// BookingUpdateStatusRequest update booking status (for owner)
// @Description Request body for updating booking status
type BookingUpdateStatusRequest struct {
	Status string `json:"status" example:"confirmed" validate:"required"`
}

func (r BookingUpdateStatusRequest) Validate() error {
	v := newValidator()
	v.required("status", r.Status)
	valid := r.Status == "confirmed" || r.Status == "rejected" || r.Status == "cancelled" || r.Status == "pending"
	if !valid {
		v.errs = append(v.errs, "status must be one of: pending, confirmed, rejected, cancelled")
	}
	return v.result()
}
