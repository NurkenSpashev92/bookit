package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/services"
)

func AuthOptional(jwtService *services.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			token = c.Cookies("jwt")
		}
		if token != "" {
			if user, err := jwtService.ValidateToken(token); err == nil {
				c.Locals("user", user)
			}
		}
		return c.Next()
	}
}
