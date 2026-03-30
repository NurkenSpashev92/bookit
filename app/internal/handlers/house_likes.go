package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type HouseLikeHandler struct {
	likeService *services.HouseLikeService
}

func NewHouseLikeHandler(likeService *services.HouseLikeService) *HouseLikeHandler {
	return &HouseLikeHandler{likeService: likeService}
}

// Like godoc
// @Summary Like a house
// @Tags Houses
// @Produce json
// @Param id path int true "House ID"
// @Success 200 {object} schemas.HouseLikeResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{id}/like [post]
func (h *HouseLikeHandler) Like(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid house id"})
	}

	resp, err := h.likeService.Like(c.Context(), user.ID, houseID)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
}

// Unlike godoc
// @Summary Unlike a house
// @Tags Houses
// @Produce json
// @Param id path int true "House ID"
// @Success 200 {object} schemas.HouseLikeResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{id}/like [delete]
func (h *HouseLikeHandler) Unlike(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid house id"})
	}

	resp, err := h.likeService.Unlike(c.Context(), user.ID, houseID)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
}

// LikeStatus godoc
// @Summary Check if user liked a house
// @Tags Houses
// @Produce json
// @Param id path int true "House ID"
// @Success 200 {object} schemas.HouseLikeResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/{id}/like [get]
func (h *HouseLikeHandler) Status(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid house id"})
	}

	resp, err := h.likeService.Status(c.Context(), user.ID, houseID)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
}

// UserLikedHouses godoc
// @Summary Get houses liked by current user
// @Tags Houses
// @Produce json
// @Success 200 {array} schemas.HouseLikeItem
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /houses/liked [get]
func (h *HouseLikeHandler) UserLikedHouses(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	houses, err := h.likeService.GetUserLikedHouses(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	if houses == nil {
		houses = []schemas.HouseLikeItem{}
	}

	return c.JSON(houses)
}
