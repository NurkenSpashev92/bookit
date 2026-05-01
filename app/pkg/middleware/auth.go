package middleware

import (
	"github.com/gofiber/fiber/v3"

	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
)

func AuthRequired(jwtService *identitysvc.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			// Fallback: check old "jwt" cookie for backward compat
			token = c.Cookies("jwt")
		}
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthenticated",
			})
		}

		user, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
