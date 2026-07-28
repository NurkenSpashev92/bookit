package middleware

import (
	"github.com/gofiber/fiber/v3"

	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

const userLocalsKey = "user"

func CurrentUser(c fiber.Ctx) (identitymodel.User, error) {
	user, ok := c.Locals(userLocalsKey).(identitymodel.User)
	if !ok {
		return identitymodel.User{}, shared.Unauthorized("unauthenticated")
	}

	return user, nil
}

func CurrentUserID(c fiber.Ctx) int {
	user, ok := c.Locals(userLocalsKey).(identitymodel.User)
	if !ok {
		return 0
	}

	return user.ID
}
