package handler

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/internal/location/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CountryHandler struct {
	countryService *service.CountryService
}

func NewCountryHandler(countryService *service.CountryService) *CountryHandler {
	return &CountryHandler{countryService: countryService}
}

// GetCountries godoc
// @Summary Get all countries
// @Tags Countries
// @Produce json
// @Success 200 {array} schema.Country
// @Failure 500 {object} shared.ErrorResponse
// @Router /countries [get]
func (h *CountryHandler) GetAll(c fiber.Ctx) error {
	countries, err := h.countryService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(countries)
}

// GetCountry godoc
// @Summary Get country by ID
// @Tags Countries
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} schema.Country
// @Failure 404 {object} shared.ErrorResponse
// @Router /countries/{id} [get]
func (h *CountryHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	country, err := h.countryService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: "country not found"})
	}

	return c.JSON(country)
}

// CreateCountry godoc
// @Summary Create a country
// @Tags Countries
// @Accept json
// @Produce json
// @Param country body schema.CountryCreateRequest true "Country data"
// @Success 201 {object} schema.Country
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /countries [post]
func (h *CountryHandler) Create(c fiber.Ctx) error {
	var req schema.CountryCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	country, err := h.countryService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(country)
}

// UpdateCountry godoc
// @Summary Update a country
// @Tags Countries
// @Accept json
// @Produce json
// @Param id path int true "Country ID"
// @Param country body schema.CountryUpdateRequest true "Country data"
// @Success 200 {object} schema.Country
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /countries/{id} [patch]
func (h *CountryHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	var req schema.CountryUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	country, err := h.countryService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(country)
}

// DeleteCountry godoc
// @Summary Delete a country
// @Tags Countries
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /countries/{id} [delete]
func (h *CountryHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid id"})
	}

	if err := h.countryService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "country deleted"})
}
