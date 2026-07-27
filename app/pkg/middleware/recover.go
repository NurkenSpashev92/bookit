package middleware

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
)

func RecoverPanic() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			log.Printf("[%s %s] panic: %v\n%s", c.Method(), c.Path(), p, debug.Stack())

			if panicErr, ok := p.(error); ok {
				err = fmt.Errorf("panic: %w", panicErr)
				return
			}

			err = fmt.Errorf("panic: %v", p)
		}()

		return c.Next()
	}
}
