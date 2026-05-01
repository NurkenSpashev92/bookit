package service

import (
	"context"
	"mime/multipart"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

// CategoryRepository describes the persistence contract CategoryService depends on.
type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]schema.CategoryPaginate, error)
	GetByID(ctx context.Context, id int) (model.Category, error)
	CreateCategory(ctx context.Context, req schema.CategoryCreateRequest) (model.Category, error)
	Update(ctx context.Context, id int, req schema.CategoryUpdateRequest, icon *string) (model.Category, *string, error)
	Delete(ctx context.Context, id int) (*string, error)
}

type CategoryService struct {
	repository CategoryRepository
	s3         *aws.AwsS3Client
	awsCfg     *configs.AwsConfig
}

func NewCategoryService(repo CategoryRepository, s3 *aws.AwsS3Client, awsCfg *configs.AwsConfig) *CategoryService {
	return &CategoryService{
		repository: repo,
		s3:         s3,
		awsCfg:     awsCfg,
	}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]schema.CategoryPaginate, error) {
	categories, err := s.repository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	for i := range categories {
		if categories[i].Icon != nil && *categories[i].Icon != "" {
			full := s.awsCfg.AwsS3URL(*categories[i].Icon)
			categories[i].Icon = &full
		}
	}

	return categories, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (model.Category, error) {
	category, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return category, err
	}

	if category.Icon != nil && *category.Icon != "" {
		full := s.awsCfg.AwsS3URL(*category.Icon)
		category.Icon = &full
	}

	return category, nil
}

func (s *CategoryService) Create(ctx context.Context, nameKz, nameRu, nameEn string, isActive bool, iconFile *multipart.FileHeader) (model.Category, error) {
	var iconURL *string

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "categories", iconFile)
		if err != nil {
			return model.Category{}, err
		}
		full := s.awsCfg.AwsS3URL(uploaded)
		iconURL = &full
	}

	req := schema.CategoryCreateRequest{
		NameKz:   nameKz,
		NameRu:   nameRu,
		NameEn:   nameEn,
		Icon:     iconURL,
		IsActive: isActive,
	}

	category, err := s.repository.CreateCategory(ctx, req)
	if err != nil {
		if iconURL != nil {
			_ = s.s3.Delete(ctx, *iconURL)
		}
		return category, err
	}

	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id int, req schema.CategoryUpdateRequest, iconFile *multipart.FileHeader) (model.Category, error) {
	var newIcon *string

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "categories", iconFile)
		if err != nil {
			return model.Category{}, err
		}
		full := s.awsCfg.AwsS3URL(uploaded)
		newIcon = &full
	}

	category, oldIcon, err := s.repository.Update(ctx, id, req, newIcon)
	if err != nil {
		if newIcon != nil {
			_ = s.s3.Delete(ctx, *newIcon)
		}
		return category, err
	}

	if newIcon != nil && oldIcon != nil {
		_ = s.s3.Delete(ctx, *oldIcon)
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	icon, err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	if icon != nil {
		_ = s.s3.Delete(ctx, *icon)
	}

	return nil
}
