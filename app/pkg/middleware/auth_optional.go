package middleware

import (
	"github.com/gofiber/fiber/v3"

	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
)

func AuthOptional(jwtService *identitysvc.JWTService) fiber.Handler {
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
