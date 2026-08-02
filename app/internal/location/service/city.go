package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
)

type CityRepository interface {
	GetAllWithCountry(ctx context.Context, search string) ([]schema.City, error)
	GetAllWithCountryPaginated(ctx context.Context, search string, limit, offset int) ([]schema.City, int, error)
	GetByIDWithCountry(ctx context.Context, id int) (schema.City, error)
	Create(ctx context.Context, req schema.CityCreateRequest) (model.City, error)
	Update(ctx context.Context, id int, req schema.CityUpdateRequest) (model.City, error)
	Delete(ctx context.Context, id int) error
}

type CityService struct {
	repository CityRepository
}

func NewCityService(repo CityRepository) *CityService {
	return &CityService{repository: repo}
}

func (s *CityService) GetAll(ctx context.Context, search string) ([]schema.City, error) {
	return s.repository.GetAllWithCountry(ctx, search)
}

func (s *CityService) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.City, int, error) {
	return s.repository.GetAllWithCountryPaginated(ctx, search, limit, offset)
}

func (s *CityService) GetByID(ctx context.Context, id int) (schema.City, error) {
	return s.repository.GetByIDWithCountry(ctx, id)
}

func (s *CityService) Create(ctx context.Context, req schema.CityCreateRequest) (model.City, error) {
	return s.repository.Create(ctx, req)
}

func (s *CityService) Update(ctx context.Context, id int, req schema.CityUpdateRequest) (model.City, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *CityService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
