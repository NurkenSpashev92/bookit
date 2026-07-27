package test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/nurkenspashev92/bookit/pkg/logger"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

func newLoggedApp(t *testing.T, quiet map[string]struct{}) (*fiber.App, *bytes.Buffer) {
	t.Helper()

	logger.Init(logger.Config{Level: "info"})

	var out bytes.Buffer
	fiberlog.SetOutput(&out)
	t.Cleanup(func() { fiberlog.SetOutput(os.Stdout) })

	app := fiber.New(fiber.Config{PassLocalsToContext: true})
	app.Use(requestid.New())
	app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{QuietPaths: quiet}))

	return app, &out
}

func TestRequestLoggerLevels(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		handler   fiber.Handler
		quiet     map[string]struct{}
		wantLevel string
		wantLine  bool
		wantError string
	}{
		{
			name:      "success is logged at info",
			path:      "/ok",
			handler:   func(c fiber.Ctx) error { return c.SendString("ok") },
			wantLevel: "[Info]",
			wantLine:  true,
		},
		{
			name:      "client error is logged at warn",
			path:      "/bad",
			handler:   func(c fiber.Ctx) error { return fiber.NewError(http.StatusBadRequest, "bad input") },
			wantLevel: "[Warn]",
			wantLine:  true,
		},
		{
			name:      "server error is logged at error",
			path:      "/boom",
			handler:   func(c fiber.Ctx) error { return errors.New("boom") },
			wantLevel: "[Error]",
			wantLine:  true,
			wantError: "error=boom",
		},
		{
			name:     "quiet path is dropped below info",
			path:     "/healthcheck",
			handler:  func(c fiber.Ctx) error { return c.SendString("ok") },
			quiet:    map[string]struct{}{"/healthcheck": {}},
			wantLine: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, out := newLoggedApp(t, tt.quiet)
			app.Get(tt.path, tt.handler)

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, tt.path, nil))
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			requestID := resp.Header.Get(fiber.HeaderXRequestID)
			if requestID == "" {
				t.Fatal("response must carry X-Request-ID")
			}

			line := out.String()

			if !tt.wantLine {
				if line != "" {
					t.Fatalf("quiet path must not be logged, got %q", line)
				}
				return
			}

			for _, want := range []string{
				tt.wantLevel,
				"[" + requestID + "]",
				"request",
				"method=GET",
				"path=" + tt.path,
				"latency_ms=",
			} {
				if !strings.Contains(line, want) {
					t.Errorf("log line %q must contain %q", line, want)
				}
			}

			if tt.wantError != "" && !strings.Contains(line, tt.wantError) {
				t.Errorf("log line %q must contain %q", line, tt.wantError)
			}
		})
	}
}

func TestRequestLoggerPropagatesRequestIDToContext(t *testing.T) {
	app, out := newLoggedApp(t, nil)
	app.Get("/ok", func(c fiber.Ctx) error {
		fiberlog.WithContext(c.Context()).Infow("handler ran")
		return c.SendString("ok")
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/ok", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	requestID := resp.Header.Get(fiber.HeaderXRequestID)
	line := out.String()

	if !strings.Contains(line, "handler ran") {
		t.Fatalf("handler record missing: %q", line)
	}
	if strings.Count(line, "["+requestID+"]") != 2 {
		t.Errorf("handler and access records must share request id %q, got %q", requestID, line)
	}
}
