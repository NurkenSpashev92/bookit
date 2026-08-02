package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/interaction/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type HouseLikeService interface {
	Like(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error)
	Unlike(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error)
	Status(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error)
	GetUserLikedHouses(ctx context.Context, userID int, search string) ([]propertyschema.HouseListItem, error)
	GetUserLikedHousesPaginated(ctx context.Context, userID int, search string, limit, offset int) ([]propertyschema.HouseListItem, int, error)
}

type HouseLikeHandler struct {
	likeService HouseLikeService
}

func NewHouseLikeHandler(likeService HouseLikeService) *HouseLikeHandler {
	return &HouseLikeHandler{likeService: likeService}
}

// Like godoc
// @Summary Like a house
// @Tags Houses
// @Produce json
// @Param slug path string true "House slug"
// @Success 200 {object} schema.HouseLikeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{slug}/like [post]
func (h *HouseLikeHandler) Like(c fiber.Ctx) error {
	user, slug, err := likeRequest(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	response, err := h.likeService.Like(c.Context(), user, slug)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(response)
}

// Unlike godoc
// @Summary Unlike a house
// @Tags Houses
// @Produce json
// @Param slug path string true "House slug"
// @Success 200 {object} schema.HouseLikeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{slug}/like [delete]
func (h *HouseLikeHandler) Unlike(c fiber.Ctx) error {
	user, slug, err := likeRequest(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	response, err := h.likeService.Unlike(c.Context(), user, slug)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(response)
}

// LikeStatus godoc
// @Summary Check if user liked a house
// @Tags Houses
// @Produce json
// @Param slug path string true "House slug"
// @Success 200 {object} schema.HouseLikeResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{slug}/like [get]
func (h *HouseLikeHandler) Status(c fiber.Ctx) error {
	user, slug, err := likeRequest(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	response, err := h.likeService.Status(c.Context(), user, slug)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(response)
}

// UserLikedHouses godoc
// @Summary Get houses liked by current user
// @Description  Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags Houses
// @Produce json
// @Param        search     query string false "Search by house name/address"
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success 200 {array} propertyschema.HouseListItem
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/liked [get]
func (h *HouseLikeHandler) UserLikedHouses(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]propertyschema.HouseListItem, error) {
			return h.likeService.GetUserLikedHouses(ctx, user.ID, search)
		},
		func(ctx context.Context, limit, offset int) ([]propertyschema.HouseListItem, int, error) {
			return h.likeService.GetUserLikedHousesPaginated(ctx, user.ID, search, limit, offset)
		},
	)
}

func likeRequest(c fiber.Ctx) (int, string, error) {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return 0, "", err
	}

	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return 0, "", err
	}

	return user.ID, slug, nil
}
