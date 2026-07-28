package configs

import (
	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func NewLogConfig() logger.Config {
	return logger.Config{
		Level: getEnv("LOG_LEVEL", "info"),
		File:  getEnv("LOG_FILE", ""),
	}
}
