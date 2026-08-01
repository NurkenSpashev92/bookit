package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func AuthAdmin(validator TokenValidator) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := accessToken(c)
		if token == "" {
			return shared.Fail(c, shared.Unauthorized("unauthenticated"))
		}

		user, err := validator.ValidateToken(token)
		if err != nil {
			return shared.Fail(c, shared.Unauthorized("invalid or expired token"))
		}

		if !user.IsSuperuser {
			return shared.Fail(c, shared.Forbidden("forbidden"))
		}

		c.Locals(userLocalsKey, user)

		return c.Next()
	}
}
