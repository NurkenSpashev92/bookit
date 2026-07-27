package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// GetAll godoc
// @Summary      Get all active categories
// @Tags         Categories
// @Produce      json
// @Success      200  {array}   schema.CategoryPaginate
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	categories, err := h.categoryService.GetAll(c.Context())
	if err != nil {
		return shared.Fail(c, http.StatusInternalServerError, err)
	}
	return c.JSON(categories)
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{Error: "invalid id: " + err.Error()})
	}

	category, err := h.categoryService.GetByID(c.Context(), id)
	if err != nil {
		return shared.FailMsg(c, http.StatusNotFound, "category not found")
	}

	return c.JSON(category)
}

// Create godoc
// @Summary      Create category
// @Tags         Categories
// @Accept       multipart/form-data
// @Produce      json
// @Param        name_kz   formData string true "Name KZ"
// @Param        name_ru   formData string true "Name RU"
// @Param        name_en   formData string true "Name EN"
// @Param        is_active formData bool   false "Is active"
// @Param        icon      formData file   false "Category icon"
// @Success      201   {object}  model.Category
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /categories [post]
func (h *CategoryHandler) Create(c fiber.Ctx) error {
	nameKz := c.FormValue("name_kz")
	nameRu := c.FormValue("name_ru")
	nameEn := c.FormValue("name_en")
	isActiveStr := c.FormValue("is_active")

	createReq := schema.CategoryCreateRequest{NameKz: nameKz, NameRu: nameRu, NameEn: nameEn}
	if err := createReq.Validate(); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	isActive := true
	if isActiveStr != "" {
		isActive = isActiveStr == "true"
	}

	file, _ := c.FormFile("icon")

	category, err := h.categoryService.Create(c.Context(), nameKz, nameRu, nameEn, isActive, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{Error: "failed to create category: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// Update godoc
// @Summary Update category
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Category ID"
// @Param name_kz formData string false "Name KZ"
// @Param name_ru formData string false "Name RU"
// @Param name_en formData string false "Name EN"
// @Param is_active formData bool false "Is active"
// @Param icon formData file false "Icon"
// @Success 200 {object} model.Category
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /categories/{id} [patch]
func (h *CategoryHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{Error: "invalid id: " + err.Error()})
	}

	var req schema.CategoryUpdateRequest
	if err := c.Bind().Form(&req); err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid form")
	}

	file, _ := c.FormFile("icon")

	category, err := h.categoryService.Update(c.Context(), id, req, file)
	if err != nil {
		return shared.FailMsg(c, http.StatusNotFound, "category not found")
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid id")
	}

	if err := h.categoryService.Delete(c.Context(), id); err != nil {
		return shared.FailMsg(c, http.StatusNotFound, "category not found")
	}

	return c.JSON(shared.MessageResponse{Message: "category deleted"})
}
