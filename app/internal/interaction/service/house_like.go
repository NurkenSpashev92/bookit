package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/interaction/schema"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
)

type HouseLikeRepository interface {
	LikeReturningCount(ctx context.Context, userID int, slug string) (int, error)
	UnlikeReturningCount(ctx context.Context, userID int, slug string) (int, error)
	StatusWithCount(ctx context.Context, userID int, slug string) (bool, int, error)
	GetUserLikedHouses(ctx context.Context, userID int) ([]propertyschema.HouseListItem, error)
	GetUserLikedHousesPaginated(ctx context.Context, userID, limit, offset int) ([]propertyschema.HouseListItem, int, error)
	GetUserLikedHouseIDs(ctx context.Context, userID int) ([]int, error)
}

type HouseLikeService struct {
	repository HouseLikeRepository
}

func NewHouseLikeService(repo HouseLikeRepository) *HouseLikeService {
	return &HouseLikeService{repository: repo}
}

func (s *HouseLikeService) Like(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error) {
	count, err := s.repository.LikeReturningCount(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	return &schema.HouseLikeResponse{Liked: true, LikeCount: count}, nil
}

func (s *HouseLikeService) Unlike(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error) {
	count, err := s.repository.UnlikeReturningCount(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	return &schema.HouseLikeResponse{Liked: false, LikeCount: count}, nil
}

func (s *HouseLikeService) Status(ctx context.Context, userID int, slug string) (*schema.HouseLikeResponse, error) {
	liked, count, err := s.repository.StatusWithCount(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	return &schema.HouseLikeResponse{Liked: liked, LikeCount: count}, nil
}

func (s *HouseLikeService) GetUserLikedHouses(ctx context.Context, userID int) ([]propertyschema.HouseListItem, error) {
	return s.repository.GetUserLikedHouses(ctx, userID)
}

func (s *HouseLikeService) GetFavoriteHouses(ctx context.Context, userID, limit, offset int) ([]propertyschema.HouseListItem, int, error) {
	return s.repository.GetUserLikedHousesPaginated(ctx, userID, limit, offset)
}
