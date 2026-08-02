package handler

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type AvatarService interface {
	Upload(ctx context.Context, userID int, file *multipart.FileHeader) (schema.AuthUser, error)
	Delete(ctx context.Context, userID int) (schema.AuthUser, error)
	UploadQR(ctx context.Context, userID int, file *multipart.FileHeader) (schema.AuthUser, error)
	DeleteQR(ctx context.Context, userID int) (schema.AuthUser, error)
}

type AvatarHandler struct {
	avatarService AvatarService
}

func NewAvatarHandler(avatarService AvatarService) *AvatarHandler {
	return &AvatarHandler{avatarService: avatarService}
}

// Upload godoc
// @Summary Upload user avatar
// @Description Upload or replace the authenticated user's avatar
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image"
// @Success 200 {object} schema.AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [post]
func (h *AvatarHandler) Upload(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return shared.Fail(c, shared.Invalid("avatar file is required"))
	}

	authUser, err := h.avatarService.Upload(c.Context(), user.ID, file)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(schema.AuthResponse{User: authUser})
}

// Delete godoc
// @Summary Delete user avatar
// @Description Remove the authenticated user's avatar
// @Tags Auth
// @Produce json
// @Success 200 {object} schema.AuthResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/avatar [delete]
func (h *AvatarHandler) Delete(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	authUser, err := h.avatarService.Delete(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(schema.AuthResponse{User: authUser})
}

// UploadQR godoc
// @Summary Upload payment QR
// @Description Upload or replace the authenticated user's payment QR image
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param qr formData file true "Payment QR image"
// @Success 200 {object} schema.AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/payment-qr [post]
func (h *AvatarHandler) UploadQR(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	file, err := c.FormFile("qr")
	if err != nil {
		return shared.Fail(c, shared.Invalid("qr file is required"))
	}

	authUser, err := h.avatarService.UploadQR(c.Context(), user.ID, file)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(schema.AuthResponse{User: authUser})
}

// DeleteQR godoc
// @Summary Delete payment QR
// @Description Remove the authenticated user's payment QR image
// @Tags Auth
// @Produce json
// @Success 200 {object} schema.AuthResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/payment-qr [delete]
func (h *AvatarHandler) DeleteQR(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	authUser, err := h.avatarService.DeleteQR(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(schema.AuthResponse{User: authUser})
}
