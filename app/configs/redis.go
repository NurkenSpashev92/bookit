package configs

import (
	"strconv"
	"time"
)

const defaultCacheTTL = 5 * time.Minute

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	}
}

type CacheConfig struct {
	Enabled bool
	TTL     time.Duration
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled: getEnv("CACHE_ENABLED", "true") != "false",
		TTL:     getCacheTTL(getEnv("CACHE_TTL", "300")),
	}
}

func getCacheTTL(val string) time.Duration {
	sec, err := strconv.Atoi(val)
	if err != nil || sec <= 0 {
		return defaultCacheTTL
	}
	return time.Duration(sec) * time.Second
}
