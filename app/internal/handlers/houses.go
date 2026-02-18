package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

// GetHouses godoc
// @Summary      Get all houses
// @Description  Returns a list of all houses
// @Tags         Houses
// @Produce      json
// @Success      200  {array}  schemas.HouseListItem
// @Failure      500  {object} schemas.ErrorResponse
// @Router       /houses [get]
func GetHouses(db *pgxpool.Pool, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewHouseRepository(db)

		houses, err := repo.GetAll(c.Context())
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		utils.FillHouseListImagesURL(cfg, houses)

		return c.JSON(houses)
	}
}

// GetHouseByID godoc
// @Summary      Get house by ID
// @Description  Returns a single house by ID
// @Tags         Houses
// @Produce      json
// @Param        id   path      int  true  "House ID"
// @Success      200  {object} models.House
// @Failure      404  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Router       /houses/{id} [get]
func GetHouseByID(db *pgxpool.Pool, cfg *configs.AwsConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		repo := repositories.NewHouseRepository(db)
		house, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "house not found: " + err.Error()})
		}
		utils.FillHouseImagesURL(cfg, house.Images)
		return c.JSON(house)
	}
}

// CreateHouse godoc
// @Summary      Create a new house
// @Description  Creates a new house. Auth required.
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        house  body  schemas.HouseCreateRequest  true  "House data"
// @Success      201    {object} models.House
// @Failure      400    {object} schemas.ErrorResponse
// @Failure      401    {object} schemas.ErrorResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /houses [post]
func CreateHouse(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(models.User)

		var req schemas.HouseCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		req.OwnerID = user.ID
		repo := repositories.NewHouseRepository(db)
		house, err := repo.Create(c.Context(), req)
		if err != nil {
			if strings.Contains(err.Error(), "slug already exists") {
				return c.Status(409).JSON(fiber.Map{
					"error": "slug already exists",
				})
			}

			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(201).JSON(house)
	}
}

// UpdateHouse godoc
// @Summary      Update house
// @Description  Updates house data. Auth required.
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "House ID"
// @Param        house  body      schemas.HouseUpdateRequest true  "House update data"
// @Success      200    {object} models.House
// @Failure      400    {object} schemas.ErrorResponse
// @Failure      401    {object} schemas.ErrorResponse
// @Failure      500    {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router       /houses/{id} [patch]
func UpdateHouse(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))

		var req schemas.HouseUpdateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewHouseRepository(db)
		house, err := repo.Update(c.Context(), id, req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
			}
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(200).JSON(house)
	}
}

// DeleteHouse godoc
// @Summary      Delete house
// @Description  Deletes a house by ID. Auth required.
// @Tags         Houses
// @Produce      json
// @Param        id   path      int  true  "House ID"
// @Success      200  {object} schemas.MessageResponse
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router       /houses/{id} [delete]
func DeleteHouse(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))

		repo := repositories.NewHouseRepository(db)
		if err := repo.Delete(c.Context(), id); err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(fiber.Map{"message": "house deleted"})
	}
}

// CheckSlug godoc
// @Summary      Check house slug availability
// @Description  Checks if a house slug is available for use
// @Tags         Houses
// @Accept       json
// @Produce      json
// @Param        slug   query     string  true  "Slug to check"
// @Success      200    {object} schemas.SlugCheckResponse
// @Failure      400    {object}  schemas.ErrorResponse
// @Failure      500    {object}  schemas.ErrorResponse
// @Router       /houses/check-slug [get]
func CheckSlug(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		s := slug.Make(c.Query("slug"))

		repo := repositories.NewHouseRepository(db)
		exists, err := repo.SlugExists(c.Context(), s)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"available": !exists,
			"slug":      s,
		})
	}
}
