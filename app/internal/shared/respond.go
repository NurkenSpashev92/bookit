package shared

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func Fail(c fiber.Ctx, status int, err error) error {
	if status >= http.StatusInternalServerError {
		log.Printf("[%s %s] %d: %v", c.Method(), c.Path(), status, err)
		return c.Status(status).JSON(ErrorResponse{Error: internalErrorMessage})
	}

	return c.Status(status).JSON(ErrorResponse{Error: err.Error()})
}

func FailMsg(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{Error: message})
}

const internalErrorMessage = "internal server error"
