package repository

import (
	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
)

func assign[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

func applyHouseUpdate(h *model.House, req schema.HouseUpdateRequest) {
	applyHouseText(h, req)
	applyHouseNumbers(h, req)
	applyHouseFlags(h, req)
}

func applyHouseText(h *model.House, req schema.HouseUpdateRequest) {
	assign(&h.NameEN, req.NameEN)
	assign(&h.NameKZ, req.NameKZ)
	assign(&h.NameRU, req.NameRU)
	assign(&h.Slug, req.Slug)
	assign(&h.DescriptionEN, req.DescriptionEN)
	assign(&h.DescriptionKZ, req.DescriptionKZ)
	assign(&h.DescriptionRU, req.DescriptionRU)
	assign(&h.AddressEN, req.AddressEN)
	assign(&h.AddressKZ, req.AddressKZ)
	assign(&h.AddressRU, req.AddressRU)
	assign(&h.DistrictEN, req.DistrictEN)
	assign(&h.DistrictKZ, req.DistrictKZ)
	assign(&h.DistrictRU, req.DistrictRU)
	assign(&h.PhoneNumber, req.PhoneNumber)
}

func applyHouseNumbers(h *model.House, req schema.HouseUpdateRequest) {
	if req.Price != nil {
		h.Price = req.Price.Int()
	}
	if req.RoomsQty != nil {
		h.RoomsQty = req.RoomsQty.Int()
	}
	if req.GuestQty != nil {
		h.GuestQty = req.GuestQty.Int()
	}
	if req.BedroomQty != nil {
		h.BedroomQty = req.BedroomQty.Int()
	}
	if req.BathQty != nil {
		h.BathQty = req.BathQty.IntPtr()
	}
	if req.Priority != nil {
		h.Priority = req.Priority.Int()
	}
	if req.TypeID != nil {
		h.TypeID = req.TypeID.Int()
	}
	if req.CityID != nil {
		h.CityID = req.CityID.IntPtr()
	}
	if req.CountryID != nil {
		h.CountryID = req.CountryID.IntPtr()
	}
	if req.Lng != nil {
		h.Lng = req.Lng.Float64Ptr()
	}
	if req.Lat != nil {
		h.Lat = req.Lat.Float64Ptr()
	}
}

func applyHouseFlags(h *model.House, req schema.HouseUpdateRequest) {
	assign(&h.IsActive, req.IsActive)
	assign(&h.GuestsWithPets, req.GuestsWithPets)
	assign(&h.GuestsWithBabies, req.GuestsWithBabies)
	assign(&h.BestHouse, req.BestHouse)
	assign(&h.Promotion, req.Promotion)
	assign(&h.IsVerified, req.IsVerified)
	assign(&h.IsSale, req.IsSale)
	assign(&h.IsNewest, req.IsNewest)
	assign(&h.IsHot, req.IsHot)
	assign(&h.IsFeatured, req.IsFeatured)
	assign(&h.IsDiscount, req.IsDiscount)
}
