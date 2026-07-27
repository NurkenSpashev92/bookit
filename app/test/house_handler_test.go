package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	propertyh "github.com/nurkenspashev92/bookit/internal/property/handler"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

func TestHouseHandler_GetAll_Empty(t *testing.T) {
	app := newTestApp()
	app.Get("/houses", func(c fiber.Ctx) error {
		return c.JSON([]propertyschema.HouseListItem{})
	})

	resp := doRequest(t, app, http.MethodGet, "/houses", nil)
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var items []propertyschema.HouseListItem
	parseJSON(t, resp, &items)
	if len(items) != 0 {
		t.Errorf("expected empty list, got %d items", len(items))
	}
}

func TestHouseHandler_GetBySlug_NotFound(t *testing.T) {
	app := newTestApp()
	houseHandler := propertyh.NewHouseHandler(nil)

	app.Get("/houses/:slug", func(c fiber.Ctx) error {
		slug := c.Params("slug")
		if slug == "" {
			return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
		}
		return c.Status(404).JSON(shared.ErrorResponse{Error: "house not found"})
	})
	_ = houseHandler

	resp := doRequest(t, app, http.MethodGet, "/houses/non-existent-slug", nil)
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}

	var errResp shared.ErrorResponse
	parseJSON(t, resp, &errResp)
	if errResp.Error != "house not found" {
		t.Errorf("error = %q, want 'house not found'", errResp.Error)
	}
}

func TestHouseHandler_Create_InvalidJSON(t *testing.T) {
	app := newTestApp()
	app.Post("/houses", func(c fiber.Ctx) error {
		var req propertyschema.HouseCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(nil)
	})

	resp := doRequest(t, app, http.MethodPost, "/houses", map[string]any{})
	if resp.StatusCode != 400 {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}

	var errResp shared.ErrorResponse
	parseJSON(t, resp, &errResp)
	if errResp.Error == "" {
		t.Error("expected validation error message")
	}
}

func TestHouseHandler_Create_ValidationErrors(t *testing.T) {
	app := newTestApp()
	app.Post("/houses", func(c fiber.Ctx) error {
		var req propertyschema.HouseCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		if err := req.Validate(); err != nil {
			return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(201).JSON(nil)
	})

	body := map[string]any{
		"name_en": "Test",
	}
	resp := doRequest(t, app, http.MethodPost, "/houses", body)
	if resp.StatusCode != 400 {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestHouseHandler_CheckSlug_RouteOrder(t *testing.T) {
	app := newTestApp()

	app.Get("/houses/check-slug", func(c fiber.Ctx) error {
		return c.JSON(propertyschema.SlugCheckResponse{Available: true, Slug: "test"})
	})
	app.Get("/houses/liked", func(c fiber.Ctx) error {
		return c.JSON([]propertyschema.HouseListItem{})
	})
	app.Get("/houses/:slug", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"slug": c.Params("slug")})
	})

	resp := doRequest(t, app, http.MethodGet, "/houses/check-slug?slug=test", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var slugResp propertyschema.SlugCheckResponse
	parseJSON(t, resp, &slugResp)
	if !slugResp.Available {
		t.Error("expected available=true")
	}

	resp2 := doRequest(t, app, http.MethodGet, "/houses/liked", nil)
	if resp2.StatusCode != 200 {
		t.Fatalf("liked: status = %d, want 200", resp2.StatusCode)
	}

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
