package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CategoryService interface {
	GetAll(ctx context.Context, search string) ([]schema.CategoryPaginate, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.CategoryPaginate, int, error)
	GetByID(ctx context.Context, id int) (model.Category, error)
	Create(ctx context.Context, req schema.CategoryCreateRequest) (model.Category, error)
	Update(ctx context.Context, id int, req schema.CategoryUpdateRequest) (model.Category, error)
	Delete(ctx context.Context, id int) error
}

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// GetAll godoc
// @Summary      Get all active categories
// @Description  Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags         Categories
// @Produce      json
// @Param        search     query string false "Search by name/slug"
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success      200  {array}   schema.CategoryPaginate
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.CategoryPaginate, error) {
			return h.categoryService.GetAll(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.CategoryPaginate, int, error) {
			return h.categoryService.GetAllPaginated(ctx, search, limit, offset)
		},
	)
}

// GetByID godoc
// @Summary Get category
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.Category
// @Failure 404 {object} shared.ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	category, err := h.categoryService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(category)
}

// Create godoc
// @Summary      Create category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        request body schema.CategoryCreateRequest true "Category"
// @Success      201   {object}  model.Category
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /categories [post]
func (h *CategoryHandler) Create(c fiber.Ctx) error {
	request, err := shared.Bind[schema.CategoryCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	category, err := h.categoryService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, category)
}

// Update godoc
// @Summary Update category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body schema.CategoryUpdateRequest true "Fields to update"
// @Success 200 {object} model.Category
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /categories/{id} [patch]
func (h *CategoryHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.CategoryUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	category, err := h.categoryService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(category)
}

// Delete godoc
// @Summary Delete category
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.categoryService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "category deleted")
}
