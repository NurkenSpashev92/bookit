package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CityService interface {
	GetAll(ctx context.Context, search string) ([]schema.City, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.City, int, error)
	GetByID(ctx context.Context, id int) (schema.City, error)
	Create(ctx context.Context, req schema.CityCreateRequest) (model.City, error)
	Update(ctx context.Context, id int, req schema.CityUpdateRequest) (model.City, error)
	Delete(ctx context.Context, id int) error
}

type CityHandler struct {
	cityService CityService
}

func NewCityHandler(cityService CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

// GetCities godoc
// @Summary Get all cities
// @Description Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags Cities
// @Produce json
// @Param search query string false "Search by name/postal code"
// @Param page query int false "Page number (enables the paginated envelope)"
// @Param page_size query int false "Items per page (default 20, max 100)"
// @Success 200 {array} schema.City
// @Failure 500 {object} shared.ErrorResponse
// @Router /cities [get]
func (h *CityHandler) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.City, error) {
			return h.cityService.GetAll(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.City, int, error) {
			return h.cityService.GetAllPaginated(ctx, search, limit, offset)
		},
	)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	city, err := h.cityService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
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
	request, err := shared.Bind[schema.CityCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	city, err := h.cityService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, city)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.CityUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	city, err := h.cityService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
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
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.cityService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "city deleted")
}
