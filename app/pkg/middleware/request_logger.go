package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
)

type RequestLoggerConfig struct {
	QuietPaths map[string]struct{}
}

func RequestLogger(cfg RequestLoggerConfig) fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		LoggerFunc: func(c fiber.Ctx, data *fiberlogger.Data, _ *fiberlogger.Config) error {
			status := c.Response().StatusCode()

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

			log := fiberlog.WithContext(c.Context())

			switch {
			case status >= http.StatusInternalServerError:
				log.Errorw("request", fields...)
			case status >= http.StatusBadRequest:
				log.Warnw("request", fields...)
			default:
				if _, quiet := cfg.QuietPaths[c.Path()]; quiet {
					log.Debugw("request", fields...)
				} else {
					log.Infow("request", fields...)
				}
			}

			return nil
		},
	})
}
