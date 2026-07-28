package main

import (
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"

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
	logFile, err := logger.Init(configs.NewLogConfig())
	if err != nil {
		fiberlog.Errorw("failed to open log file, falling back to stdout", "error", err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	app := &apiserver.ApiApp{
		App: fiber.New(initializers.NewFiberConfig()),
	}
	app.Run()
}
