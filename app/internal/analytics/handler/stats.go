package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/analytics/schema"
	"github.com/nurkenspashev92/bookit/internal/analytics/service"
	identitymodel "github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Dashboard godoc
// @Summary      Owner dashboard stats
// @Tags         Stats
// @Produce      json
// @Success      200  {object} schema.DashboardStats
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/dashboard [get]
func (h *StatsHandler) Dashboard(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	stats, err := h.statsService.GetDashboard(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(stats)
}

// HouseStats godoc
// @Summary      Per-house statistics
// @Tags         Stats
// @Produce      json
// @Success      200  {array}  schema.HouseStatsItem
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/houses [get]
func (h *StatsHandler) HouseStats(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	items, err := h.statsService.GetHouseStats(c.Context(), user.ID)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	if items == nil {
		items = []schema.HouseStatsItem{}
	}

	return c.JSON(items)
}

// Charts godoc
// @Summary      Chart data for owner dashboard
// @Tags         Stats
// @Produce      json
// @Param        days query int false "Days for likes chart" default(30)
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/charts [get]
func (h *StatsHandler) Charts(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	days := 30
	if d, err := strconv.Atoi(c.Query("days", "30")); err == nil && d > 0 && d <= 365 {
		days = d
	}

	charts, err := h.statsService.GetCharts(c.Context(), user.ID, days)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(charts)
}

// HouseDetail godoc
// @Summary      Detailed stats for a single house
// @Tags         Stats
// @Produce      json
// @Param        slug path string true "House slug"
// @Success      200  {object} schema.HouseDetailStats
// @Failure      401  {object} shared.ErrorResponse
// @Failure      404  {object} shared.ErrorResponse
// @Failure      500  {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /stats/houses/{slug} [get]
func (h *StatsHandler) HouseDetail(c fiber.Ctx) error {
	user := c.Locals("user").(identitymodel.User)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(shared.ErrorResponse{Error: "slug is required"})
	}

	stats, err := h.statsService.GetHouseDetailStats(c.Context(), user.ID, slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(404).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(stats)
}
