package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func UploadLimits(maxBody int) fiber.Handler {
	return func(c fiber.Ctx) error {
		if len(c.Body()) > maxBody {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "file too large",
			})
		}

		return c.Next()
	}
}
