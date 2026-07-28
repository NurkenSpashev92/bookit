package handler

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

type ImageService interface {
	UploadHouseImages(ctx context.Context, slug string, files []*multipart.FileHeader) error
	DeleteHouseImage(ctx context.Context, imageID int) error
}

type ImageHandler struct {
	imageService ImageService
}

func NewImageHandler(imageService ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// Upload godoc
// @Summary Upload house images
// @Tags Houses
// @Accept multipart/form-data
// @Produce json
// @Param slug path string true "House slug"
// @Param files formData []file true "Images"
// @Success 200 {object} shared.MessageResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{slug}/images [post]
func (h *ImageHandler) Upload(c fiber.Ctx) error {
	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return shared.Fail(c, err)
	}

	files, err := uploadedFiles(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.imageService.UploadHouseImages(c.Context(), slug, files); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "images uploaded")
}

// Delete godoc
// @Summary Delete house image
// @Tags Houses
// @Produce json
// @Param image_id path int true "Image ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/images/{image_id} [delete]
func (h *ImageHandler) Delete(c fiber.Ctx) error {
	imageID, err := shared.ParamInt(c, "image_id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.imageService.DeleteHouseImage(c.Context(), imageID); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "image deleted")
}

func uploadedFiles(c fiber.Ctx) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, shared.Invalid("invalid form")
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		files = form.File["files"]
	}
	if len(files) == 0 {
		return nil, shared.Invalid("no files")
	}

	return files, nil
}
