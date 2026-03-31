package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Dashboard godoc
// @Summary      Owner dashboard stats
// @Description  Returns overall statistics for the authenticated user's houses
// @Tags         Stats
// @Produce      json
// @Success      200  {object} schemas.DashboardStats
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/dashboard [get]
func (h *StatsHandler) Dashboard(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	stats, err := h.statsService.GetDashboard(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(stats)
}

// HouseStats godoc
// @Summary      Per-house statistics
// @Description  Returns views, likes, price for each house owned by user
// @Tags         Stats
// @Produce      json
// @Success      200  {array}  schemas.HouseStatsItem
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/houses [get]
func (h *StatsHandler) HouseStats(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	items, err := h.statsService.GetHouseStats(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	if items == nil {
		items = []schemas.HouseStatsItem{}
	}

	return c.JSON(items)
}

// Charts godoc
// @Summary      Chart data for owner dashboard
// @Description  Returns top by views, top by likes, likes over time, price distribution
// @Tags         Stats
// @Produce      json
// @Param        days query int false "Days for likes chart" default(30)
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/charts [get]
func (h *StatsHandler) Charts(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	days := 30
	if d, err := strconv.Atoi(c.Query("days", "30")); err == nil && d > 0 && d <= 365 {
		days = d
	}

	charts, err := h.statsService.GetCharts(c.Context(), user.ID, days)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(charts)
}

// HouseDetail godoc
// @Summary      Detailed stats for a single house
// @Description  Returns viewers, likers, bookings, views/likes per day charts for a house owned by current user
// @Tags         Stats
// @Produce      json
// @Param        slug path string true "House slug"
// @Success      200  {object} schemas.HouseDetailStats
// @Failure      401  {object} schemas.ErrorResponse
// @Failure      404  {object} schemas.ErrorResponse
// @Failure      500  {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/houses/{slug} [get]
func (h *StatsHandler) HouseDetail(c fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "slug is required"})
	}

	stats, err := h.statsService.GetHouseDetailStats(c.Context(), user.ID, slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(stats)
}
