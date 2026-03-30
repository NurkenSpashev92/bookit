package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func CorsHandler(c fiber.Ctx) error {
	origin := c.Get("Origin")
	if origin == "" {
		return c.Next()
	}

	allowedOrigins := getAllowedOrigins()

	allowed := false
	for _, o := range allowedOrigins {
		if o == origin || o == "*" {
			allowed = true
			break
		}
	}

	if !allowed {
		return c.Next()
	}

	c.Set("Access-Control-Allow-Origin", origin)
	c.Set("Access-Control-Allow-Credentials", "true")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	c.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Set("Access-Control-Max-Age", "86400")

	if c.Method() == fiber.MethodOptions {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Next()
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:3000", "http://localhost:5173"}
	}
	return strings.Split(origins, ",")
}
