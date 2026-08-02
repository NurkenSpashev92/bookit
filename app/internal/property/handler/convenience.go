package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type ConvenienceService interface {
	GetAll(ctx context.Context, search string) ([]schema.ConveniencePaginate, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.ConveniencePaginate, int, error)
	GetByID(ctx context.Context, id int) (schema.Convenience, error)
	Create(ctx context.Context, req schema.ConvenienceCreateRequest) (schema.Convenience, error)
	Update(ctx context.Context, id int, req schema.ConvenienceUpdateRequest) (schema.Convenience, error)
	Delete(ctx context.Context, id int) error
}

type ConvenienceHandler struct {
	convenienceService ConvenienceService
}

func NewConvenienceHandler(convenienceService ConvenienceService) *ConvenienceHandler {
	return &ConvenienceHandler{convenienceService: convenienceService}
}

// GetAll godoc
// @Summary      Get all active conveniences
// @Description  Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags         Conveniences
// @Produce      json
// @Param        search     query string false "Search by name/slug"
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success      200  {array}   schema.ConveniencePaginate
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /conveniences [get]
func (h *ConvenienceHandler) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.ConveniencePaginate, error) {
			return h.convenienceService.GetAll(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.ConveniencePaginate, int, error) {
			return h.convenienceService.GetAllPaginated(ctx, search, limit, offset)
		},
	)
}

// GetByID godoc
// @Summary Get convenience
// @Tags Conveniences
// @Produce json
// @Param id path int true "Convenience ID"
// @Success 200 {object} schema.Convenience
// @Failure 404 {object} shared.ErrorResponse
// @Router /conveniences/{id} [get]
func (h *ConvenienceHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	convenience, err := h.convenienceService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(convenience)
}

// Create godoc
// @Summary      Create convenience
// @Tags         Conveniences
// @Accept       json
// @Produce      json
// @Param        request body schema.ConvenienceCreateRequest true "Convenience"
// @Success      201   {object}  schema.Convenience
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /conveniences [post]
func (h *ConvenienceHandler) Create(c fiber.Ctx) error {
	request, err := shared.Bind[schema.ConvenienceCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	convenience, err := h.convenienceService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, convenience)
}

// Update godoc
// @Summary Update convenience
// @Tags Conveniences
// @Accept json
// @Produce json
// @Param id path int true "Convenience ID"
// @Param request body schema.ConvenienceUpdateRequest true "Fields to update"
// @Success 200 {object} schema.Convenience
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /conveniences/{id} [patch]
func (h *ConvenienceHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.ConvenienceUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	convenience, err := h.convenienceService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(convenience)
}

// Delete godoc
// @Summary Delete convenience
// @Tags Conveniences
// @Produce json
// @Param id path int true "Convenience ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /conveniences/{id} [delete]
func (h *ConvenienceHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.convenienceService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "convenience deleted")
}
