package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]schema.CategoryPaginate, error)
	GetByID(ctx context.Context, id int) (model.Category, error)
	CreateCategory(ctx context.Context, req schema.CategoryCreateRequest) (model.Category, error)
	Update(ctx context.Context, id int, req schema.CategoryUpdateRequest) (model.Category, error)
	Delete(ctx context.Context, id int) error
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repository: repo}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]schema.CategoryPaginate, error) {
	return s.repository.GetCategories(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (model.Category, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, req schema.CategoryCreateRequest) (model.Category, error) {
	return s.repository.CreateCategory(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, id int, req schema.CategoryUpdateRequest) (model.Category, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
