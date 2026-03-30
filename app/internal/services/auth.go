package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
)

type JWTService struct {
	secretKey      string
	accessExpire   time.Duration
	refreshExpire  time.Duration
}

func NewJWTService(cnf *configs.AuthConfig) *JWTService {
	return &JWTService{
		secretKey:     cnf.JWTSecret,
		accessExpire:  15 * time.Minute,
		refreshExpire: 7 * 24 * time.Hour,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (j *JWTService) GenerateTokenPair(user models.User) (*TokenPair, error) {
	accessToken, err := j.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWTService) generateAccessToken(user models.User) (string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"type": "access",
		"jti":  jti,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(j.accessExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) generateRefreshToken(user models.User) (string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"type": "refresh",
		"jti":  jti,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(j.refreshExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateAccessToken(tokenStr string) (int, error) {
	return j.validateToken(tokenStr, "access")
}

func (j *JWTService) ValidateRefreshToken(tokenStr string) (int, error) {
	return j.validateToken(tokenStr, "refresh")
}

func (j *JWTService) validateToken(tokenStr, expectedType string) (int, error) {
	if tokenStr == "" {
		return 0, errors.New("empty token")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != expectedType {
		return 0, errors.New("wrong token type")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid sub claim")
	}

	return int(sub), nil
}

// GenerateToken backward compat for existing code that expects old signature.
func (j *JWTService) GenerateToken(user models.User) (string, error) {
	return j.generateAccessToken(user)
}

// ValidateToken backward compat — returns user with ID only.
func (j *JWTService) ValidateToken(tokenStr string) (models.User, error) {
	userID, err := j.ValidateAccessToken(tokenStr)
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: userID}, nil
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
