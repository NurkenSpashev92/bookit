package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the application and database are running
// @Tags         HealthCheck
// @Produce      json
// @Success      200  {object}  interface{}
// @Failure      503  {object}  interface{}
// @Router       /healthcheck [get]
func HealthCheck(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "fail",
				"message": "database not reachable",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "success",
		})
	}
}
