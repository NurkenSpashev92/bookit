package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type AvatarHandler struct {
	avatarService *services.AvatarService
}

func NewAvatarHandler(avatarService *services.AvatarService) *AvatarHandler {
	return &AvatarHandler{avatarService: avatarService}
}

// Upload godoc
// @Summary Upload user avatar
// @Description Upload or replace the authenticated user's avatar
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image"
// @Success 200 {object} schemas.MessageResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [post]
func (h *AvatarHandler) Upload(c fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.ErrorResponse{Error: "unauthenticated"})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{Error: "avatar file is required"})
	}

	if _, err := h.avatarService.Upload(c.Context(), user.ID, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "avatar uploaded"})
}

// Delete godoc
// @Summary Delete user avatar
// @Description Remove the authenticated user's avatar
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.MessageResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [delete]
func (h *AvatarHandler) Delete(c fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.ErrorResponse{Error: "unauthenticated"})
	}

	if err := h.avatarService.Delete(c.Context(), user.ID); err != nil {
		if errors.Is(err, services.ErrAvatarNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "avatar deleted"})
}
