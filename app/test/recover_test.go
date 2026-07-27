package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"

	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/logger"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

func TestRecoverPanicRendersServerError(t *testing.T) {
	logger.Init(logger.Config{Level: "info"})

	var out bytes.Buffer
	fiberlog.SetOutput(&out)
	t.Cleanup(func() { fiberlog.SetOutput(os.Stdout) })

	app := fiber.New(fiber.Config{ErrorHandler: shared.ErrorHandler, PassLocalsToContext: true})
	app.Use(middleware.RecoverPanic())
	app.Get("/boom", func(_ fiber.Ctx) error { panic("db exploded") })

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/boom", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "internal server error") {
		t.Errorf("body %q must not leak panic details", body)
	}

	line := out.String()
	for _, want := range []string{"panic recovered", "panic=db exploded", "path=/boom"} {
		if !strings.Contains(line, want) {
			t.Errorf("log %q must contain %q", line, want)
		}
	}
}

func TestUploadLimitsRejectsOversizedBody(t *testing.T) {
	app := fiber.New()
	app.Post("/upload", middleware.UploadLimits(8), func(c fiber.Ctx) error {
		return c.SendString("stored")
	})

	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(strings.Repeat("x", 64)))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", resp.StatusCode)
	}
}
