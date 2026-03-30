package services

import (
	"context"
	"mime/multipart"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

type CategoryService struct {
	repository *repositories.CategoryRepository
	s3         *aws.AwsS3Client
	awsCfg     *configs.AwsConfig
}

func NewCategoryService(repo *repositories.CategoryRepository, s3 *aws.AwsS3Client, awsCfg *configs.AwsConfig) *CategoryService {
	return &CategoryService{
		repository: repo,
		s3:         s3,
		awsCfg:     awsCfg,
	}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]schemas.CategoryPaginate, error) {
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

func (s *CategoryService) GetByID(ctx context.Context, id int) (models.Category, error) {
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

func (s *CategoryService) Create(ctx context.Context, nameKz, nameRu, nameEn string, isActive bool, iconFile *multipart.FileHeader) (models.Category, error) {
	var iconURL *string

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "categories", iconFile)
		if err != nil {
			return models.Category{}, err
		}
		full := s.awsCfg.AwsS3URL(uploaded)
		iconURL = &full
	}

	req := schemas.CategoryCreateRequest{
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

func (s *CategoryService) Update(ctx context.Context, id int, req schemas.CategoryUpdateRequest, iconFile *multipart.FileHeader) (models.Category, error) {
	var newIcon *string

	if iconFile != nil {
		uploaded, err := s.s3.Upload(ctx, "categories", iconFile)
		if err != nil {
			return models.Category{}, err
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
