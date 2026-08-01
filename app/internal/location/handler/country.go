package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CountryService interface {
	GetAll(ctx context.Context) ([]model.Country, error)
	GetAllPaginated(ctx context.Context, limit, offset int) ([]model.Country, int, error)
	GetByID(ctx context.Context, id int) (model.Country, error)
	Create(ctx context.Context, req schema.CountryCreateRequest) (model.Country, error)
	Update(ctx context.Context, id int, req schema.CountryUpdateRequest) (model.Country, error)
	Delete(ctx context.Context, id int) error
}

type CountryHandler struct {
	countryService CountryService
}

func NewCountryHandler(countryService CountryService) *CountryHandler {
	return &CountryHandler{countryService: countryService}
}

// GetCountries godoc
// @Summary Get all countries
// @Description Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags Countries
// @Produce json
// @Param page query int false "Page number (enables the paginated envelope)"
// @Param page_size query int false "Items per page (default 20, max 100)"
// @Success 200 {array} schema.Country
// @Failure 500 {object} shared.ErrorResponse
// @Router /countries [get]
func (h *CountryHandler) GetAll(c fiber.Ctx) error {
	return shared.ListMaybePaginated(c, h.countryService.GetAll, h.countryService.GetAllPaginated)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	country, err := h.countryService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
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
	request, err := shared.Bind[schema.CountryCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	country, err := h.countryService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, country)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.CountryUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	country, err := h.countryService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.countryService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "country deleted")
}
