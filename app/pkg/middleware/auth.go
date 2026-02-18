package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/services"
)

func AuthRequired(jwtService *services.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("jwt")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthenticated",
			})
		}

		user, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
