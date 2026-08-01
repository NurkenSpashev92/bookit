package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type TypeService interface {
	GetAll(ctx context.Context) ([]schema.TypeResponse, error)
	GetAllPaginated(ctx context.Context, limit, offset int) ([]schema.TypeResponse, int, error)
	GetByID(ctx context.Context, id int) (schema.TypeResponse, error)
	Create(ctx context.Context, req schema.TypeCreateRequest) (schema.TypeResponse, error)
	Update(ctx context.Context, id int, req schema.TypeUpdateRequest) (schema.TypeResponse, error)
	Delete(ctx context.Context, id int) error
}

type TypeHandler struct {
	typeService TypeService
}

func NewTypeHandler(typeService TypeService) *TypeHandler {
	return &TypeHandler{typeService: typeService}
}

// GetAll godoc
// @Summary Get all types
// @Description Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags Types
// @Produce json
// @Param page query int false "Page number (enables the paginated envelope)"
// @Param page_size query int false "Items per page (default 20, max 100)"
// @Success 200 {array} schema.TypeResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /types [get]
func (h *TypeHandler) GetAll(c fiber.Ctx) error {
	return shared.ListMaybePaginated(c, h.typeService.GetAll, h.typeService.GetAllPaginated)
}

// GetByID godoc
// @Summary Get type by ID
// @Tags Types
// @Produce json
// @Param id path int true "Type ID"
// @Success 200 {object} schema.TypeResponse
// @Failure 404 {object} shared.ErrorResponse
// @Router /types/{id} [get]
func (h *TypeHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	propertyType, err := h.typeService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(propertyType)
}

// Create godoc
// @Summary Create a type
// @Tags Types
// @Accept json
// @Produce json
// @Param request body schema.TypeCreateRequest true "Type"
// @Success 201 {object} schema.TypeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /types [post]
func (h *TypeHandler) Create(c fiber.Ctx) error {
	request, err := shared.Bind[schema.TypeCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	created, err := h.typeService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, created)
}

// Update godoc
// @Summary Update a type
// @Tags Types
// @Accept json
// @Produce json
// @Param id path int true "Type ID"
// @Param request body schema.TypeUpdateRequest true "Fields to update"
// @Success 200 {object} schema.TypeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /types/{id} [patch]
func (h *TypeHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.TypeUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	updated, err := h.typeService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(updated)
}

// Delete godoc
// @Summary Delete a type
// @Tags Types
// @Produce json
// @Param id path int true "Type ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /types/{id} [delete]
func (h *TypeHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.typeService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "type deleted")
}
