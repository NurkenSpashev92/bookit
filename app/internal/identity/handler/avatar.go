package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/identity/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type AvatarHandler struct {
	avatarService *service.AvatarService
}

func NewAvatarHandler(avatarService *service.AvatarService) *AvatarHandler {
	return &AvatarHandler{avatarService: avatarService}
}

// Upload godoc
// @Summary Upload user avatar
// @Description Upload or replace the authenticated user's avatar
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image"
// @Success 200 {object} shared.MessageResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [post]
func (h *AvatarHandler) Upload(c fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{Error: "unauthenticated"})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{Error: "avatar file is required"})
	}

	if _, err := h.avatarService.Upload(c.Context(), user.ID, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "avatar uploaded"})
}

// Delete godoc
// @Summary Delete user avatar
// @Description Remove the authenticated user's avatar
// @Tags Auth
// @Produce json
// @Success 200 {object} shared.MessageResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [delete]
func (h *AvatarHandler) Delete(c fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{Error: "unauthenticated"})
	}

	if err := h.avatarService.Delete(c.Context(), user.ID); err != nil {
		if errors.Is(err, service.ErrAvatarNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "avatar deleted"})
}
