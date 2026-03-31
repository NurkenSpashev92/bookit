package schemas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type HouseFilter struct {
	Name           *string
	MinPrice       *int
	MaxPrice       *int
	GuestCount     *int
	RoomsQty       *int
	BedroomQty     *int
	BedQty         *int
	BathQty        *int
	GuestsWithPets *bool
	CategoryID     *int
	TypeID         *int
	CountryID      *int
	CityID         *int
	OwnerID        *int // internal, not parsed from query
}

func ParseHouseFilter(c fiber.Ctx) HouseFilter {
	f := HouseFilter{}
	f.Name = queryString(c, "house_name")
	f.MinPrice = queryInt(c, "min_price")
	f.MaxPrice = queryInt(c, "max_price")
	f.GuestCount = queryInt(c, "guest_count")
	f.RoomsQty = queryInt(c, "rooms_qty")
	f.BedroomQty = queryInt(c, "bedroom_qty")
	f.BedQty = queryInt(c, "bed_qty")
	f.BathQty = queryInt(c, "bath_qty")
	f.CategoryID = queryInt(c, "category")
	f.TypeID = queryInt(c, "house_type")
	f.CountryID = queryInt(c, "country")
	f.CityID = queryInt(c, "city")
	f.GuestsWithPets = queryBool(c, "guests_with_pets")
	return f
}

func (f HouseFilter) IsEmpty() bool {
	return f.Name == nil && f.MinPrice == nil && f.MaxPrice == nil && f.GuestCount == nil &&
		f.RoomsQty == nil && f.BedroomQty == nil && f.BedQty == nil && f.BathQty == nil &&
		f.GuestsWithPets == nil && f.CategoryID == nil && f.TypeID == nil &&
		f.CountryID == nil && f.CityID == nil
}

func (f HouseFilter) CacheKey(limit, offset int) string {
	var parts []string
	if f.Name != nil {
		parts = append(parts, fmt.Sprintf("n%s", *f.Name))
	}
	if f.MinPrice != nil {
		parts = append(parts, fmt.Sprintf("mn%d", *f.MinPrice))
	}
	if f.MaxPrice != nil {
		parts = append(parts, fmt.Sprintf("mx%d", *f.MaxPrice))
	}
	if f.GuestCount != nil {
		parts = append(parts, fmt.Sprintf("g%d", *f.GuestCount))
	}
	if f.RoomsQty != nil {
		parts = append(parts, fmt.Sprintf("r%d", *f.RoomsQty))
	}
	if f.BedroomQty != nil {
		parts = append(parts, fmt.Sprintf("br%d", *f.BedroomQty))
	}
	if f.BedQty != nil {
		parts = append(parts, fmt.Sprintf("bd%d", *f.BedQty))
	}
	if f.BathQty != nil {
		parts = append(parts, fmt.Sprintf("bt%d", *f.BathQty))
	}
	if f.GuestsWithPets != nil && *f.GuestsWithPets {
		parts = append(parts, "pets")
	}
	if f.CategoryID != nil {
		parts = append(parts, fmt.Sprintf("cat%d", *f.CategoryID))
	}
	if f.TypeID != nil {
		parts = append(parts, fmt.Sprintf("t%d", *f.TypeID))
	}
	if f.CountryID != nil {
		parts = append(parts, fmt.Sprintf("cn%d", *f.CountryID))
	}
	if f.CityID != nil {
		parts = append(parts, fmt.Sprintf("ct%d", *f.CityID))
	}
	filter := strings.Join(parts, ".")
	return fmt.Sprintf("houses:%s:%d:%d", filter, limit, offset)
}

func queryString(c fiber.Ctx, key string) *string {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	return &s
}

func queryInt(c fiber.Ctx, key string) *int {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func queryBool(c fiber.Ctx, key string) *bool {
	s := c.Query(key)
	if s == "" {
		return nil
	}
	v := s == "true" || s == "1"
	return &v
}
