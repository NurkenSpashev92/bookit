package test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"

	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	interactionschema "github.com/nurkenspashev92/bookit/internal/interaction/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

func TestHouseLikeHandler_Like_RequiresAuth(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Post("/houses/:id/like", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		return c.JSON(interactionschema.HouseLikeResponse{Liked: true, LikeCount: 1})
	})

	// No cookie
	resp := doRequest(t, app, http.MethodPost, "/houses/1/like", nil)
	if resp.StatusCode != 401 {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestHouseLikeHandler_Like_WithAuth(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Post("/houses/:id/like", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		user := c.Locals("user").(identitymodel.User)
		return c.JSON(interactionschema.HouseLikeResponse{Liked: true, LikeCount: user.ID})
	})

	token := generateTestToken(jwtSvc, 5, "user@test.com")
	cookie := &http.Cookie{Name: "jwt", Value: token}

	resp := doRequest(t, app, http.MethodPost, "/houses/1/like", nil, cookie)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body interactionschema.HouseLikeResponse
	parseJSON(t, resp, &body)
	if !body.Liked {
		t.Error("expected liked=true")
	}
	if body.LikeCount != 5 {
		t.Errorf("like_count = %d, want 5 (user ID)", body.LikeCount)
	}
}

func TestHouseLikeHandler_UserLikedHouses_EmptyList(t *testing.T) {
	jwtSvc := newTestJWTService()
	app := newTestApp()

	app.Get("/houses/liked", middleware.AuthRequired(jwtSvc), func(c fiber.Ctx) error {
		return c.JSON([]propertyschema.HouseListItem{})
	})

	token := generateTestToken(jwtSvc, 1, "u@t.com")
	cookie := &http.Cookie{Name: "jwt", Value: token}

	resp := doRequest(t, app, http.MethodGet, "/houses/liked", nil, cookie)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var items []propertyschema.HouseListItem
	parseJSON(t, resp, &items)
	if len(items) != 0 {
		t.Errorf("expected empty list, got %d", len(items))
	}
}
