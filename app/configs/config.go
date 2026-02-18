package configs

import (
	"os"
	"strconv"
	"time"
)

type AuthConfig struct {
	JWTSecret string
	JWTExpire time.Duration
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		JWTExpire: getJWTExpire(getEnv("JWT_EXPIRE", "86400")),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getJWTExpire(val string) time.Duration {
	sec, err := strconv.Atoi(val)
	if err != nil {
		sec = 86400
	}
	return time.Duration(sec) * time.Second
}
