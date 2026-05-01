package service

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/analytics/schema"
)

// StatsRepository describes the persistence contract StatsService depends on.
type StatsRepository interface {
	GetDashboard(ctx context.Context, ownerID int) (schema.DashboardStats, error)
	GetHouseStats(ctx context.Context, ownerID int) ([]schema.HouseStatsItem, error)
	GetHouseDetailStats(ctx context.Context, ownerID int, slug string) (schema.HouseDetailStats, error)
	GetTopByViews(ctx context.Context, ownerID, limit int) ([]schema.TopHouseStats, error)
	GetTopByLikes(ctx context.Context, ownerID, limit int) ([]schema.TopHouseStats, error)
	GetLikesOverTime(ctx context.Context, ownerID, days int) ([]schema.LikesOverTime, error)
	GetPriceDistribution(ctx context.Context, ownerID int) ([]schema.PriceDistribution, error)
}

type StatsService struct {
	repository StatsRepository
}

func NewStatsService(repo StatsRepository) *StatsService {
	return &StatsService{repository: repo}
}

func (s *StatsService) GetDashboard(ctx context.Context, ownerID int) (schema.DashboardStats, error) {
	return s.repository.GetDashboard(ctx, ownerID)
}

func (s *StatsService) GetHouseStats(ctx context.Context, ownerID int) ([]schema.HouseStatsItem, error) {
	return s.repository.GetHouseStats(ctx, ownerID)
}

func (s *StatsService) GetHouseDetailStats(ctx context.Context, ownerID int, slug string) (schema.HouseDetailStats, error) {
	return s.repository.GetHouseDetailStats(ctx, ownerID, slug)
}

func (s *StatsService) GetCharts(ctx context.Context, ownerID, days int) (map[string]interface{}, error) {
	g, gctx := errgroup.WithContext(ctx)

	var topViews []schema.TopHouseStats
	var topLikes []schema.TopHouseStats
	var likesOverTime []schema.LikesOverTime
	var priceDistribution []schema.PriceDistribution

	g.Go(func() error {
		var err error
		topViews, err = s.repository.GetTopByViews(gctx, ownerID, 10)
		return err
	})

	g.Go(func() error {
		var err error
		topLikes, err = s.repository.GetTopByLikes(gctx, ownerID, 10)
		return err
	})

	g.Go(func() error {
		var err error
		likesOverTime, err = s.repository.GetLikesOverTime(gctx, ownerID, days)
		return err
	})

	g.Go(func() error {
		var err error
		priceDistribution, err = s.repository.GetPriceDistribution(gctx, ownerID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if topViews == nil {
		topViews = []schema.TopHouseStats{}
	}
	if topLikes == nil {
		topLikes = []schema.TopHouseStats{}
	}
	if likesOverTime == nil {
		likesOverTime = []schema.LikesOverTime{}
	}
	if priceDistribution == nil {
		priceDistribution = []schema.PriceDistribution{}
	}

	return map[string]interface{}{
		"top_by_views":       topViews,
		"top_by_likes":       topLikes,
		"likes_over_time":    likesOverTime,
		"price_distribution": priceDistribution,
	}, nil
}
