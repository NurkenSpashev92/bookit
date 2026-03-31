package handlers

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type HouseHandler struct {
	houseService *services.HouseService
}

func NewHouseHandler(houseService *services.HouseService) *HouseHandler {
	return &HouseHandler{houseService: houseService}
}

// GetHouses godoc
// @Summary      Get all houses
// @Description  Returns a filtered, paginated list of all houses
// @Tags         Houses
// @Produce      json
// @Param        page            query int    false "Page number" default(1)
// @Param        page_size       query int    false "Items per page" default(10)
// @Param        min_price       query int    false "Minimum price"
// @Param        max_price       query int    false "Maximum price"
// @Param        guest_count     query int    false "Minimum guest capacity"
// @Param        rooms_qty       query int    false "Minimum rooms"
// @Param        bedroom_qty     query int    false "Minimum bedrooms"
// @Param        bed_qty         query int    false "Minimum beds"
// @Param        bath_qty        query int    false "Minimum bathrooms"
// @Param        guests_with_pets query bool  false "Allows pets"
// @Param        category        query int    false "Category ID"
// @Param        house_type      query int    false "House type ID"
// @Param        country         query int    false "Country ID"
// @Param        city            query int    false "City ID"
// @Success      200  {object} schemas.PaginatedResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Router       /houses [get]
func (h *HouseHandler) GetAll(c fiber.Ctx) error {
	var userID int
	if user, ok := c.Locals("user").(models.User); ok {
		userID = user.ID
	}

	p := schemas.ParsePagination(c)
	filter := schemas.ParseHouseFilter(c)

	houses, total, err := h.houseService.GetAllPaginated(c.Context(), userID, filter, p.PageSize, p.Offset)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []schemas.HouseListItem{}
	}

	totalPages := total / p.PageSize
	if total%p.PageSize > 0 {
		totalPages++
	}

	return c.JSON(schemas.PaginatedResponse{
		Data:       houses,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	})
}

// MyHouses godoc
// @Summary      Get current user's houses
// @Description  Returns paginated houses owned by the authenticated user
// @Tags         Houses
// @Produce      json
// @Param        page      query int false "Page number" default(1)
// @Param        page_size query int false "Items per page" default(10)
// @Success      200  {object} schemas.PaginatedResponse
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /my-houses [get]
func (h *HouseHandler) MyHouses(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	p := schemas.ParsePagination(c)

	houses, total, err := h.houseService.GetMyHouses(c.Context(), user.ID, p.PageSize, p.Offset)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []schemas.HouseListItem{}
	}

	totalPages := total / p.PageSize
	if total%p.PageSize > 0 {
		totalPages++
	}

	return c.JSON(schemas.PaginatedResponse{
		Data:       houses,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	})
}

// GetHouseBySlug godoc
// @Summary      Get house by slug
// @Description  Returns a single house by slug
// @Tags         Houses
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Success      200  {object} models.House
// @Failure      404  {object} schemas.ErrorResponse
// @Router       /houses/{slug} [get]
func (h *HouseHandler) GetBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "slug is required"})
	}

	house, err := h.houseService.GetBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: "house not found"})
	}

	return c.JSON(house)
}

// CreateHouse godoc
// @Summary      Create a new house
// @Description  Creates a new house. Auth required.
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        house  body  schemas.HouseCreateRequest  true  "House data"
// @Success      201    {object} models.House
// @Failure      400    {object} schemas.ErrorResponse
// @Failure      401    {object} schemas.ErrorResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses [post]
func (h *HouseHandler) Create(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var req schemas.HouseCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	house, err := h.houseService.Create(c.Context(), req, user.ID)
	if err != nil {
		if errors.Is(err, services.ErrSlugExists) {
			return c.Status(409).JSON(schemas.ErrorResponse{Error: "slug already exists"})
		}
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(house)
}

// UpdateHouse godoc
// @Summary      Update house
// @Description  Updates house data. Auth required.
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Param        house  body      schemas.HouseUpdateRequest true  "House update data"
// @Success      200    {object} models.House
// @Failure      400    {object} schemas.ErrorResponse
// @Failure      401    {object} schemas.ErrorResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{slug} [patch]
func (h *HouseHandler) Update(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "slug is required"})
	}

	var req schemas.HouseUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	house, err := h.houseService.Update(c.Context(), slug, req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(house)
}

// DeleteHouse godoc
// @Summary      Delete house
// @Description  Deletes a house by slug. Auth required.
// @Tags         Houses
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Success      200  {object} schemas.MessageResponse
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{slug} [delete]
func (h *HouseHandler) Delete(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "slug is required"})
	}

	if err := h.houseService.Delete(c.Context(), slug); err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "house deleted"})
}

// CheckSlug godoc
// @Summary      Check house slug availability
// @Description  Checks if a house slug is available for use
// @Tags         Houses
// @Produce      json
// @Param        slug   query     string  true  "Slug to check"
// @Success      200    {object} schemas.SlugCheckResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Router       /houses/check-slug [get]
func (h *HouseHandler) CheckSlug(c fiber.Ctx) error {
	available, normalized, err := h.houseService.CheckSlug(c.Context(), c.Query("slug"))
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(fiber.Map{
		"available": available,
		"slug":      normalized,
	})
}
