package main

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/cmd/apiserver"
	"github.com/nurkenspashev92/bookit/internal/initializers"
)

// @title Bookit API
// @version 1.0
// @description Bookit API Documentation
// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := &apiserver.ApiApp{
		App: fiber.New(initializers.NewFiberConfig()),
	}
	app.Run()
}
