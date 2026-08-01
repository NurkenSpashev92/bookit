package cache

import (
	"strings"
	"sync"
	"time"
)

type localEntry struct {
	data      []byte
	expiresAt time.Time
}

type LocalCache struct {
	mu      sync.RWMutex
	entries map[string]localEntry
	ttl     time.Duration
	limit   int
}

func NewLocalCache(ttl time.Duration, limit int) *LocalCache {
	if ttl <= 0 || limit <= 0 {
		return nil
	}

	return &LocalCache{
		entries: make(map[string]localEntry, limit),
		ttl:     ttl,
		limit:   limit,
	}
}

func (l *LocalCache) Get(key string) ([]byte, bool) {
	if l == nil {
		return nil, false
	}

	l.mu.RLock()
	entry, found := l.entries[key]
	l.mu.RUnlock()

	if !found || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.data, true
}

func (l *LocalCache) Set(key string, data []byte, ttl time.Duration) {
	if l == nil || len(data) == 0 {
		return
	}

	if ttl <= 0 || ttl > l.ttl {
		ttl = l.ttl
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.entries[key]; !exists && len(l.entries) >= l.limit {
		l.evict()
	}

	l.entries[key] = localEntry{data: data, expiresAt: time.Now().Add(ttl)}
}

func (l *LocalCache) Delete(key string) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.entries, key)
}

func (l *LocalCache) DeletePrefix(prefix string) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for key := range l.entries {
		if strings.HasPrefix(key, prefix) {
			delete(l.entries, key)
		}
	}
}

func (l *LocalCache) Clear() {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make(map[string]localEntry, l.limit)
}

func (l *LocalCache) evict() {
	now := time.Now()

	for key, entry := range l.entries {
		if now.After(entry.expiresAt) {
			delete(l.entries, key)
		}
	}

	for key := range l.entries {
		if len(l.entries) < l.limit {
			return
		}
		delete(l.entries, key)
	}
}
