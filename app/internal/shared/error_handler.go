package shared

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
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
		fiberlog.WithContext(c.Context()).Errorw("unhandled error",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"error", err,
		)
	}

	return c.Status(status).JSON(ErrorResponse{Error: message})
}
