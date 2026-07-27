package main

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/cmd/apiserver"
	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/initializers"
	"github.com/nurkenspashev92/bookit/pkg/logger"
)

// @title Bookit API
// @version 1.0
// @description Bookit API Documentation
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger.Init(configs.NewLogConfig())

	app := &apiserver.ApiApp{
		App: fiber.New(initializers.NewFiberConfig()),
	}
	app.Run()
}
