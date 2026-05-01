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

// Get fetches and unmarshals a value into dest. Caller-supplied context is honored
// so a slow Redis cannot block the request beyond its own timeout.
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) bool {
	getCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	data, err := c.client.Get(getCtx, key).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false
	}
	return true
}

// Set marshals and stores a value under TTL. Errors are silenced — caching is best-effort.
func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	setCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	c.client.Set(setCtx, key, data, c.ttl)
}

func (c *Cache) Delete(ctx context.Context, key string) {
	c.client.Del(ctx, key)
}

func (c *Cache) DeleteByPrefix(prefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
