package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
)

type FAQRepository interface {
	GetAll(ctx context.Context) ([]schema.FAQ, error)
	GetByID(ctx context.Context, id int) (schema.FAQ, error)
	Create(ctx context.Context, req schema.FAQCreateRequest) (schema.FAQ, error)
	Update(ctx context.Context, id int, req schema.FAQUpdateRequest) (schema.FAQ, error)
	Delete(ctx context.Context, id int) error
}

type FAQService struct {
	repository FAQRepository
}

func NewFAQService(repo FAQRepository) *FAQService {
	return &FAQService{repository: repo}
}

func (s *FAQService) GetAll(ctx context.Context) ([]schema.FAQ, error) {
	return s.repository.GetAll(ctx)
}

func (s *FAQService) GetByID(ctx context.Context, id int) (schema.FAQ, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *FAQService) Create(ctx context.Context, req schema.FAQCreateRequest) (schema.FAQ, error) {
	return s.repository.Create(ctx, req)
}

func (s *FAQService) Update(ctx context.Context, id int, req schema.FAQUpdateRequest) (schema.FAQ, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *FAQService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
