package handler

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v3"

	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type HouseHandler struct {
	houseService *service.HouseService
}

func NewHouseHandler(houseService *service.HouseService) *HouseHandler {
	return &HouseHandler{houseService: houseService}
}

// GetHouses godoc
// @Summary      Get all houses
// @Description  Returns a filtered, paginated list of all houses
// @Tags         Houses
// @Produce      json
// @Param        page            query int    false "Page number" default(1)
// @Param        page_size       query int    false "Items per page" default(10)
// @Success      200  {object} shared.PaginatedResponse
// @Failure      500  {object} shared.ErrorResponse
// @Router       /houses [get]
func (h *HouseHandler) GetAll(c fiber.Ctx) error {
	var userID int
	if user, ok := c.Locals("user").(identitymodel.User); ok {
		userID = user.ID
	}

	p := shared.ParsePagination(c)
	filter := schema.ParseHouseFilter(c)

	houses, total, err := h.houseService.GetAllPaginated(c.Context(), userID, filter, p.PageSize, p.Offset)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []schema.HouseListItem{}
	}

	totalPages := total / p.PageSize
	if total%p.PageSize > 0 {
		totalPages++
	}

	return c.JSON(shared.PaginatedResponse{
		Data:       houses,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	})
}

// MyHouses godoc
// @Summary      Get current user's houses
// @Tags         Houses
// @Produce      json
// @Success      200  {object} shared.PaginatedResponse
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /my-houses [get]
func (h *HouseHandler) MyHouses(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	p := shared.ParsePagination(c)

	houses, total, err := h.houseService.GetMyHouses(c.Context(), user.ID, p.PageSize, p.Offset)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []schema.HouseListItem{}
	}

	totalPages := total / p.PageSize
	if total%p.PageSize > 0 {
		totalPages++
	}

	return c.JSON(shared.PaginatedResponse{
		Data:       houses,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	})
}

// GetBySlug godoc
// @Summary      Get house by slug
// @Tags         Houses
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Success      200  {object} schema.HouseDetailResponse
// @Failure      404  {object} shared.ErrorResponse
// @Router       /houses/{slug} [get]
func (h *HouseHandler) GetBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	var userID int
	if user, ok := c.Locals("user").(identitymodel.User); ok {
		userID = user.ID
	}

	house, err := h.houseService.GetBySlug(c.Context(), slug, userID, c.IP())
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: "house not found"})
	}

	return c.JSON(house)
}

// Create godoc
// @Summary      Create a new house
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        house  body  schema.HouseCreateRequest  true  "House data"
// @Success      201    {object} model.House
// @Failure      400    {object} shared.ErrorResponse
// @Failure      401    {object} shared.ErrorResponse
// @Failure      500    {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses [post]
func (h *HouseHandler) Create(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	var req schema.HouseCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	house, err := h.houseService.Create(c.Context(), req, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrSlugExists) {
			return c.Status(409).JSON(shared.ErrorResponse{Error: "slug already exists"})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(house)
}

// Update godoc
// @Summary      Update house
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Param        house  body      schema.HouseUpdateRequest true  "House update data"
// @Success      200    {object} model.House
// @Failure      400    {object} shared.ErrorResponse
// @Failure      401    {object} shared.ErrorResponse
// @Failure      500    {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{slug} [patch]
func (h *HouseHandler) Update(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	var req schema.HouseUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	house, err := h.houseService.Update(c.Context(), slug, req)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(house)
}

// Delete godoc
// @Summary      Delete house
// @Tags         Houses
// @Produce      json
// @Param        slug   path      string  true  "House slug"
// @Success      200  {object} shared.MessageResponse
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{slug} [delete]
func (h *HouseHandler) Delete(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	if err := h.houseService.Delete(c.Context(), slug); err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "house deleted"})
}

// CheckSlug godoc
// @Summary      Check house slug availability
// @Tags         Houses
// @Produce      json
// @Param        slug   query     string  true  "Slug to check"
// @Success      200    {object} schema.SlugCheckResponse
// @Failure      500    {object} shared.ErrorResponse
// @Router       /houses/check-slug [get]
func (h *HouseHandler) CheckSlug(c fiber.Ctx) error {
	available, normalized, err := h.houseService.CheckSlug(c.Context(), c.Query("slug"))
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(fiber.Map{
		"available": available,
		"slug":      normalized,
	})
}
