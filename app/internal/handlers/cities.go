package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type CityHandler struct {
	cityService *services.CityService
}

func NewCityHandler(cityService *services.CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

// GetCities godoc
// @Summary Get all cities
// @Tags Cities
// @Produce json
// @Success 200 {array} schemas.City
// @Failure 500 {object} schemas.ErrorResponse
// @Router /cities [get]
func (h *CityHandler) GetAll(c fiber.Ctx) error {
	cities, err := h.cityService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(cities)
}

// GetCity godoc
// @Summary Get city by ID
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} schemas.City
// @Failure 404 {object} schemas.ErrorResponse
// @Router /cities/{id} [get]
func (h *CityHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	city, err := h.cityService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(city)
}

// CreateCity godoc
// @Summary Create a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param city body schemas.CityCreateRequest true "City data"
// @Success 201 {object} schemas.City
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities [post]
func (h *CityHandler) Create(c fiber.Ctx) error {
	var req schemas.CityCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	city, err := h.cityService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(city)
}

// UpdateCity godoc
// @Summary Update a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param id path int true "City ID"
// @Param city body schemas.CityUpdateRequest true "City data"
// @Success 200 {object} schemas.City
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities/{id} [patch]
func (h *CityHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	var req schemas.CityUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	city, err := h.cityService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(city)
}

// DeleteCity godoc
// @Summary Delete a city
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities/{id} [delete]
func (h *CityHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	if err := h.cityService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "city deleted"})
}
