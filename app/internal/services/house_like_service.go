package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type HouseLikeService struct {
	repository HouseLikeRepository
}

func NewHouseLikeService(repo HouseLikeRepository) *HouseLikeService {
	return &HouseLikeService{repository: repo}
}

func (s *HouseLikeService) Like(ctx context.Context, userID, houseID int) (*schemas.HouseLikeResponse, error) {
	count, err := s.repository.LikeReturningCount(ctx, userID, houseID)
	if err != nil {
		return nil, err
	}
	return &schemas.HouseLikeResponse{Liked: true, LikeCount: count}, nil
}

func (s *HouseLikeService) Unlike(ctx context.Context, userID, houseID int) (*schemas.HouseLikeResponse, error) {
	count, err := s.repository.UnlikeReturningCount(ctx, userID, houseID)
	if err != nil {
		return nil, err
	}
	return &schemas.HouseLikeResponse{Liked: false, LikeCount: count}, nil
}

func (s *HouseLikeService) Status(ctx context.Context, userID, houseID int) (*schemas.HouseLikeResponse, error) {
	liked, count, err := s.repository.StatusWithCount(ctx, userID, houseID)
	if err != nil {
		return nil, err
	}
	return &schemas.HouseLikeResponse{Liked: liked, LikeCount: count}, nil
}

func (s *HouseLikeService) GetUserLikedHouses(ctx context.Context, userID int) ([]schemas.HouseLikeItem, error) {
	return s.repository.GetUserLikedHouses(ctx, userID)
}
