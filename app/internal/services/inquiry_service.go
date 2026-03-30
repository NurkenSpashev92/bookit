package services

import (
	"context"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type InquiryService struct {
	repository InquiryRepository
}

func NewInquiryService(repo InquiryRepository) *InquiryService {
	return &InquiryService{repository: repo}
}

func (s *InquiryService) GetAll(ctx context.Context) ([]schemas.Inquiry, error) {
	return s.repository.GetAll(ctx)
}

func (s *InquiryService) GetByID(ctx context.Context, id int) (schemas.Inquiry, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *InquiryService) Create(ctx context.Context, req schemas.InquiryCreateRequest) (schemas.Inquiry, error) {
	return s.repository.Create(ctx, req)
}

func (s *InquiryService) Update(ctx context.Context, id int, req schemas.InquiryUpdateRequest) (schemas.Inquiry, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *InquiryService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
