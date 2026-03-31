package services

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type StatsService struct {
	repository *repositories.StatsRepository
}

func NewStatsService(repo *repositories.StatsRepository) *StatsService {
	return &StatsService{repository: repo}
}

func (s *StatsService) GetDashboard(ctx context.Context, ownerID int) (schemas.DashboardStats, error) {
	return s.repository.GetDashboard(ctx, ownerID)
}

func (s *StatsService) GetHouseStats(ctx context.Context, ownerID int) ([]schemas.HouseStatsItem, error) {
	return s.repository.GetHouseStats(ctx, ownerID)
}

func (s *StatsService) GetHouseDetailStats(ctx context.Context, ownerID int, slug string) (schemas.HouseDetailStats, error) {
	return s.repository.GetHouseDetailStats(ctx, ownerID, slug)
}

func (s *StatsService) GetCharts(ctx context.Context, ownerID, days int) (map[string]interface{}, error) {
	g, gctx := errgroup.WithContext(ctx)

	var topViews []schemas.TopHouseStats
	var topLikes []schemas.TopHouseStats
	var likesOverTime []schemas.LikesOverTime
	var priceDistribution []schemas.PriceDistribution

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
		topViews = []schemas.TopHouseStats{}
	}
	if topLikes == nil {
		topLikes = []schemas.TopHouseStats{}
	}
	if likesOverTime == nil {
		likesOverTime = []schemas.LikesOverTime{}
	}
	if priceDistribution == nil {
		priceDistribution = []schemas.PriceDistribution{}
	}

	return map[string]interface{}{
		"top_by_views":       topViews,
		"top_by_likes":       topLikes,
		"likes_over_time":    likesOverTime,
		"price_distribution": priceDistribution,
	}, nil
}
