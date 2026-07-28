package middleware

import (
	"github.com/gofiber/fiber/v3"

	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type TokenValidator interface {
	ValidateToken(token string) (identitymodel.User, error)
}

func AuthRequired(validator TokenValidator) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := accessToken(c)
		if token == "" {
			return shared.Fail(c, shared.Unauthorized("unauthenticated"))
		}

		user, err := validator.ValidateToken(token)
		if err != nil {
			return shared.Fail(c, shared.Unauthorized("invalid or expired token"))
		}

		c.Locals(userLocalsKey, user)

		return c.Next()
	}
}

func accessToken(c fiber.Ctx) string {
	if token := c.Cookies("access_token"); token != "" {
		return token
	}

	return c.Cookies("jwt")
}
