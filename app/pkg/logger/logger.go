package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

type Config struct {
	Level     string
	Format    string
	AddSource bool
}

func New(cfg Config, w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     ParseLevel(cfg.Level),
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if strings.EqualFold(cfg.Format, FormatText) {
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}

	return slog.New(handler)
}

func Init(cfg Config) *slog.Logger {
	log := New(cfg, os.Stdout)
	slog.SetDefault(log)

	return log
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
