package handler

import (
	"context"

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
	GetUserLikedHouses(ctx context.Context, userID int) ([]propertyschema.HouseListItem, error)
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
// @Tags Houses
// @Produce json
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

	houses, err := h.likeService.GetUserLikedHouses(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.List(c, houses)
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
