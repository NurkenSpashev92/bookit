package handlers

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

// GetCategories godoc
// @Summary      Get all active categories
// @Description  Returns a list of active categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   schemas.CategoryPaginate
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /categories [get]
func GetCategories(db *pgxpool.Pool, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewCategoryRepository(db)

		categories, err := repo.GetCategories(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		var wg sync.WaitGroup
		for i := range categories {
			if categories[i].Icon != nil && *categories[i].Icon != "" {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					full := cfg.AwsS3URL(*categories[i].Icon)
					categories[i].Icon = &full
				}(i)
			}
		}
		wg.Wait()

		return c.JSON(categories)
	}
}

// GetCategory godoc
// @Summary Get category
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 404 {object} schemas.ErrorResponse
// @Router /categories/{id} [get]
func GetCategory(db *pgxpool.Pool, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{
				Error: "invalid id: " + err.Error(),
			})
		}

		repo := repositories.NewCategoryRepository(db)

		category, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "category not found: " + err.Error()})
		}

		if category.Icon != nil && *category.Icon != "" {
			full := cfg.AwsS3URL(*category.Icon)
			category.Icon = &full
		}

		return c.JSON(category)
	}
}

// CreateCategory godoc
// @Summary      Create category
// @Description  Creates a new category
// @Tags         Categories
// @Accept       multipart/form-data
// @Produce      json
// @Param        name_kz   formData string true "Name KZ"
// @Param        name_ru   formData string true "Name RU"
// @Param        name_en   formData string true "Name EN"
// @Param        is_active formData bool   false "Is active"
// @Param        icon   formData  file false "Category icon"
// @Success      201   {object}  models.Category
// @Failure      400   {object}  schemas.ErrorResponse
// @Failure      500   {object}  schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router       /categories [post]
func CreateCategory(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		nameKz := c.FormValue("name_kz")
		nameRu := c.FormValue("name_ru")
		nameEn := c.FormValue("name_en")
		isActiveStr := c.FormValue("is_active")

		if nameKz == "" || nameRu == "" || nameEn == "" {
			return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{
				Error: "all names required",
			})
		}

		isActive := true
		if isActiveStr != "" {
			isActive = isActiveStr == "true"
		}

		var iconURL *string
		file, err := c.FormFile("icon")
		if err == nil && file != nil {
			uploaded, err := s3.Upload(c.Context(), "categories", file)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.ErrorResponse{
					Error: "failed upload image " + err.Error(),
				})
			}
			full := cfg.AwsS3URL(uploaded)
			iconURL = &full
		}

		req := schemas.CategoryCreateRequest{
			NameKz:   nameKz,
			NameRu:   nameRu,
			NameEn:   nameEn,
			Icon:     iconURL,
			IsActive: isActive,
		}

		repo := repositories.NewCategoryRepository(db)
		category, err := repo.CreateCategory(c.Context(), req)
		if err != nil {
			if iconURL != nil {
				_ = s3.Delete(c.Context(), *iconURL)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(schemas.ErrorResponse{
				Error: "failed to create category: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(category)
	}
}

// UpdateCategory godoc
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
// @Success 200 {object} models.Category
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /categories/{id} [patch]
func UpdateCategory(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{
				Error: "invalid id: " + err.Error(),
			})
		}

		var req schemas.CategoryUpdateRequest
		if err := c.Bind().Form(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid form"})
		}

		var newIcon *string
		file, err := c.FormFile("icon")
		if err == nil && file != nil {
			uploaded, err := s3.Upload(c.Context(), "categories", file)
			if err != nil {
				return c.Status(500).JSON(schemas.ErrorResponse{Error: "upload failed: " + err.Error()})
			}
			full := cfg.AwsS3URL(uploaded)
			newIcon = &full
		}

		repo := repositories.NewCategoryRepository(db)
		category, oldIcon, err := repo.Update(c.Context(), id, req, newIcon)
		if err != nil {
			if newIcon != nil {
				_ = s3.Delete(c.Context(), *newIcon)
			}
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "category not found: " + err.Error()})
		}

		if newIcon != nil && oldIcon != nil {
			_ = s3.Delete(c.Context(), *oldIcon)
		}

		return c.JSON(category)
	}
}

// DeleteCategory godoc
// @Summary Delete category
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /categories/{id} [delete]
func DeleteCategory(db *pgxpool.Pool, s3 *aws.AwsS3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id: " + err.Error()})
		}

		repo := repositories.NewCategoryRepository(db)
		icon, err := repo.Delete(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{
				Error: "category not found: " + err.Error(),
			})
		}

		if icon != nil {
			_ = s3.Delete(c.Context(), *icon)
		}

		return c.JSON(schemas.MessageResponse{
			Message: "category deleted",
		})
	}
}
