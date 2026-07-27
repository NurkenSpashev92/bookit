package shared

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
)

func Fail(c fiber.Ctx, status int, err error) error {
	if status >= http.StatusInternalServerError {
		fiberlog.WithContext(c.Context()).Errorw("request failed",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"error", err,
		)

		return c.Status(status).JSON(ErrorResponse{Error: internalErrorMessage})
	}

	return c.Status(status).JSON(ErrorResponse{Error: err.Error()})
}

func FailMsg(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: message})
}

const internalErrorMessage = "internal server error"
