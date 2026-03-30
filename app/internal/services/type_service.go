package services

import (
	"context"
	"mime/multipart"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

type TypeService struct {
	repository *repositories.TypeRepository
	s3         *aws.AwsS3Client
	awsCfg     *configs.AwsConfig
}

func NewTypeService(repo *repositories.TypeRepository, s3 *aws.AwsS3Client, awsCfg *configs.AwsConfig) *TypeService {
	return &TypeService{
		repository: repo,
		s3:         s3,
		awsCfg:     awsCfg,
	}
}

func (s *TypeService) fillIconURL(t *models.Type) {
	if t.Icon != nil && *t.Icon != "" {
		url := s.awsCfg.AwsS3URL(*t.Icon)
		t.Icon = &url
	}
}

func (s *TypeService) GetAll(ctx context.Context) ([]models.Type, error) {
	types, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range types {
		s.fillIconURL(&types[i])
	}

	return types, nil
}

func (s *TypeService) GetByID(ctx context.Context, id int) (models.Type, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return t, err
	}

	s.fillIconURL(&t)
	return t, nil
}

func (s *TypeService) Create(ctx context.Context, name string, isActive bool, iconFile *multipart.FileHeader) (models.Type, error) {
	var iconPtr *string

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "types", iconFile)
		if err != nil {
			return models.Type{}, err
		}
		iconPtr = &uploaded
	}

	t := models.Type{
		Name:     name,
		Icon:     iconPtr,
		IsActive: isActive,
	}

	created, err := s.repository.Create(ctx, t)
	if err != nil {
		if iconPtr != nil {
			_ = s.s3.Delete(ctx, *iconPtr)
		}
		return created, err
	}

	s.fillIconURL(&created)
	return created, nil
}

func (s *TypeService) Update(ctx context.Context, id int, name *string, isActive *bool, iconFile *multipart.FileHeader) (models.Type, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return t, err
	}

	if name != nil {
		t.Name = *name
	}
	if isActive != nil {
		t.IsActive = *isActive
	}

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "types", iconFile)
		if err != nil {
			return t, err
		}

		if t.Icon != nil && *t.Icon != "" {
			_ = s.s3.Delete(ctx, *t.Icon)
		}
		t.Icon = &uploaded
	}

	updated, err := s.repository.Update(ctx, id, t)
	if err != nil {
		return updated, err
	}

	s.fillIconURL(&updated)
	return updated, nil
}

func (s *TypeService) Delete(ctx context.Context, id int) error {
	icon, err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	if icon != "" {
		_ = s.s3.Delete(ctx, icon)
	}

	return nil
}
