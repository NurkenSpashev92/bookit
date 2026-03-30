package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type FAQService struct {
	repository FAQRepository
}

func NewFAQService(repo FAQRepository) *FAQService {
	return &FAQService{repository: repo}
}

func (s *FAQService) GetAll(ctx context.Context) ([]schemas.FAQ, error) {
	return s.repository.GetAll(ctx)
}

func (s *FAQService) GetByID(ctx context.Context, id int) (schemas.FAQ, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *FAQService) Create(ctx context.Context, req schemas.FAQCreateRequest) (schemas.FAQ, error) {
	return s.repository.Create(ctx, req)
}

func (s *FAQService) Update(ctx context.Context, id int, req schemas.FAQUpdateRequest) (schemas.FAQ, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *FAQService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
