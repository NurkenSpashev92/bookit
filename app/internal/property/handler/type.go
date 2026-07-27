package handler

import (
	"encoding/json"
	"net/http"
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
		return shared.Fail(c, http.StatusInternalServerError, err)
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
		return shared.FailMsg(c, http.StatusNotFound, "type not found")
	}

	return c.JSON(t)
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
	var req schema.TypeCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid body")
	}

	if err := req.Validate(); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	created, err := h.typeService.Create(c.Context(), req)
	if err != nil {
		return shared.Fail(c, http.StatusInternalServerError, err)
	}

	return c.Status(http.StatusCreated).JSON(created)
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
	id, _ := strconv.Atoi(c.Params("id"))

	var req schema.TypeUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return shared.FailMsg(c, http.StatusBadRequest, "invalid body")
	}

	if err := req.Validate(); err != nil {
		return shared.Fail(c, http.StatusBadRequest, err)
	}

	updated, err := h.typeService.Update(c.Context(), id, req)
	if err != nil {
		return shared.FailMsg(c, http.StatusNotFound, "type not found")
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
		return shared.Fail(c, http.StatusNotFound, err)
	}

	return c.JSON(shared.MessageResponse{Message: "type deleted"})
}
