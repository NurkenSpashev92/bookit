package initializers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ServerHeader:        "Bookit",
		AppName:             "Bookit App v0.1-beta",
		CaseSensitive:       true,
		ErrorHandler:        shared.ErrorHandler,
		PassLocalsToContext: true,
		BodyLimit:           80 * 1024 * 1024,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        15 * time.Second,
		IdleTimeout:         60 * time.Second,
	}
}
