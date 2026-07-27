package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/pkg/logger"
)

func RecoverPanic() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			logger.FromContext(c.Context()).ErrorContext(c.Context(), "panic recovered",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.Any("panic", p),
				slog.String("stack", string(debug.Stack())),
			)

			if panicErr, ok := p.(error); ok {
				err = fmt.Errorf("panic: %w", panicErr)
				return
			}

			err = fmt.Errorf("panic: %v", p)
		}()

		return c.Next()
	}
}
