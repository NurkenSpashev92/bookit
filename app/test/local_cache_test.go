package test

import (
	"testing"
	"time"

	"github.com/nurkenspashev92/bookit/pkg/cache"
)

func TestLocalCacheReturnsStoredValue(t *testing.T) {
	local := cache.NewLocalCache(time.Minute, 10)

	local.Set("resp:houses:key", []byte("payload"), 0)

	data, found := local.Get("resp:houses:key")
	if !found {
		t.Fatal("value must be found before it expires")
	}
	if string(data) != "payload" {
		t.Errorf("got %q, want %q", data, "payload")
	}
}

func TestLocalCacheExpires(t *testing.T) {
	local := cache.NewLocalCache(time.Minute, 10)

	local.Set("key", []byte("payload"), time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	if _, found := local.Get("key"); found {
		t.Error("expired value must not be returned")
	}
}

func TestLocalCacheCapsEntryTTLAtConfiguredTTL(t *testing.T) {
	local := cache.NewLocalCache(time.Millisecond, 10)

	local.Set("key", []byte("payload"), time.Hour)
	time.Sleep(5 * time.Millisecond)

	if _, found := local.Get("key"); found {
		t.Error("entry ttl must not exceed the cache ttl")
	}
}

func TestLocalCacheDeletePrefix(t *testing.T) {
	local := cache.NewLocalCache(time.Minute, 10)

	local.Set("resp:houses:one", []byte("a"), 0)
	local.Set("resp:houses:two", []byte("b"), 0)
	local.Set("resp:cities:one", []byte("c"), 0)

	local.DeletePrefix("resp:houses:")

	if _, found := local.Get("resp:houses:one"); found {
		t.Error("prefixed key must be dropped")
	}
	if _, found := local.Get("resp:houses:two"); found {
		t.Error("prefixed key must be dropped")
	}
	if _, found := local.Get("resp:cities:one"); !found {
		t.Error("other namespaces must survive")
	}
}

func TestLocalCacheRespectsLimit(t *testing.T) {
	const limit = 4

	local := cache.NewLocalCache(time.Minute, limit)

	for _, key := range []string{"a", "b", "c", "d", "e", "f"} {
		local.Set(key, []byte(key), 0)
	}

	stored := 0
	for _, key := range []string{"a", "b", "c", "d", "e", "f"} {
		if _, found := local.Get(key); found {
			stored++
		}
	}

	if stored > limit {
		t.Errorf("cache holds %d entries, limit is %d", stored, limit)
	}
	if stored == 0 {
		t.Error("cache must keep entries after eviction")
	}
}

func TestLocalCacheDisabled(t *testing.T) {
	local := cache.NewLocalCache(0, 10)

	local.Set("key", []byte("payload"), 0)

	if _, found := local.Get("key"); found {
		t.Error("disabled cache must not store values")
	}
}

func TestLocalCacheIgnoresEmptyPayload(t *testing.T) {
	local := cache.NewLocalCache(time.Minute, 10)

	local.Set("key", nil, 0)

	if _, found := local.Get("key"); found {
		t.Error("empty payload must not be cached")
	}
}
