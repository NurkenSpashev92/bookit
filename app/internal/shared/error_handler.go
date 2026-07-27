package shared

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	status := http.StatusInternalServerError
	message := internalErrorMessage

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code

		if status < http.StatusInternalServerError {
			message = fiberErr.Message
		}
	}

	if status >= http.StatusInternalServerError {
		logger.FromContext(c.Context()).ErrorContext(c.Context(), "unhandled error",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			logger.Err(err),
		)
	}

	return c.Status(status).JSON(ErrorResponse{Error: message})
}
