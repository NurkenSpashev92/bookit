package handlers

import (
	"encoding/json"
	"errors"
	"strconv"

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
// @Description  Returns a list of all houses
// @Tags         Houses
// @Produce      json
// @Success      200  {array}  schemas.HouseListItem
// @Failure      500  {object} schemas.ErrorResponse
// @Router       /houses [get]
func (h *HouseHandler) GetAll(c fiber.Ctx) error {
	houses, err := h.houseService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(houses)
}

// GetHouseByID godoc
// @Summary      Get house by ID
// @Description  Returns a single house by ID
// @Tags         Houses
// @Produce      json
// @Param        id   path      int  true  "House ID"
// @Success      200  {object} models.House
// @Failure      404  {object} schemas.ErrorResponse
// @Router       /houses/{id} [get]
func (h *HouseHandler) GetByID(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	house, err := h.houseService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: "house not found: " + err.Error()})
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
// @Param        id     path      int  true  "House ID"
// @Param        house  body      schemas.HouseUpdateRequest true  "House update data"
// @Success      200    {object} models.House
// @Failure      400    {object} schemas.ErrorResponse
// @Failure      401    {object} schemas.ErrorResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{id} [patch]
func (h *HouseHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var req schemas.HouseUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	house, err := h.houseService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(house)
}

// DeleteHouse godoc
// @Summary      Delete house
// @Description  Deletes a house by ID. Auth required.
// @Tags         Houses
// @Produce      json
// @Param        id   path      int  true  "House ID"
// @Success      200  {object} schemas.MessageResponse
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/{id} [delete]
func (h *HouseHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.houseService.Delete(c.Context(), id); err != nil {
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
