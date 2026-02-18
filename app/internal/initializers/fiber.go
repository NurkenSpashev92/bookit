package initializers

import (
	"github.com/gofiber/fiber/v3"
)

func NewFiberConfig() fiber.Config {
	cfg := fiber.Config{
		ServerHeader:  "Bookit",
		AppName:       "Bookit App v0.1-beta",
		CaseSensitive: true,
	}

	return cfg
}
