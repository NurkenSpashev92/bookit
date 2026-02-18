package middleware

import "github.com/gofiber/fiber/v3"

func CorsHandler(c fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	c.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

	return c.Next()
}
