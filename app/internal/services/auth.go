package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
)

type JWTService struct {
	SecretKey string
	Expire    time.Duration
}

func NewJWTService(cnf *configs.AuthConfig) *JWTService {
	return &JWTService{
		SecretKey: cnf.JWTSecret,
		Expire:    cnf.JWTExpire,
	}
}

func (j *JWTService) GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"middle_name": user.MiddleName,
		"avatar":      user.Avatar,
		"exp":         time.Now().Add(j.Expire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (models.User, error) {
	var user models.User

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["user_id"].(float64); ok {
			user.ID = int(id)
		}
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		}
		if firstName, ok := claims["first_name"].(string); ok {
			user.FirstName = firstName
		}
		if lastName, ok := claims["last_name"].(string); ok {
			user.LastName = lastName
		}
		if middleName, ok := claims["middle_name"].(string); ok {
			user.MiddleName = middleName
		}
		if avatar, ok := claims["avatar"].(string); ok {
			user.Avatar = avatar
		}
		return user, nil
	}

	return user, errors.New("invalid token claims")
}
