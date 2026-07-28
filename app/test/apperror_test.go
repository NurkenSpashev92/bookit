package test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	propertymodel "github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

func TestStatusOf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"invalid", shared.Invalid("bad input"), http.StatusBadRequest},
		{"unauthorized", shared.Unauthorized("unauthenticated"), http.StatusUnauthorized},
		{"forbidden", shared.Forbidden("not allowed"), http.StatusForbidden},
		{"not found", shared.NotFound("house not found"), http.StatusNotFound},
		{"conflict", shared.Conflict("slug already exists"), http.StatusConflict},
		{"wrapped domain error", fmt.Errorf("repository: %w", propertymodel.ErrHouseNotFound), http.StatusNotFound},
		{"validation errors", shared.ValidationErrors{"title is required"}, http.StatusBadRequest},
		{"fiber error", fiber.NewError(http.StatusRequestEntityTooLarge, "too large"), http.StatusRequestEntityTooLarge},
		{"unknown error", errors.New("connection reset"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shared.StatusOf(tt.err); got != tt.want {
				t.Errorf("StatusOf() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMessageOf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"domain message", shared.NotFound("house not found"), "house not found"},
		{"wrapped keeps client message", fmt.Errorf("slug %q: %w", "villa", propertymodel.ErrSlugExists), "slug already exists"},
		{"validation errors joined", shared.ValidationErrors{"title is required", "price is required"}, "title is required; price is required"},
		{"unknown error stays neutral", errors.New("pq: connection refused"), "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shared.MessageOf(tt.err); got != tt.want {
				t.Errorf("MessageOf() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDomainErrorsRemainComparable(t *testing.T) {
	wrapped := fmt.Errorf("repository: %w", propertymodel.ErrSlugExists)

	if !errors.Is(wrapped, propertymodel.ErrSlugExists) {
		t.Error("expected wrapped error to match its sentinel")
	}
	if errors.Is(wrapped, propertymodel.ErrHouseNotFound) {
		t.Error("expected distinct sentinels not to match")
	}
}

func TestFailUsesErrorKind(t *testing.T) {
	app := newTestApp()
	app.Get("/houses/:slug", func(c fiber.Ctx) error {
		return shared.Fail(c, fmt.Errorf("repository: %w", propertymodel.ErrHouseNotFound))
	})

	resp := doRequest(t, app, http.MethodGet, "/houses/missing", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	var body shared.ErrorResponse
	parseJSON(t, resp, &body)
	if body.Error != "house not found" {
		t.Errorf("error = %q, want 'house not found'", body.Error)
	}
}

func TestFailHidesInternalDetails(t *testing.T) {
	app := newTestApp()
	app.Get("/houses", func(c fiber.Ctx) error {
		return shared.Fail(c, errors.New("pq: relation \"houses\" does not exist"))
	})

	resp := doRequest(t, app, http.MethodGet, "/houses", nil)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	var body shared.ErrorResponse
	parseJSON(t, resp, &body)
	if body.Error != "internal server error" {
		t.Errorf("error = %q, want 'internal server error'", body.Error)
	}
}
