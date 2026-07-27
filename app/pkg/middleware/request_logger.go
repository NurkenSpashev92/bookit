package middleware

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

const HeaderRequestID = "X-Request-ID"

const maxInboundRequestIDLen = 64

type RequestLoggerConfig struct {
	Logger     *slog.Logger
	QuietPaths map[string]struct{}
}

func RequestLogger(cfg RequestLoggerConfig) fiber.Handler {
	root := cfg.Logger
	if root == nil {
		root = slog.Default()
	}

	return func(c fiber.Ctx) error {
		start := time.Now()

		requestID := resolveRequestID(c)
		c.Set(HeaderRequestID, requestID)

		log := root.With(slog.String(logger.KeyRequestID, requestID))

		ctx := logger.WithRequestID(c.Context(), requestID)
		c.SetContext(logger.WithContext(ctx, log))

		err := c.Next()

		status := responseStatus(c, err)

		attrs := []any{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("ip", c.IP()),
			slog.Int("bytes", len(c.Response().Body())),
		}

		if query := string(c.Request().URI().QueryString()); query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		if err != nil {
			attrs = append(attrs, logger.Err(err))
		}

		log.Log(c.Context(), levelFor(c.Path(), status, cfg.QuietPaths), "request", attrs...)

		return err
	}
}

func resolveRequestID(c fiber.Ctx) string {
	if id := c.Get(HeaderRequestID); id != "" && len(id) <= maxInboundRequestIDLen {
		return id
	}

	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:8], rand.Uint64())
	binary.LittleEndian.PutUint64(buf[8:], rand.Uint64())

	return hex.EncodeToString(buf[:])
}

func responseStatus(c fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return http.StatusInternalServerError
}

func levelFor(path string, status int, quiet map[string]struct{}) slog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return slog.LevelError
	case status >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		if _, ok := quiet[path]; ok {
			return slog.LevelDebug
		}

		return slog.LevelInfo
	}
}
