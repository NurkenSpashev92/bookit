package apiserver

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"

	"github.com/nurkenspashev92/bookit/cmd/router"
	_ "github.com/nurkenspashev92/bookit/docs"
)

const shutdownTimeout = 10 * time.Second

type ApiApp struct {
	*fiber.App
}

func (app *ApiApp) Run() {
	deps, err := newContainer()
	if err != nil {
		fiberlog.Fatalw("failed to initialize dependencies", "error", err)
	}
	defer deps.Close()

	app.App = router.RegisterRoutes(app.App, deps.db, deps.services)

	startPprof()

	done := make(chan bool, 1)
	go app.Shutdown(done)
	go app.listen()

	<-done
	fiberlog.Info("graceful shutdown complete")
}

func (app *ApiApp) listen() {
	addr := "0.0.0.0:" + os.Getenv("APP_PORT")
	fiberlog.Infow("http server starting", "addr", addr)

	if err := app.Listen(addr); err != nil {
		fiberlog.Fatalw("http server error", "error", err)
	}
}

func (app *ApiApp) Shutdown(done chan<- bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	fiberlog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		fiberlog.Errorw("server forced to shutdown", "error", err)
	}

	fiberlog.Info("server exiting")

	done <- true
}

func startPprof() {
	port := os.Getenv("PPROF_PORT")
	if port == "" {
		return
	}

	go func() {
		addr := "0.0.0.0:" + port
		fiberlog.Infow("pprof listening", "addr", addr)

		if err := http.ListenAndServe(addr, nil); err != nil {
			fiberlog.Errorw("pprof server stopped", "error", err)
		}
	}()
}
