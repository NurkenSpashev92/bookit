package handlers

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

const maxHouseImages = 15

// UploadHouseImages godoc
// @Summary Upload house images
// @Description Upload multiple images for a house (max 15 images per house)
// @Tags Houses
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "House ID"
// @Param files formData []file true "Images"
// @Success 200 {object} schemas.MessageResponse
// @Failure 400 {object} schemas.ErrorResponse "Invalid request or max images exceeded"
// @Failure 401 {object} schemas.ErrorResponse "Unauthorized"
// @Failure 404 {object} schemas.ErrorResponse "House not found"
// @Failure 500 {object} schemas.ErrorResponse "Upload or database error"
// @Security ApiKeyAuth
// @Router /houses/{id}/images [post]
func UploadHouseImages(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
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

		repo := repositories.NewHouseImageRepository(db)

		count, err := repo.CountByHouse(c.Context(), houseID)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		if count+len(files) > maxHouseImages {
			return c.Status(400).JSON(schemas.ErrorResponse{
				Error: "maximum 15 images allowed",
			})
		}

		type uploaded struct {
			key      string
			thumbKey string
			size     int
			mime     string
			width    int
			height   int
		}

		var (
			wg      sync.WaitGroup
			mutex   sync.Mutex
			results []uploaded
			upErr   error
		)

		for _, file := range files {
			wg.Add(1)

			go func(f *multipart.FileHeader) {
				defer wg.Done()

				img, err := services.Process(f)
				if err != nil {
					upErr = err
					return
				}

				originalKey := fmt.Sprintf("houses/original/%d_%d.jpg", houseID, time.Now().UnixNano())
				thumbKey := fmt.Sprintf("houses/thumbnail/%d_%d.webp", houseID, time.Now().UnixNano())

				_, err = s3.UploadCompressed(c.Context(), originalKey, img.Original, "image/jpeg")
				if err != nil {
					upErr = err
					return
				}

				_, err = s3.UploadCompressed(c.Context(), thumbKey, img.Thumbnail, "image/webp")
				if err != nil {
					upErr = err
					return
				}

				mutex.Lock()
				results = append(results, uploaded{
					key:      originalKey,
					thumbKey: thumbKey,
					size:     img.Size,
					mime:     img.Mime,
					width:    img.Width,
					height:   img.Height,
				})
				mutex.Unlock()

			}(file)
		}

		wg.Wait()

		if upErr != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: upErr.Error()})
		}

		for _, r := range results {
			img := models.Image{
				Original:  r.key,
				Thumbnail: r.thumbKey,
				MimeType:  r.mime,
				Width:     &r.width,
				Height:    &r.height,
				Size:      &r.size,
				HouseID:   &houseID,
			}

			if err := repo.Create(c.Context(), &img); err != nil {
				for _, r := range results {
					_ = s3.Delete(c.Context(), r.key)
				}
				return c.Status(500).JSON(schemas.ErrorResponse{Error: "db save failed"})
			}
		}

		return c.JSON(fiber.Map{"message": "images uploaded"})
	}
}

// DeleteHouseImage godoc
// @Summary Delete house image
// @Description Deletes image from database and AWS S3
// @Tags Houses
// @Produce json
// @Param image_id path int true "Image ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 400 {object} schemas.ErrorResponse "Invalid ID"
// @Failure 401 {object} schemas.ErrorResponse "Unauthorized"
// @Failure 404 {object} schemas.ErrorResponse "Image not found"
// @Failure 500 {object} schemas.ErrorResponse "Delete error"
// @Security ApiKeyAuth
// @Router /houses/images/{image_id} [delete]
func DeleteHouseImage(db *pgxpool.Pool, s3 *aws.AwsS3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		imageID, err := strconv.Atoi(c.Params("image_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(schemas.ErrorResponse{Error: "invalid image id"})
		}

		repo := repositories.NewHouseImageRepository(db)

		key, err := repo.Delete(c.Context(), imageID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).
				JSON(schemas.ErrorResponse{Error: "image not found"})
		}

		if key != nil && *key != "" {
			_ = s3.Delete(c.Context(), *key)
		}

		return c.JSON(schemas.MessageResponse{
			Message: "image deleted",
		})
	}
}
