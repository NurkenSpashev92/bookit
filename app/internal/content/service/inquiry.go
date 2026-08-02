package service

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
)

type InquiryRepository interface {
	GetAll(ctx context.Context, search string) ([]schema.Inquiry, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.Inquiry, int, error)
	GetByID(ctx context.Context, id int) (schema.Inquiry, error)
	Create(ctx context.Context, req schema.InquiryCreateRequest) (schema.Inquiry, error)
	Update(ctx context.Context, id int, req schema.InquiryUpdateRequest) (schema.Inquiry, error)
	Delete(ctx context.Context, id int) error
}

type InquiryService struct {
	repository InquiryRepository
}

func NewInquiryService(repo InquiryRepository) *InquiryService {
	return &InquiryService{repository: repo}
}

func (s *InquiryService) GetAll(ctx context.Context, search string) ([]schema.Inquiry, error) {
	return s.repository.GetAll(ctx, search)
}

func (s *InquiryService) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.Inquiry, int, error) {
	return s.repository.GetAllPaginated(ctx, search, limit, offset)
}

func (s *InquiryService) GetByID(ctx context.Context, id int) (schema.Inquiry, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *InquiryService) Create(ctx context.Context, req schema.InquiryCreateRequest) (schema.Inquiry, error) {
	return s.repository.Create(ctx, req)
}

func (s *InquiryService) Update(ctx context.Context, id int, req schema.InquiryUpdateRequest) (schema.Inquiry, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *InquiryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
