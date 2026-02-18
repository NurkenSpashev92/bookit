package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

// GetCountries godoc
// @Summary Get all countries
// @Tags Countries
// @Produce json
// @Success 200 {array} schemas.Country
// @Failure 500 {object} schemas.ErrorResponse
// @Router /countries [get]
func GetCountries(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewCountryRepository(db)
		countries, err := repo.GetAll(c.Context())
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(countries)
	}
}

// GetCountry godoc
// @Summary Get country by ID
// @Tags Countries
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} schemas.Country
// @Failure 404 {object} schemas.ErrorResponse
// @Router /countries/{id} [get]
func GetCountry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewCountryRepository(db)
		country, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "country not found"})
		}

		return c.JSON(country)
	}
}

// CreateCountry godoc
// @Summary Create a country
// @Tags Countries
// @Accept json
// @Produce json
// @Param country body schemas.CountryCreateRequest true "Country data"
// @Success 201 {object} schemas.Country
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /countries [post]
func CreateCountry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req schemas.CountryCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		if req.NameKZ == "" || req.NameEN == "" || req.NameRU == "" {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "all names are required"})
		}

		repo := repositories.NewCountryRepository(db)
		country, err := repo.Create(c.Context(), req)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(country)
	}
}

// UpdateCountry godoc
// @Summary Update a country
// @Tags Countries
// @Accept json
// @Produce json
// @Param id path int true "Country ID"
// @Param country body schemas.CountryUpdateRequest true "Country data"
// @Success 200 {object} schemas.Country
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /countries/{id} [patch]
func UpdateCountry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		var req schemas.CountryUpdateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewCountryRepository(db)
		country, err := repo.Update(c.Context(), id, req)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(country)
	}
}

// DeleteCountry godoc
// @Summary Delete a country
// @Tags Countries
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /countries/{id} [delete]
func DeleteCountry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewCountryRepository(db)
		if err := repo.Delete(c.Context(), id); err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(schemas.MessageResponse{Message: "country deleted"})
	}
}
