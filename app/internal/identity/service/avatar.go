package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/imageproc"
)

type AvatarService struct {
	repository UserRepository
	s3         *aws.AwsS3Client
}

func NewAvatarService(repo UserRepository, s3 *aws.AwsS3Client) *AvatarService {
	return &AvatarService{
		repository: repo,
		s3:         s3,
	}
}

func (s *AvatarService) Upload(ctx context.Context, userID int, file *multipart.FileHeader) (string, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if user.Avatar != "" {
		_ = s.s3.Delete(ctx, user.Avatar)
	}

	result, err := imageproc.Process(file)
	if err != nil {
		return "", fmt.Errorf("failed to process image: %w", err)
	}

	key := fmt.Sprintf("avatars/%d_%d.jpg", userID, time.Now().UnixNano())
	if _, err := s.s3.UploadCompressed(ctx, key, result.Original, "image/jpeg"); err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	if err := s.repository.UpdateAvatar(ctx, userID, key); err != nil {
		_ = s.s3.Delete(ctx, key)
		return "", fmt.Errorf("failed to save avatar: %w", err)
	}

	return key, nil
}

func (s *AvatarService) Delete(ctx context.Context, userID int) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.Avatar == "" {
		return ErrAvatarNotFound
	}

	if err := s.repository.UpdateAvatar(ctx, userID, ""); err != nil {
		return fmt.Errorf("failed to remove avatar: %w", err)
	}

	_ = s.s3.Delete(ctx, user.Avatar)
	return nil
}

func (s *AvatarService) GetByUserID(ctx context.Context, userID int) (model.User, error) {
	return s.repository.GetByID(ctx, userID)
}
