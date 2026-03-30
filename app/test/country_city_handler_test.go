package test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

func TestCountryHandler_Create_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/countries", func(c fiber.Ctx) error {
		var req schemas.CountryCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(req)
	})

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{"valid", map[string]any{"name_kz": "Қ", "name_en": "K", "name_ru": "К"}, 201},
		{"missing name_kz", map[string]any{"name_en": "K", "name_ru": "К"}, 400},
		{"missing name_en", map[string]any{"name_kz": "Қ", "name_ru": "К"}, 400},
		{"name too long", map[string]any{"name_kz": strings.Repeat("a", 256), "name_en": "K", "name_ru": "К"}, 400},
		{"code too long", map[string]any{"name_kz": "Қ", "name_en": "K", "name_ru": "К", "code": strings.Repeat("X", 11)}, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/countries", tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestCityHandler_Create_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/cities", func(c fiber.Ctx) error {
		var req schemas.CityCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(req)
	})

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{"valid", map[string]any{"name_kz": "А", "name_en": "A", "name_ru": "А", "country_id": 1}, 201},
		{"missing name", map[string]any{"name_en": "A", "name_ru": "А", "country_id": 1}, 400},
		{"missing country_id", map[string]any{"name_kz": "А", "name_en": "A", "name_ru": "А"}, 400},
		{"long postall_code", map[string]any{"name_kz": "А", "name_en": "A", "name_ru": "А", "country_id": 1, "postall_code": strings.Repeat("0", 21)}, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/cities", tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
