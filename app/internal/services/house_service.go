package services

import (
	"context"
	"strings"

	"github.com/gosimple/slug"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type HouseService struct {
	repository     HouseRepository
	likeRepository HouseLikeRepository
}

func NewHouseService(repo HouseRepository, likeRepo HouseLikeRepository) *HouseService {
	return &HouseService{repository: repo, likeRepository: likeRepo}
}

func (s *HouseService) GetAllPaginated(ctx context.Context, userID, limit, offset int) ([]schemas.HouseListItem, int, error) {
	houses, total, err := s.repository.GetAllPaginated(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if userID > 0 {
		likedIDs, err := s.likeRepository.GetUserLikedHouseIDs(ctx, userID)
		if err != nil {
			return houses, total, nil
		}
		liked := make(map[int]bool, len(likedIDs))
		for _, id := range likedIDs {
			liked[id] = true
		}
		for i := range houses {
			houses[i].IsLiked = liked[houses[i].ID]
		}
	}

	return houses, total, nil
}

func (s *HouseService) GetMyHouses(ctx context.Context, ownerID, limit, offset int) ([]schemas.HouseListItem, int, error) {
	return s.repository.GetByOwnerPaginated(ctx, ownerID, limit, offset)
}

func (s *HouseService) GetBySlug(ctx context.Context, slug string) (models.House, error) {
	return s.repository.GetBySlug(ctx, slug)
}

func (s *HouseService) Create(ctx context.Context, req schemas.HouseCreateRequest, ownerID int) (models.House, error) {
	req.OwnerID = ownerID
	house, err := s.repository.Create(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "slug already exists") {
			return house, ErrSlugExists
		}
		return house, err
	}
	return house, nil
}

func (s *HouseService) Update(ctx context.Context, slug string, req schemas.HouseUpdateRequest) (models.House, error) {
	return s.repository.Update(ctx, slug, req)
}

func (s *HouseService) Delete(ctx context.Context, slug string) error {
	return s.repository.Delete(ctx, slug)
}

func (s *HouseService) CheckSlug(ctx context.Context, rawSlug string) (bool, string, error) {
	normalized := slug.Make(rawSlug)
	exists, err := s.repository.SlugExists(ctx, normalized)
	if err != nil {
		return false, "", err
	}
	return !exists, normalized, nil
}
