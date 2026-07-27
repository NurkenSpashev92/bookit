package configs

import (
	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func NewLogConfig() logger.Config {
	return logger.Config{
		Level:     getEnv("LOG_LEVEL", "info"),
		Format:    getEnv("LOG_FORMAT", logger.FormatJSON),
		AddSource: getEnv("LOG_SOURCE", "false") == "true",
	}
}
