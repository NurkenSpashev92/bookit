package configs

import (
	"strconv"
	"time"
)

const (
	defaultCacheTTL      = 5 * time.Minute
	defaultLocalCacheTTL = 2 * time.Second
	defaultLocalEntries  = 1024
)

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
	Enabled      bool
	TTL          time.Duration
	LocalTTL     time.Duration
	LocalEntries int
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:      getEnv("CACHE_ENABLED", "true") != "false",
		TTL:          getCacheTTL(getEnv("CACHE_TTL", "300")),
		LocalTTL:     getLocalTTL(getEnv("CACHE_LOCAL_TTL", "")),
		LocalEntries: getLocalEntries(getEnv("CACHE_LOCAL_ENTRIES", "")),
	}
}

func getCacheTTL(val string) time.Duration {
	sec, err := strconv.Atoi(val)
	if err != nil || sec <= 0 {
		return defaultCacheTTL
	}
	return time.Duration(sec) * time.Second
}

func getLocalTTL(val string) time.Duration {
	sec, err := strconv.Atoi(val)
	if err != nil || sec < 0 {
		return defaultLocalCacheTTL
	}
	return time.Duration(sec) * time.Second
}

func getLocalEntries(val string) int {
	count, err := strconv.Atoi(val)
	if err != nil || count < 0 {
		return defaultLocalEntries
	}
	return count
}
