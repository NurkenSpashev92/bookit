package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type HouseService interface {
	GetAllPaginated(ctx context.Context, userID int, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error)
	GetMyHouses(ctx context.Context, ownerID, limit, offset int) ([]schema.HouseListItem, int, error)
	GetForModeration(ctx context.Context, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error)
	GetBySlug(ctx context.Context, slug string, userID int, ip string) (schema.HouseDetailResponse, error)
	Create(ctx context.Context, req schema.HouseCreateRequest, ownerID int, isAdmin bool) (model.House, error)
	Update(ctx context.Context, slug string, req schema.HouseUpdateRequest, actorID int, isAdmin bool) (model.House, error)
	Delete(ctx context.Context, slug string, actorID int, isAdmin bool) error
	CheckSlug(ctx context.Context, rawSlug string) (bool, string, error)
}

type HouseHandler struct {
	houseService HouseService
}

func NewHouseHandler(houseService HouseService) *HouseHandler {
	return &HouseHandler{houseService: houseService}
}

// GetHouses godoc
// @Summary      Get all houses
// @Description  Returns a filtered, paginated list of active (approved) houses
// @Tags         Houses
// @Produce      json
// @Param        page            query int    false "Page number" default(1)
// @Param        page_size       query int    false "Items per page" default(10)
// @Success      200  {object} shared.PaginatedResponse
// @Failure      500  {object} shared.ErrorResponse
// @Router       /houses [get]
func (h *HouseHandler) GetAll(c fiber.Ctx) error {
	page := shared.Page(c)
	filter := schema.ParseHouseFilter(c)

	houses, total, err := h.houseService.GetAllPaginated(c.Context(), middleware.CurrentUserID(c), filter, page.Limit, page.Start())
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(shared.Paginated(shared.Items(houses), total, page))
}

// MyHouses godoc
// @Summary      Get current user's houses
// @Description  Returns the current user's houses, including inactive (pending approval) listings
// @Tags         Houses
// @Produce      json
// @Success      200  {object} shared.PaginatedResponse
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /my-houses [get]
func (h *HouseHandler) MyHouses(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	page := shared.Page(c)

	houses, total, err := h.houseService.GetMyHouses(c.Context(), user.ID, page.Limit, page.Start())
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(shared.Paginated(shared.Items(houses), total, page))
}

// Moderation godoc
// @Summary      List houses for moderation (admin)
// @Description  Admin only. Returns ALL houses including inactive/pending ones. Filter by ?is_active=true|false; omit to see all. Activate a house via PATCH /houses/{slug} with {"is_active": true}.
// @Tags         Houses
// @Produce      json
// @Param        is_active  query bool  false "Filter by active status (omit = all)"
// @Param        page       query int   false "Page"
// @Param        page_size  query int   false "Items per page"
// @Success      200  {object} shared.PaginatedResponse
// @Failure      401  {object} shared.ErrorResponse
// @Failure      403  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses/moderation [get]
func (h *HouseHandler) Moderation(c fiber.Ctx) error {
	page := shared.Page(c)
	filter := schema.ParseHouseFilter(c)

	houses, total, err := h.houseService.GetForModeration(c.Context(), filter, page.Limit, page.Start())
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(shared.Paginated(shared.Items(houses), total, page))
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
	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return shared.Fail(c, err)
	}

	house, err := h.houseService.GetBySlug(c.Context(), slug, middleware.CurrentUserID(c), c.IP())
	if err != nil {
		return shared.Fail(c, err)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.HouseCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	house, err := h.houseService.Create(c.Context(), request, user.ID, user.IsSuperuser)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, house)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.HouseUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	house, err := h.houseService.Update(c.Context(), slug, request, user.ID, user.IsSuperuser)
	if err != nil {
		return shared.Fail(c, err)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.houseService.Delete(c.Context(), slug, user.ID, user.IsSuperuser); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "house deleted")
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
		return shared.Fail(c, err)
	}

	return c.JSON(schema.SlugCheckResponse{Available: available, Slug: normalized})
}
