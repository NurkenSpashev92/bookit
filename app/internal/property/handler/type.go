package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type TypeHandler struct {
	typeService *service.TypeService
}

func NewTypeHandler(typeService *service.TypeService) *TypeHandler {
	return &TypeHandler{typeService: typeService}
}

// GetAll godoc
// @Summary Get all types
// @Tags Types
// @Produce json
// @Success 200 {array} schema.TypeResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /types [get]
func (h *TypeHandler) GetAll(c fiber.Ctx) error {
	types, err := h.typeService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(types)
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
	id, _ := strconv.Atoi(c.Params("id"))

	t, err := h.typeService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: "type not found"})
	}

	return c.JSON(t)
}

// Create godoc
// @Summary Create a type
// @Tags Types
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name"
// @Param is_active formData bool false "Is Active"
// @Param icon formData file false "Icon"
// @Success 201 {object} schema.TypeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /types [post]
func (h *TypeHandler) Create(c fiber.Ctx) error {
	name := c.FormValue("name")
	isActiveStr := c.FormValue("is_active")

	createReq := schema.TypeCreateRequest{Name: name}
	if err := createReq.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	isActive := true
	if isActiveStr != "" {
		isActive = isActiveStr == "true"
	}

	file, _ := c.FormFile("icon")

	created, err := h.typeService.Create(c.Context(), name, isActive, file)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(created)
}

// Update godoc
// @Summary Update a type
// @Tags Types
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Type ID"
// @Param name formData string false "Name"
// @Param is_active formData bool false "Is Active"
// @Param icon formData file false "Icon"
// @Success 200 {object} schema.TypeResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /types/{id} [patch]
func (h *TypeHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var name *string
	if n := c.FormValue("name"); n != "" {
		name = &n
	}

	var isActive *bool
	if ia := c.FormValue("is_active"); ia != "" {
		v := ia == "true"
		isActive = &v
	}

	file, _ := c.FormFile("icon")

	updated, err := h.typeService.Update(c.Context(), id, name, isActive, file)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
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
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.typeService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "type deleted"})
}
