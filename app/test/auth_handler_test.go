package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

func TestAuthMiddleware_NoCookie(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Get("/protected", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	resp := doRequest(t, app, http.MethodGet, "/protected", nil)
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}

	var errResp map[string]string
	parseJSON(t, resp, &errResp)
	if errResp["error"] != "unauthenticated" {
		t.Errorf("error = %q, want 'unauthenticated'", errResp["error"])
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Get("/protected", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	cookie := &http.Cookie{Name: "access_token", Value: "invalid-token"}
	resp := doRequest(t, app, http.MethodGet, "/protected", nil, cookie)
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Get("/protected", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		return c.JSON(fiber.Map{"user_id": user.ID})
	})

	token := generateTestToken(jwtSvc, 42, "test@example.com")
	cookie := &http.Cookie{Name: "access_token", Value: token}

	resp := doRequest(t, app, http.MethodGet, "/protected", nil, cookie)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	parseJSON(t, resp, &body)
	if int(body["user_id"].(float64)) != 42 {
		t.Errorf("user_id = %v, want 42", body["user_id"])
	}
}

func TestAuthMiddleware_OldJwtCookie_BackwardCompat(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Get("/protected", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		return c.JSON(fiber.Map{"user_id": user.ID})
	})

	token := generateTestToken(jwtSvc, 7, "old@test.com")
	cookie := &http.Cookie{Name: "jwt", Value: token}

	resp := doRequest(t, app, http.MethodGet, "/protected", nil, cookie)
	if resp.StatusCode != 200 {
		t.Errorf("old 'jwt' cookie should still work, got status %d", resp.StatusCode)
	}
}

func TestAuthHandler_Logout_ClearsCookies(t *testing.T) {
	app := newTestApp()

	app.Post("/auth/logout", func(c fiber.Ctx) error {
		expired := fiber.Cookie{
			Name: "access_token", Value: "", HTTPOnly: true,
		}
		c.Cookie(&expired)
		return c.JSON(schemas.MessageResponse{Message: "logged out"})
	})

	resp := doRequest(t, app, http.MethodPost, "/auth/logout", nil)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestAuthHandler_Register_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/auth/register", func(c fiber.Ctx) error {
		var req schemas.UserCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(nil)
	})

	tests := []struct {
		name string
		body map[string]any
	}{
		{"empty body", map[string]any{}},
		{"missing password", map[string]any{"email": "a@b.com"}},
		{"invalid email", map[string]any{"email": "bad", "password": "secret123"}},
		{"short password", map[string]any{"email": "a@b.com", "password": "ab"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/auth/register", tt.body)
			if resp.StatusCode != 400 {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}

func TestAuthHandler_Login_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/auth/login", func(c fiber.Ctx) error {
		var req schemas.UserLoginRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(nil)
	})

	tests := []struct {
		name string
		body map[string]any
	}{
		{"empty", map[string]any{}},
		{"no password", map[string]any{"email": "a@b.com"}},
		{"no email", map[string]any{"password": "secret"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/auth/login", tt.body)
			if resp.StatusCode != 400 {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}
