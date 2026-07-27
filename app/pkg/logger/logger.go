package logger

import (
	"os"
	"strings"

	fiberlog "github.com/gofiber/fiber/v3/log"
)

type Config struct {
	Level string
}

func Init(cfg Config) {
	fiberlog.SetLevel(ParseLevel(cfg.Level))
	fiberlog.SetOutput(os.Stdout)
	fiberlog.MustSetContextTemplate(fiberlog.ContextConfig{Format: fiberlog.RequestIDFormat})
}

func ParseLevel(level string) fiberlog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return fiberlog.LevelTrace
	case "debug":
		return fiberlog.LevelDebug
	case "warn", "warning":
		return fiberlog.LevelWarn
	case "error":
		return fiberlog.LevelError
	case "fatal":
		return fiberlog.LevelFatal
	default:
		return fiberlog.LevelInfo
	}
}
