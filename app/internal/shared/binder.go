package shared

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Validatable interface {
	Validate() error
}

func Bind[T any](c fiber.Ctx) (T, error) {
	var request T

	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return request, Invalid("invalid body")
	}

	if validatable, ok := any(&request).(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return request, err
		}
	}

	return request, nil
}

func ParamInt(c fiber.Ctx, name string) (int, error) {
	value, err := strconv.Atoi(c.Params(name))
	if err != nil {
		return 0, Invalid("invalid " + name)
	}

	return value, nil
}

func ParamString(c fiber.Ctx, name string) (string, error) {
	value := c.Params(name)
	if value == "" {
		return "", Invalid(name + " is required")
	}

	return value, nil
}

func QueryIntInRange(c fiber.Ctx, name string, fallback, min, max int) int {
	value, err := strconv.Atoi(c.Query(name))
	if err != nil || value < min || value > max {
		return fallback
	}

	return value
}
