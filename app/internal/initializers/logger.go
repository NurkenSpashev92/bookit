package initializers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func NewLogger() fiber.Handler {
	return logger.New(logger.Config{
		Next:          nil,
		Done:          nil,
		Format:        "${time} | ${pid} | ${ip}:${port} | ${status}| ${method} | ${path} |${latency} | ${error}\n",
		TimeFormat:    "15:04:05",
		TimeZone:      "Local",
		TimeInterval:  500 * time.Millisecond,
		DisableColors: false,
	})
}
