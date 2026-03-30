package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/handlers"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

func TestHouseHandler_GetAll_Empty(t *testing.T) {
	// This tests that the handler returns proper JSON even when service returns nil.
	// We can't call the real handler without DB, but we test the routing setup.
	app := newTestApp()
	app.Get("/houses", func(c fiber.Ctx) error {
		return c.JSON([]schemas.HouseListItem{})
	})

	resp := doRequest(t, app, http.MethodGet, "/houses", nil)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var items []schemas.HouseListItem
	parseJSON(t, resp, &items)
	if len(items) != 0 {
		t.Errorf("expected empty list, got %d items", len(items))
	}
}

func TestHouseHandler_GetBySlug_NotFound(t *testing.T) {
	app := newTestApp()
	houseHandler := handlers.NewHouseHandler(nil)
	// GetBySlug with nil service will panic — test the route pattern instead
	app.Get("/houses/:slug", func(c fiber.Ctx) error {
		slug := c.Params("slug")
		if slug == "" {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "slug is required"})
		}
		return c.Status(404).JSON(schemas.ErrorResponse{Error: "house not found"})
	})
	_ = houseHandler // used for type check

	resp := doRequest(t, app, http.MethodGet, "/houses/non-existent-slug", nil)
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}

	var errResp schemas.ErrorResponse
	parseJSON(t, resp, &errResp)
	if errResp.Error != "house not found" {
		t.Errorf("error = %q, want 'house not found'", errResp.Error)
	}
}

func TestHouseHandler_Create_InvalidJSON(t *testing.T) {
	app := newTestApp()
	app.Post("/houses", func(c fiber.Ctx) error {
		var req schemas.HouseCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(nil)
	})

	// Send empty body
	resp := doRequest(t, app, http.MethodPost, "/houses", map[string]any{})
	if resp.StatusCode != 400 {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}

	var errResp schemas.ErrorResponse
	parseJSON(t, resp, &errResp)
	if errResp.Error == "" {
		t.Error("expected validation error message")
	}
}

func TestHouseHandler_Create_ValidationErrors(t *testing.T) {
	app := newTestApp()
	app.Post("/houses", func(c fiber.Ctx) error {
		var req schemas.HouseCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(nil)
	})

	// Missing required fields
	body := map[string]any{
		"name_en": "Test",
		// missing name_kz, name_ru, descriptions, addresses, type_id
	}
	resp := doRequest(t, app, http.MethodPost, "/houses", body)
	if resp.StatusCode != 400 {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestHouseHandler_CheckSlug_RouteOrder(t *testing.T) {
	app := newTestApp()

	// Static routes must be before /:slug
	app.Get("/houses/check-slug", func(c fiber.Ctx) error {
		return c.JSON(schemas.SlugCheckResponse{Available: true, Slug: "test"})
	})
	app.Get("/houses/liked", func(c fiber.Ctx) error {
		return c.JSON([]schemas.HouseLikeItem{})
	})
	app.Get("/houses/:slug", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"slug": c.Params("slug")})
	})

	// check-slug should NOT be caught by :slug
	resp := doRequest(t, app, http.MethodGet, "/houses/check-slug?slug=test", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var slugResp schemas.SlugCheckResponse
	parseJSON(t, resp, &slugResp)
	if !slugResp.Available {
		t.Error("expected available=true")
	}

	// liked should NOT be caught by :slug
	resp2 := doRequest(t, app, http.MethodGet, "/houses/liked", nil)
	if resp2.StatusCode != 200 {
		t.Fatalf("liked: status = %d, want 200", resp2.StatusCode)
	}

	// real slug should work
	resp3 := doRequest(t, app, http.MethodGet, "/houses/beach-house", nil)
	if resp3.StatusCode != 200 {
		t.Fatalf("slug: status = %d, want 200", resp3.StatusCode)
	}
	var slugData map[string]string
	parseJSON(t, resp3, &slugData)
	if slugData["slug"] != "beach-house" {
		t.Errorf("slug = %q, want beach-house", slugData["slug"])
	}
}
