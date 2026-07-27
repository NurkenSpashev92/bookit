package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const ResponsePrefix = "resp:"

const opTimeout = 200 * time.Millisecond

type Cache struct {
	client  *redis.Client
	ttl     time.Duration
	enabled bool
}

func New(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		client:  client,
		ttl:     ttl,
		enabled: true,
	}
}

func NewDisabled() *Cache {
	return &Cache{enabled: false}
}

func (c *Cache) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Cache) Close() error {
	if !c.Enabled() {
		return nil
	}
	return c.client.Close()
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) bool {
	data, ok := c.GetBytes(ctx, key)
	if !ok {
		return false
	}
	return json.Unmarshal(data, dest) == nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	if !c.enabled {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	c.SetBytes(ctx, key, data)
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	getCtx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()

	data, err := c.client.Get(getCtx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return data, true
}

func (c *Cache) SetBytes(ctx context.Context, key string, data []byte) {
	if !c.enabled {
		return
	}

	setCtx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()

	c.client.Set(setCtx, key, data, c.ttl)
}

func (c *Cache) Delete(ctx context.Context, key string) {
	if !c.enabled {
		return
	}
	c.client.Del(ctx, key)
}

func (c *Cache) InvalidateNamespace(ns string) {
	c.DeleteByPrefix(ns + ":")
	c.DeleteByPrefix(ResponsePrefix + ns + ":")
}

func (c *Cache) DeleteByPrefix(prefix string) {
	if !c.enabled {
		return
	}

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
	if !c.enabled {
		return
	}
	c.client.FlushDB(context.Background())
}
