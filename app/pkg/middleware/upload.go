package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func UploadLimits(maxBody int, timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		if len(c.Body()) > maxBody {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "file too large",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()
		c.SetContext(ctx)

		return c.Next()
	}
}
