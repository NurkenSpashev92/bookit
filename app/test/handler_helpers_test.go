package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/configs"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
)

func newTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		CaseSensitive: true,
	})
}

func newTestJWTService() *identitysvc.JWTService {
	return identitysvc.NewJWTService(&configs.AuthConfig{
		JWTSecret: "test-secret",
		JWTExpire: 24 * time.Hour,
	})
}

func generateTestToken(jwtSvc *identitysvc.JWTService, userID int, email string) string {
	token, _ := jwtSvc.GenerateToken(identitymodel.User{
		ID:    userID,
		Email: email,
	})
	return token
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body any, cookies ...*http.Cookie) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func parseJSON(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
}
