package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/identity/schema"
)

type UserRepository interface {
	Create(ctx context.Context, req schema.UserCreateRequest) (model.User, error)
	GetByID(ctx context.Context, id int) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (model.User, error)
	Update(ctx context.Context, userID int, req schema.UserUpdateRequest) (model.User, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	UpdateAvatar(ctx context.Context, userID int, avatar string) error
	UpdatePaymentQR(ctx context.Context, userID int, key string) error
	ListAll(ctx context.Context, search string) ([]model.User, error)
	ListPaginated(ctx context.Context, search string, limit, offset int) ([]model.User, int, error)
	UpdateUser(ctx context.Context, id int, req schema.UserAdminUpdateRequest) (model.User, error)
}

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

func (s *UserService) Register(ctx context.Context, req schema.UserCreateRequest) (*schema.AuthResponse, error) {
	req.Email = normalizeEmail(req.Email)

	user, err := s.repository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schema.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req schema.UserLoginRequest) (*schema.AuthResponse, error) {
	var user model.User
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
	return &schema.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (*schema.AuthResponse, error) {
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
	return &schema.AuthResponse{
		User:         authUser,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID int, req schema.ChangePasswordRequest) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrWrongPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repository.UpdatePassword(ctx, userID, string(hashed))
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int, req schema.UserUpdateRequest) (*schema.AuthUser, error) {
	user, err := s.repository.Update(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &authUser, nil
}

func (s *UserService) Me(ctx context.Context, accessToken string) (*schema.AuthResponse, error) {
	userID, err := s.jwtService.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	authUser := s.mapper.ToAuthUser(user, s.awsCfg)
	return &schema.AuthResponse{
		User:        authUser,
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, search string) ([]schema.AdminUser, error) {
	users, err := s.repository.ListAll(ctx, search)
	if err != nil {
		return nil, err
	}

	result := make([]schema.AdminUser, 0, len(users))
	for _, user := range users {
		result = append(result, s.mapper.ToAdminUser(user, s.awsCfg))
	}

	return result, nil
}

func (s *UserService) ListUsersPaginated(ctx context.Context, search string, limit, offset int) ([]schema.AdminUser, int, error) {
	users, total, err := s.repository.ListPaginated(ctx, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]schema.AdminUser, 0, len(users))
	for _, user := range users {
		result = append(result, s.mapper.ToAdminUser(user, s.awsCfg))
	}

	return result, total, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req schema.UserAdminUpdateRequest) (*schema.AdminUser, error) {
	if req.Email != nil {
		normalized := normalizeEmail(*req.Email)
		req.Email = &normalized
	}

	user, err := s.repository.UpdateUser(ctx, id, req)
	if err != nil {
		return nil, err
	}

	adminUser := s.mapper.ToAdminUser(user, s.awsCfg)
	return &adminUser, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

type UserMapper struct{}

func (m *UserMapper) ToAuthUser(user model.User, awsCfg *configs.AwsConfig) schema.AuthUser {
	var phoneNumber string
	if user.PhoneNumber != nil {
		phoneNumber = *user.PhoneNumber
	}

	var dateOfBirth string
	if user.DateOfBirth != nil {
		dateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	var paymentPhone string
	if user.PaymentPhone != nil {
		paymentPhone = *user.PaymentPhone
	}

	return schema.AuthUser{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		MiddleName:       user.MiddleName,
		PhoneNumber:      phoneNumber,
		DateOfBirth:      dateOfBirth,
		Avatar:           awsCfg.AwsS3URL(user.Avatar),
		PaymentQR:        awsCfg.AwsS3URL(user.PaymentQR),
		PaymentPhone:     paymentPhone,
		IsSuperuser:      user.IsSuperuser,
		IsActive:         user.IsActive,
		SubscriptionType: user.SubscriptionType,
	}
}

func (m *UserMapper) ToAdminUser(user model.User, awsCfg *configs.AwsConfig) schema.AdminUser {
	var phoneNumber string
	if user.PhoneNumber != nil {
		phoneNumber = *user.PhoneNumber
	}

	return schema.AdminUser{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		MiddleName:       user.MiddleName,
		Avatar:           awsCfg.AwsS3URL(user.Avatar),
		PhoneNumber:      phoneNumber,
		IsSuperuser:      user.IsSuperuser,
		IsActive:         user.IsActive,
		SubscriptionType: user.SubscriptionType,
	}
}
