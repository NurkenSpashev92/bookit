package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type ConvenienceRepository interface {
	GetConveniences(ctx context.Context, search string) ([]schema.ConveniencePaginate, error)
	GetConveniencesPaginated(ctx context.Context, search string, limit, offset int) ([]schema.ConveniencePaginate, int, error)
	GetByID(ctx context.Context, id int) (schema.Convenience, error)
	CreateConvenience(ctx context.Context, req schema.ConvenienceCreateRequest) (schema.Convenience, error)
	Update(ctx context.Context, id int, req schema.ConvenienceUpdateRequest) (schema.Convenience, error)
	Delete(ctx context.Context, id int) error
}

type ConvenienceService struct {
	repository ConvenienceRepository
}

func NewConvenienceService(repo ConvenienceRepository) *ConvenienceService {
	return &ConvenienceService{repository: repo}
}

func (s *ConvenienceService) GetAll(ctx context.Context, search string) ([]schema.ConveniencePaginate, error) {
	return s.repository.GetConveniences(ctx, search)
}

func (s *ConvenienceService) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.ConveniencePaginate, int, error) {
	return s.repository.GetConveniencesPaginated(ctx, search, limit, offset)
}

func (s *ConvenienceService) GetByID(ctx context.Context, id int) (schema.Convenience, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *ConvenienceService) Create(ctx context.Context, req schema.ConvenienceCreateRequest) (schema.Convenience, error) {
	req.Slug = utils.GenerateSlug(req.Slug, req.Name, "", "")
	return s.repository.CreateConvenience(ctx, req)
}

func (s *ConvenienceService) Update(ctx context.Context, id int, req schema.ConvenienceUpdateRequest) (schema.Convenience, error) {
	if req.Slug != nil {
		slug := utils.GenerateSlug(*req.Slug, "", "", "")
		req.Slug = &slug
	}
	return s.repository.Update(ctx, id, req)
}

func (s *ConvenienceService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
