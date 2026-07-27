package healthcheck

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
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
	return healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			return db.Ping(c.Context()) == nil
		},
		ResponseFormat: healthcheck.FormatJSON,
	})
}
