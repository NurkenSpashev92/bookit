package handler

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/internal/location/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CityHandler struct {
	cityService *service.CityService
}

func NewCityHandler(cityService *service.CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

// GetCities godoc
// @Summary Get all cities
// @Tags Cities
// @Produce json
// @Success 200 {array} schema.City
// @Failure 500 {object} shared.ErrorResponse
// @Router /cities [get]
func (h *CityHandler) GetAll(c fiber.Ctx) error {
	cities, err := h.cityService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(cities)
}

// GetCity godoc
// @Summary Get city by ID
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} schema.City
// @Failure 404 {object} shared.ErrorResponse
// @Router /cities/{id} [get]
func (h *CityHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	city, err := h.cityService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(city)
}

// CreateCity godoc
// @Summary Create a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param city body schema.CityCreateRequest true "City data"
// @Success 201 {object} schema.City
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities [post]
func (h *CityHandler) Create(c fiber.Ctx) error {
	var req schema.CityCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	city, err := h.cityService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(city)
}

// UpdateCity godoc
// @Summary Update a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param id path int true "City ID"
// @Param city body schema.CityUpdateRequest true "City data"
// @Success 200 {object} schema.City
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities/{id} [patch]
func (h *CityHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	var req schema.CityUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	city, err := h.cityService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(city)
}

// DeleteCity godoc
// @Summary Delete a city
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /cities/{id} [delete]
func (h *CityHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	if err := h.cityService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "city deleted"})
}
