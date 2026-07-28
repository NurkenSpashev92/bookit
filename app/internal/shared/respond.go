package shared

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
)

const internalErrorMessage = "internal server error"

func Fail(c fiber.Ctx, err error) error {
	status := StatusOf(err)

	if status >= http.StatusInternalServerError {
		fiberlog.WithContext(c.Context()).Errorw("request failed",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"error", err,
		)

		return c.Status(status).JSON(ErrorResponse{Error: internalErrorMessage})
	}

	return c.Status(status).JSON(ErrorResponse{Error: MessageOf(err)})
}

func OK(c fiber.Ctx, message string) error {
	return c.JSON(MessageResponse{Message: message})
}

func Created(c fiber.Ctx, payload any) error {
	return c.Status(http.StatusCreated).JSON(payload)
}

func Items[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}

func List[T any](c fiber.Ctx, items []T) error {
	return c.JSON(Items(items))
}
