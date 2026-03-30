package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type CountryService struct {
	repository CountryRepository
}

func NewCountryService(repo CountryRepository) *CountryService {
	return &CountryService{repository: repo}
}

func (s *CountryService) GetAll(ctx context.Context) ([]models.Country, error) {
	return s.repository.GetAll(ctx)
}

func (s *CountryService) GetByID(ctx context.Context, id int) (models.Country, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *CountryService) Create(ctx context.Context, req schemas.CountryCreateRequest) (models.Country, error) {
	return s.repository.Create(ctx, req)
}

func (s *CountryService) Update(ctx context.Context, id int, req schemas.CountryUpdateRequest) (models.Country, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *CountryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
