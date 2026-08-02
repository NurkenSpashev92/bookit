package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type TypeRepository interface {
	GetAll(ctx context.Context, search string) ([]model.Type, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]model.Type, int, error)
	GetByID(ctx context.Context, id int) (model.Type, error)
	Create(ctx context.Context, t model.Type) (model.Type, error)
	Update(ctx context.Context, id int, t model.Type) (model.Type, error)
	Delete(ctx context.Context, id int) error
}

type TypeService struct {
	repository TypeRepository
	mapper     TypeMapper
}

func NewTypeService(repo TypeRepository) *TypeService {
	return &TypeService{repository: repo}
}

func (s *TypeService) GetAll(ctx context.Context, search string) ([]schema.TypeResponse, error) {
	types, err := s.repository.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	responses := make([]schema.TypeResponse, 0, len(types))
	for i := range types {
		responses = append(responses, s.mapper.ToResponse(types[i]))
	}

	return responses, nil
}

func (s *TypeService) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.TypeResponse, int, error) {
	types, total, err := s.repository.GetAllPaginated(ctx, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]schema.TypeResponse, 0, len(types))
	for i := range types {
		responses = append(responses, s.mapper.ToResponse(types[i]))
	}

	return responses, total, nil
}

func (s *TypeService) GetByID(ctx context.Context, id int) (schema.TypeResponse, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return schema.TypeResponse{}, err
	}

	return s.mapper.ToResponse(t), nil
}

func (s *TypeService) Create(ctx context.Context, req schema.TypeCreateRequest) (schema.TypeResponse, error) {
	t := model.Type{IsActive: true}
	s.mapper.ApplyCreateRequest(&t, req)

	created, err := s.repository.Create(ctx, t)
	if err != nil {
		return schema.TypeResponse{}, err
	}

	return s.mapper.ToResponse(created), nil
}

func (s *TypeService) Update(ctx context.Context, id int, req schema.TypeUpdateRequest) (schema.TypeResponse, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return schema.TypeResponse{}, err
	}

	s.mapper.ApplyUpdateRequest(&t, req)

	updated, err := s.repository.Update(ctx, id, t)
	if err != nil {
		return schema.TypeResponse{}, err
	}

	return s.mapper.ToResponse(updated), nil
}

func (s *TypeService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

type TypeMapper struct{}

func (m *TypeMapper) ApplyCreateRequest(t *model.Type, req schema.TypeCreateRequest) {
	t.NameKz = req.NameKz
	t.NameRu = req.NameRu
	t.NameEn = req.NameEn
	t.Slug = utils.GenerateSlug(req.Slug, req.NameEn, req.NameKz, req.NameRu)

	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
}

func (m *TypeMapper) ApplyUpdateRequest(t *model.Type, req schema.TypeUpdateRequest) {
	if req.NameKz != nil {
		t.NameKz = *req.NameKz
	}
	if req.NameRu != nil {
		t.NameRu = *req.NameRu
	}
	if req.NameEn != nil {
		t.NameEn = *req.NameEn
	}
	if req.Slug != nil {
		t.Slug = utils.GenerateSlug(*req.Slug, "", "", "")
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
}

func (m *TypeMapper) ToResponse(t model.Type) schema.TypeResponse {
	return schema.TypeResponse{
		ID:       t.ID,
		NameKz:   t.NameKz,
		NameRu:   t.NameRu,
		NameEn:   t.NameEn,
		Slug:     t.Slug,
		IsActive: t.IsActive,
	}
}
