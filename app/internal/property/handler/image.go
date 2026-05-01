package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type ImageHandler struct {
	imageService *service.ImageService
}

func NewImageHandler(imageService *service.ImageService) *ImageHandler {
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
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "invalid form"})
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		files = form.File["files"]
	}
	if len(files) == 0 {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "no files"})
	}

	if err := h.imageService.UploadHouseImages(c.Context(), slug, files); err != nil {
		if errors.Is(err, service.ErrMaxImagesExceeded) {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "images uploaded"})
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
	imageID, err := strconv.Atoi(c.Params("image_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{Error: "invalid image id"})
	}

	if err := h.imageService.DeleteHouseImage(c.Context(), imageID); err != nil {
		if errors.Is(err, service.ErrImageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "image deleted"})
}
