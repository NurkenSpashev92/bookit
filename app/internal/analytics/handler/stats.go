package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/analytics/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const (
	defaultChartDays = 30
	minChartDays     = 1
	maxChartDays     = 365
)

type StatsService interface {
	GetDashboard(ctx context.Context, ownerID int) (schema.DashboardStats, error)
	GetHouseStats(ctx context.Context, ownerID int) ([]schema.HouseStatsItem, error)
	GetHouseDetailStats(ctx context.Context, ownerID int, slug string) (schema.HouseDetailStats, error)
	GetCharts(ctx context.Context, ownerID, days int) (map[string]interface{}, error)
}

type StatsHandler struct {
	statsService StatsService
}

func NewStatsHandler(statsService StatsService) *StatsHandler {
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	stats, err := h.statsService.GetDashboard(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	items, err := h.statsService.GetHouseStats(c.Context(), user.ID)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.List(c, items)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	days := shared.QueryIntInRange(c, "days", defaultChartDays, minChartDays, maxChartDays)

	charts, err := h.statsService.GetCharts(c.Context(), user.ID, days)
	if err != nil {
		return shared.Fail(c, err)
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
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	slug, err := shared.ParamString(c, "slug")
	if err != nil {
		return shared.Fail(c, err)
	}

	stats, err := h.statsService.GetHouseDetailStats(c.Context(), user.ID, slug)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(stats)
}
