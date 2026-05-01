package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	contentschema "github.com/nurkenspashev92/bookit/internal/content/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

func TestFAQHandler_Create_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/faqs", func(c fiber.Ctx) error {
		var req contentschema.FAQCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(req)
	})

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{
			"valid",
			map[string]any{
				"question_kz": "Q?", "answer_kz": "A",
				"question_ru": "Q?", "answer_ru": "A",
				"question_en": "Q?", "answer_en": "A",
			},
			201,
		},
		{"empty", map[string]any{}, 400},
		{
			"missing answer",
			map[string]any{
				"question_kz": "Q?",
				"question_ru": "Q?",
				"question_en": "Q?",
			},
			400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/faqs", tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestInquiryHandler_Create_Validation(t *testing.T) {
	app := newTestApp()
	app.Post("/inquiry", func(c fiber.Ctx) error {
		var req contentschema.InquiryCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(req)
	})

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{"valid", map[string]any{"email": "a@b.com", "text": "Hello"}, 201},
		{"missing email", map[string]any{"text": "Hello"}, 400},
		{"invalid email", map[string]any{"email": "bad", "text": "Hello"}, 400},
		{"missing text", map[string]any{"email": "a@b.com"}, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodPost, "/inquiry", tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
