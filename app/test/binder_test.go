package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	locationschema "github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

func TestBind(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
		wantError  string
	}{
		{
			name:       "valid payload",
			body:       locationschema.CountryCreateRequest{NameKZ: "Qazaqstan", NameEN: "Kazakhstan", NameRU: "Казахстан", Code: "KZ"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "body of the wrong shape",
			body:       []string{"Kazakhstan"},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid body",
		},
		{
			name:       "validation failure",
			body:       locationschema.CountryCreateRequest{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			app.Post("/countries", func(c fiber.Ctx) error {
				request, err := shared.Bind[locationschema.CountryCreateRequest](c)
				if err != nil {
					return shared.Fail(c, err)
				}
				return c.JSON(request)
			})

			resp := doRequest(t, app, http.MethodPost, "/countries", tt.body)
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantError != "" {
				var errResp shared.ErrorResponse
				parseJSON(t, resp, &errResp)
				if errResp.Error != tt.wantError {
					t.Errorf("error = %q, want %q", errResp.Error, tt.wantError)
				}
			}
		})
	}
}

func TestParamInt(t *testing.T) {
	app := newTestApp()
	app.Get("/countries/:id", func(c fiber.Ctx) error {
		id, err := shared.ParamInt(c, "id")
		if err != nil {
			return shared.Fail(c, err)
		}
		return c.JSON(fiber.Map{"id": id})
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"numeric id", "/countries/7", http.StatusOK},
		{"non numeric id", "/countries/abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, app, http.MethodGet, tt.path, nil)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
