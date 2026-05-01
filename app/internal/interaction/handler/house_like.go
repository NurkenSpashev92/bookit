package handler

import (
	"github.com/gofiber/fiber/v3"

	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/interaction/service"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type HouseLikeHandler struct {
	likeService *service.HouseLikeService
}

func NewHouseLikeHandler(likeService *service.HouseLikeService) *HouseLikeHandler {
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
	user := c.Locals("user").(identitymodel.User)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	resp, err := h.likeService.Like(c.Context(), user.ID, slug)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
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
	user := c.Locals("user").(identitymodel.User)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	resp, err := h.likeService.Unlike(c.Context(), user.ID, slug)
	if err != nil {
		return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
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
	user := c.Locals("user").(identitymodel.User)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	resp, err := h.likeService.Status(c.Context(), user.ID, slug)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
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
	user := c.Locals("user").(identitymodel.User)

	houses, err := h.likeService.GetUserLikedHouses(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []propertyschema.HouseListItem{}
	}

	return c.JSON(houses)
}
