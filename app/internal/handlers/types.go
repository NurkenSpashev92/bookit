package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

// GetTypes godoc
// @Summary Get all types
// @Tags Types
// @Produce json
// @Success 200 {array} schemas.TypeResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /types [get]
func GetTypes(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewTypeRepository(db)
		types, err := repo.GetAll(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		for i := range types {
			if types[i].Icon != "" {
				full := cfg.AwsS3URL(types[i].Icon)
				types[i].Icon = full
			}
		}

		return c.JSON(types)
	}
}

// GetTypeByID godoc
// @Summary Get type by ID
// @Tags Types
// @Produce json
// @Param id path int true "Type ID"
// @Success 200 {object} schemas.TypeResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Router /types/{id} [get]
func GetTypeByID(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		repo := repositories.NewTypeRepository(db)
		t, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "type not found"})
		}

		if t.Icon != "" {
			t.Icon = cfg.AwsS3URL(t.Icon)
		}

		return c.JSON(t)
	}
}

// CreateType godoc
// @Summary Create a type
// @Tags Types
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name"
// @Param is_active formData bool false "Is Active"
// @Param icon formData file false "Icon"
// @Success 201 {object} schemas.TypeResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /types [post]
func CreateType(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		name := c.FormValue("name")
		isActiveStr := c.FormValue("is_active")
		isActive := true
		if isActiveStr != "" {
			isActive = isActiveStr == "true"
		}

		var iconPath string
		file, err := c.FormFile("icon")
		if err == nil && file != nil {
			uploaded, err := s3.Upload(c.Context(), "types", file)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			iconPath = uploaded
		}

		t := models.Type{
			Name:     name,
			Icon:     iconPath,
			IsActive: isActive,
		}

		repo := repositories.NewTypeRepository(db)
		created, err := repo.Create(c.Context(), t)
		if err != nil {
			if iconPath != "" {
				_ = s3.Delete(c.Context(), iconPath)
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if created.Icon != "" {
			full := cfg.AwsS3URL(created.Icon)
			created.Icon = full
		}

		return c.Status(201).JSON(created)
	}
}

// UpdateType godoc
// @Summary Update a type
// @Tags Types
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Type ID"
// @Param name formData string false "Name"
// @Param is_active formData bool false "Is Active"
// @Param icon formData file false "Icon"
// @Success 200 {object} schemas.TypeResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /types/{id} [patch]
func UpdateType(db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		repo := repositories.NewTypeRepository(db)
		t, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "type not found"})
		}

		name := c.FormValue("name")
		if name != "" {
			t.Name = name
		}

		isActiveStr := c.FormValue("is_active")
		if isActiveStr != "" {
			t.IsActive = isActiveStr == "true"
		}

		file, err := c.FormFile("icon")
		if err == nil && file != nil {
			uploaded, err := s3.Upload(c.Context(), "types", file)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			if t.Icon != "" {
				_ = s3.Delete(c.Context(), t.Icon)
			}

			t.Icon = uploaded
		}

		updated, err := repo.Update(c.Context(), id, t)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if updated.Icon != "" {
			updated.Icon = cfg.AwsS3URL(updated.Icon)
		}

		return c.JSON(updated)
	}
}

// DeleteType godoc
// @Summary Delete a type
// @Tags Types
// @Produce json
// @Param id path int true "Type ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /types/{id} [delete]
func DeleteType(db *pgxpool.Pool, s3 *aws.AwsS3Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		repo := repositories.NewTypeRepository(db)
		icon, err := repo.Delete(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		if icon != "" {
			_ = s3.Delete(c.Context(), icon)
		}

		return c.JSON(fiber.Map{"message": "type deleted"})
	}
}
