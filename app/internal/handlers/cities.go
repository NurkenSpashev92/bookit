package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

// GetCities godoc
// @Summary Get all cities
// @Tags Cities
// @Produce json
// @Success 200 {array} schemas.City
// @Failure 500 {object} schemas.ErrorResponse
// @Router /cities [get]
func GetCities(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewCityRepository(db)
		cities, err := repo.GetAllWithCountry(c.Context())
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(cities)
	}
}

// GetCity godoc
// @Summary Get city by ID
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} schemas.City
// @Failure 404 {object} schemas.ErrorResponse
// @Router /cities/{id} [get]
func GetCity(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewCityRepository(db)
		city, err := repo.GetByIDWithCountry(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(city)
	}
}

// CreateCity godoc
// @Summary Create a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param city body schemas.CityCreateRequest true "City data"
// @Success 201 {object} schemas.City
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /cities [post]
func CreateCity(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req schemas.CityCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		if req.NameKZ == "" || req.NameEN == "" || req.NameRU == "" || req.CountryID == 0 {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "all fields are required"})
		}

		repo := repositories.NewCityRepository(db)
		city, err := repo.Create(c.Context(), req)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(201).JSON(city)
	}
}

// UpdateCity godoc
// @Summary Update a city
// @Tags Cities
// @Accept json
// @Produce json
// @Param id path int true "City ID"
// @Param city body schemas.CityUpdateRequest true "City data"
// @Success 200 {object} schemas.City
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /cities/{id} [patch]
func UpdateCity(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		var req schemas.CityUpdateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewCityRepository(db)
		city, err := repo.Update(c.Context(), id, req)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(city)
	}
}

// DeleteCity godoc
// @Summary Delete a city
// @Tags Cities
// @Produce json
// @Param id path int true "City ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /cities/{id} [delete]
func DeleteCity(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewCityRepository(db)
		if err := repo.Delete(c.Context(), id); err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(schemas.MessageResponse{Message: "city deleted"})
	}
}
