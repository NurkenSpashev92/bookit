package middleware

import (
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

const corsMaxAge = 86400

func Cors() fiber.Handler {
	origins := allowedOrigins()

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowCredentials: !slices.Contains(origins, "*"),
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			fiber.HeaderAccept,
			fiber.HeaderContentType,
			fiber.HeaderAuthorization,
		},
		MaxAge: corsMaxAge,
	})
}

func allowedOrigins() []string {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:3000", "http://localhost:5173"}
	}

	parts := strings.Split(origins, ",")
	trimmed := make([]string, 0, len(parts))
	for _, origin := range parts {
		if origin = strings.TrimSpace(origin); origin != "" {
			trimmed = append(trimmed, origin)
		}
	}

	return trimmed
}
