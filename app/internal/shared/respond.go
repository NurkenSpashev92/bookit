package shared

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func Fail(c fiber.Ctx, status int, err error) error {
	if status >= http.StatusInternalServerError {
		logger.FromContext(c.Context()).ErrorContext(c.Context(), "request failed",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			logger.Err(err),
		)

		return c.Status(status).JSON(ErrorResponse{Error: internalErrorMessage})
	}

	return c.Status(status).JSON(ErrorResponse{Error: err.Error()})
}

func FailMsg(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: message})
}

const internalErrorMessage = "internal server error"
