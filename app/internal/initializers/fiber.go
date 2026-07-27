package initializers

import (
	"time"

	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ServerHeader:  "Bookit",
		AppName:       "Bookit App v0.1-beta",
		CaseSensitive: true,
		BodyLimit:     5 * 1024 * 1024, // 5 MB default
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  15 * time.Second,
		IdleTimeout:   60 * time.Second,
		// goccy/go-json — drop-in replacement, ~2-3× faster than encoding/json on
		// Marshal-heavy workloads like list endpoints.
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}
}
