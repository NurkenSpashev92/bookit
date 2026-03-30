package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type UserService struct {
	repository *repositories.UserRepository
	jwtService *JWTService
	mapper     UserMapper
}

func NewUserService(repo *repositories.UserRepository, jwtService *JWTService) *UserService {
	return &UserService{
		repository: repo,
		jwtService: jwtService,
	}
}

func (s *UserService) Register(ctx context.Context, req schemas.UserCreateRequest) (*schemas.AuthResponse, error) {
	user, err := s.repository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{User: authUser, Token: token}, nil
}

func (s *UserService) Login(ctx context.Context, req schemas.UserLoginRequest) (*schemas.AuthResponse, error) {
	user, err := s.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{User: authUser, Token: token}, nil
}

func (s *UserService) ValidateTokenAndGetUser(token string) (*schemas.AuthResponse, error) {
	user, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{User: authUser, Token: token}, nil
}

func (s *UserService) Me(ctx context.Context, token string) (*schemas.AuthResponse, error) {
	tokenUser, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	user, err := s.repository.GetByEmail(ctx, tokenUser.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	authUser := s.mapper.ToAuthUser(user)
	return &schemas.AuthResponse{User: authUser, Token: token}, nil
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
