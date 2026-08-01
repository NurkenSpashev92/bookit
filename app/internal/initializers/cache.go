package initializers

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/pkg/cache"
)

func NewCache(cacheCfg *configs.CacheConfig, redisCfg *configs.RedisConfig) (*cache.Cache, error) {
	if !cacheCfg.Enabled {
		return cache.NewDisabled(), nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, err
	}

	local := cache.NewLocalCache(cacheCfg.LocalTTL, cacheCfg.LocalEntries)

	return cache.NewWithLocal(client, cacheCfg.TTL, local), nil
}
