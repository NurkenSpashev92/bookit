package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		client: client,
		ttl:    ttl,
	}
}

func (c *Cache) Get(key string, dest interface{}) bool {
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false
	}
	return true
}

func (c *Cache) Set(key string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	c.client.Set(context.Background(), key, data, c.ttl)
}

func (c *Cache) Delete(key string) {
	c.client.Del(context.Background(), key)
}

func (c *Cache) DeleteByPrefix(prefix string) {
	ctx := context.Background()
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return
		}
		if len(keys) > 0 {
			c.client.Del(ctx, keys...)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
}

func (c *Cache) Flush() {
	c.client.FlushDB(context.Background())
}
