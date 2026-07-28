package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	fiberlog "github.com/gofiber/fiber/v3/log"
)

const logFileMode = 0o644

type Config struct {
	Level string
	File  string
}

func Init(cfg Config) (io.Closer, error) {
	fiberlog.SetLevel(ParseLevel(cfg.Level))
	fiberlog.MustSetContextTemplate(fiberlog.ContextConfig{Format: fiberlog.RequestIDFormat})

	if cfg.File == "" {
		fiberlog.SetOutput(os.Stdout)
		return nil, nil
	}

	file, err := openLogFile(cfg.File)
	if err != nil {
		fiberlog.SetOutput(os.Stdout)
		return nil, err
	}

	fiberlog.SetOutput(io.MultiWriter(os.Stdout, file))

	return file, nil
}

func openLogFile(path string) (*os.File, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create log directory %q: %w", dir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFileMode)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", path, err)
	}

	return file, nil
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
