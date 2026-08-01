package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/identity/model"
)

type JWTService struct {
	secretKey     string
	accessExpire  time.Duration
	refreshExpire time.Duration
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

func (j *JWTService) GenerateTokenPair(user model.User) (*TokenPair, error) {
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

func (j *JWTService) generateAccessToken(user model.User) (string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":          user.ID,
		"type":         "access",
		"jti":          jti,
		"is_superuser": user.IsSuperuser,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(j.accessExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) generateRefreshToken(user model.User) (string, error) {
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
	claims, err := j.parseClaims(tokenStr, expectedType)
	if err != nil {
		return 0, err
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid sub claim")
	}

	return int(sub), nil
}

func (j *JWTService) parseClaims(tokenStr, expectedType string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != expectedType {
		return nil, errors.New("wrong token type")
	}

	return claims, nil
}

func (j *JWTService) GenerateToken(user model.User) (string, error) {
	return j.generateAccessToken(user)
}

func (j *JWTService) ValidateToken(tokenStr string) (model.User, error) {
	claims, err := j.parseClaims(tokenStr, "access")
	if err != nil {
		return model.User{}, err
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return model.User{}, errors.New("invalid sub claim")
	}

	isSuperuser, _ := claims["is_superuser"].(bool)

	return model.User{ID: int(sub), IsSuperuser: isSuperuser}, nil
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
