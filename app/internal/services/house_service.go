package services

import (
	"context"
	"strings"

	"github.com/gosimple/slug"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type HouseService struct {
	repository *repositories.HouseRepository
	awsCfg     *configs.AwsConfig
}

func NewHouseService(repo *repositories.HouseRepository, awsCfg *configs.AwsConfig) *HouseService {
	return &HouseService{
		repository: repo,
		awsCfg:     awsCfg,
	}
}

func (s *HouseService) GetAll(ctx context.Context) ([]schemas.HouseListItem, error) {
	return s.repository.GetAll(ctx)
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

func (s *HouseService) Update(ctx context.Context, id int, req schemas.HouseUpdateRequest) (models.House, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *HouseService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *HouseService) CheckSlug(ctx context.Context, rawSlug string) (bool, string, error) {
	normalized := slug.Make(rawSlug)
	exists, err := s.repository.SlugExists(ctx, normalized)
	if err != nil {
		return false, "", err
	}
	return !exists, normalized, nil
}
