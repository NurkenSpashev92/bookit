package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
)

type CountryRepository interface {
	GetAll(ctx context.Context) ([]model.Country, error)
	GetByID(ctx context.Context, id int) (model.Country, error)
	Create(ctx context.Context, req schema.CountryCreateRequest) (model.Country, error)
	Update(ctx context.Context, id int, req schema.CountryUpdateRequest) (model.Country, error)
	Delete(ctx context.Context, id int) error
}

type CountryService struct {
	repository CountryRepository
}

func NewCountryService(repo CountryRepository) *CountryService {
	return &CountryService{repository: repo}
}

func (s *CountryService) GetAll(ctx context.Context) ([]model.Country, error) {
	return s.repository.GetAll(ctx)
}

func (s *CountryService) GetByID(ctx context.Context, id int) (model.Country, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *CountryService) Create(ctx context.Context, req schema.CountryCreateRequest) (model.Country, error) {
	return s.repository.Create(ctx, req)
}

func (s *CountryService) Update(ctx context.Context, id int, req schema.CountryUpdateRequest) (model.Country, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *CountryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
