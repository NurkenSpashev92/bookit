package services

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type UserService struct {
	repository UserRepository
	jwtService *JWTService
	mapper     UserMapper
}

func NewUserService(repo UserRepository, jwtService *JWTService) *UserService {
	return &UserService{
		repository: repo,
		jwtService: jwtService,
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

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req schemas.UserLoginRequest) (*schemas.AuthResponse, error) {
	req.Email = normalizeEmail(req.Email)

	user, err := s.repository.GetByEmail(ctx, req.Email)
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

	authUser := s.mapper.ToAuthUser(user)
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

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
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

	authUser := s.mapper.ToAuthUser(user)
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

func (m *UserMapper) ToAuthUser(user models.User) schemas.AuthUser {
	return schemas.AuthUser{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Avatar:     user.Avatar,
	}
}
