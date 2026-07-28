package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func AuthOptional(validator TokenValidator) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := accessToken(c)
		if token == "" {
			return c.Next()
		}

		if user, err := validator.ValidateToken(token); err == nil {
			c.Locals(userLocalsKey, user)
		}

		return c.Next()
	}
}
