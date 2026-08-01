package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

type RequestLoggerConfig struct {
	QuietPaths map[string]struct{}
}

func RequestLogger(cfg RequestLoggerConfig) fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		LoggerFunc: func(c fiber.Ctx, data *fiberlogger.Data, _ *fiberlogger.Config) error {
			status := c.Response().StatusCode()

			level := recordLevel(status, c.Path(), cfg.QuietPaths)
			if !logger.Enabled(level) {
				return nil
			}

			fields := []any{
				"method", c.Method(),
				"path", c.Path(),
				"status", status,
				"latency_ms", data.Stop.Sub(data.Start).Milliseconds(),
				"ip", c.IP(),
				"bytes", len(c.Response().Body()),
			}

			if query := string(c.Request().URI().QueryString()); query != "" {
				fields = append(fields, "query", query)
			}

			if data.ChainErr != nil {
				fields = append(fields, "error", data.ChainErr)
			}

			writeRecord(fiberlog.WithContext(c.Context()), level, fields)

			return nil
		},
	})
}

func recordLevel(status int, path string, quiet map[string]struct{}) fiberlog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return fiberlog.LevelError
	case status >= http.StatusBadRequest:
		return fiberlog.LevelWarn
	}

	if _, muted := quiet[path]; muted {
		return fiberlog.LevelDebug
	}

	return fiberlog.LevelInfo
}

func writeRecord(log fiberlog.CommonLogger, level fiberlog.Level, fields []any) {
	switch level {
	case fiberlog.LevelError:
		log.Errorw("request", fields...)
	case fiberlog.LevelWarn:
		log.Warnw("request", fields...)
	case fiberlog.LevelDebug:
		log.Debugw("request", fields...)
	default:
		log.Infow("request", fields...)
	}
}
