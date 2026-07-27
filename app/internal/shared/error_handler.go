package shared

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// ErrorHandler renders every error escaping a handler or middleware as a JSON
// ErrorResponse: client errors keep their message, server errors are logged
// with details and answered with a neutral message.
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
		log.Printf("[%s %s] %d: %v", c.Method(), c.Path(), status, err)
	}

	return c.Status(status).JSON(ErrorResponse{Error: message})
}
