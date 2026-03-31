package services

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type UserService struct {
	repository UserRepository
	jwtService *JWTService
	mapper     UserMapper
	awsCfg     *configs.AwsConfig
}

func NewUserService(repo UserRepository, jwtService *JWTService, awsCfg *configs.AwsConfig) *UserService {
	return &UserService{
		repository: repo,
		jwtService: jwtService,
		awsCfg:     awsCfg,
	}
}

func (s *UserService) Register(ctx context.Context, req schemas.UserCreateRequest) (*schemas.AuthResponse, error) {
	req.Email = normalizeEmail(req.Email)

	user, err := s.repository.Create(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schemas.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req schemas.UserLoginRequest) (*schemas.AuthResponse, error) {
	var user models.User
	var err error

	if req.Email != "" {
		req.Email = normalizeEmail(req.Email)
		user, err = s.repository.GetByEmail(ctx, req.Email)
	} else {
		user, err = s.repository.GetByPhoneNumber(ctx, req.PhoneNumber)
	}
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schemas.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (*schemas.AuthResponse, error) {
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schemas.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID int, req schemas.ChangePasswordRequest) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repository.UpdatePassword(ctx, userID, string(hashed))
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int, req schemas.UserUpdateRequest) (*schemas.AuthUser, error) {
	user, err := s.repository.Update(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &authUser, nil
}

func (s *UserService) Me(ctx context.Context, accessToken string) (*schemas.AuthResponse, error) {
	userID, err := s.jwtService.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schemas.AuthResponse{
		User:        authUser,
		AccessToken: accessToken,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// UserMapper handles User model to DTO conversions.
type UserMapper struct{}

func (m *UserMapper) ToAuthUser(user models.User, awsCfg *configs.AwsConfig) schemas.AuthUser {
	var phoneNumber string
	if user.PhoneNumber != nil {
		phoneNumber = *user.PhoneNumber
	}

	var dateOfBirth string
	if user.DateOfBirth != nil {
		dateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	return schemas.AuthUser{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MiddleName:  user.MiddleName,
		PhoneNumber: phoneNumber,
		DateOfBirth: dateOfBirth,
		Avatar:      awsCfg.AwsS3URL(user.Avatar),
	}
}
