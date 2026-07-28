package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/imageproc"
)

type AvatarService struct {
	repository UserRepository
	s3         *aws.AwsS3Client
	mapper     UserMapper
	awsCfg     *configs.AwsConfig
}

func NewAvatarService(repo UserRepository, s3 *aws.AwsS3Client, awsCfg *configs.AwsConfig) *AvatarService {
	return &AvatarService{
		repository: repo,
		s3:         s3,
		awsCfg:     awsCfg,
	}
}

// Upload replaces the user avatar and returns the updated profile,
// so the client can refresh its state without an extra /auth/me call.
func (s *AvatarService) Upload(ctx context.Context, userID int, file *multipart.FileHeader) (schema.AuthUser, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return schema.AuthUser{}, ErrUserNotFound
	}

	result, err := imageproc.Process(file)
	if err != nil {
		return schema.AuthUser{}, fmt.Errorf("failed to process image: %w", err)
	}

	key := fmt.Sprintf("avatars/%d_%d.jpg", userID, time.Now().UnixNano())
	if _, err := s.s3.UploadCompressed(ctx, key, result.Original, "image/jpeg"); err != nil {
		return schema.AuthUser{}, fmt.Errorf("failed to upload avatar: %w", err)
	}

	previous := user.Avatar
	if err := s.repository.UpdateAvatar(ctx, userID, key); err != nil {
		_ = s.s3.Delete(ctx, key)
		return schema.AuthUser{}, fmt.Errorf("failed to save avatar: %w", err)
	}

	if previous != "" {
		_ = s.s3.Delete(ctx, previous)
	}

	user.Avatar = key
	return s.mapper.ToAuthUser(user, s.awsCfg), nil
}

// Delete removes the user avatar and returns the updated profile.
func (s *AvatarService) Delete(ctx context.Context, userID int) (schema.AuthUser, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return schema.AuthUser{}, ErrUserNotFound
	}

	if user.Avatar == "" {
		return schema.AuthUser{}, ErrAvatarNotFound
	}

	if err := s.repository.UpdateAvatar(ctx, userID, ""); err != nil {
		return schema.AuthUser{}, fmt.Errorf("failed to remove avatar: %w", err)
	}

	_ = s.s3.Delete(ctx, user.Avatar)

	user.Avatar = ""
	return s.mapper.ToAuthUser(user, s.awsCfg), nil
}

func (s *AvatarService) GetByUserID(ctx context.Context, userID int) (model.User, error) {
	return s.repository.GetByID(ctx, userID)
}
