package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func UploadLimits(maxBody int, timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Override body limit for this route
		if len(c.Body()) > maxBody {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "file too large",
			})
		}

		// Extended timeout via context deadline
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()
		c.SetContext(ctx)

		return c.Next()
	}
}
