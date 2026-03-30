package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type CityService struct {
	repository *repositories.CityRepository
}

func NewCityService(repo *repositories.CityRepository) *CityService {
	return &CityService{repository: repo}
}

func (s *CityService) GetAll(ctx context.Context) ([]schemas.City, error) {
	return s.repository.GetAllWithCountry(ctx)
}

func (s *CityService) GetByID(ctx context.Context, id int) (schemas.City, error) {
	return s.repository.GetByIDWithCountry(ctx, id)
}

func (s *CityService) Create(ctx context.Context, req schemas.CityCreateRequest) (models.City, error) {
	return s.repository.Create(ctx, req)
}

func (s *CityService) Update(ctx context.Context, id int, req schemas.CityUpdateRequest) (models.City, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *CityService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
