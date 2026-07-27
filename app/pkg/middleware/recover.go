package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func RecoverPanic() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, p any) {
			fiberlog.WithContext(c.Context()).Errorw("panic recovered",
				"method", c.Method(),
				"path", c.Path(),
				"panic", p,
				"stack", string(debug.Stack()),
			)
		},
	})
}
