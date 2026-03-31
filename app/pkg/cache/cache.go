package cache

import (
	"sync"
	"time"
)

type entry struct {
	data      interface{}
	expiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]entry
	ttl     time.Duration
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]entry),
		ttl:     ttl,
	}
	go c.cleanup()
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.data, true
}

func (c *Cache) Set(key string, data interface{}) {
	c.mu.Lock()
	c.entries[key] = entry{data: data, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

func (c *Cache) DeleteByPrefix(prefix string) {
	c.mu.Lock()
	for k := range c.entries {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(c.entries, k)
		}
	}
	c.mu.Unlock()
}

func (c *Cache) Flush() {
	c.mu.Lock()
	c.entries = make(map[string]entry)
	c.mu.Unlock()
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.entries {
			if now.After(e.expiresAt) {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
