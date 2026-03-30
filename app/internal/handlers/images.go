package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type ImageHandler struct {
	imageService *services.ImageService
}

func NewImageHandler(imageService *services.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// UploadHouseImages godoc
// @Summary Upload house images
// @Description Upload multiple images for a house (max 15 images per house)
// @Tags Houses
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "House ID"
// @Param files formData []file true "Images"
// @Success 200 {object} schemas.MessageResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{id}/images [post]
func (h *ImageHandler) Upload(c fiber.Ctx) error {
	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid house id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "no files"})
	}

	if err := h.imageService.UploadHouseImages(c.Context(), houseID, files); err != nil {
		if errors.Is(err, services.ErrMaxImagesExceeded) {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "images uploaded"})
}

// DeleteHouseImage godoc
// @Summary Delete house image
// @Description Deletes image from database and AWS S3
// @Tags Houses
// @Produce json
// @Param image_id path int true "Image ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/images/{image_id} [delete]
func (h *ImageHandler) Delete(c fiber.Ctx) error {
	imageID, err := strconv.Atoi(c.Params("image_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{Error: "invalid image id"})
	}

	if err := h.imageService.DeleteHouseImage(c.Context(), imageID); err != nil {
		if errors.Is(err, services.ErrImageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "image deleted"})
}
