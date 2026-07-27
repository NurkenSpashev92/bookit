package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPath = "/api/v1/houses"

type memoryStore struct {
	mu    sync.Mutex
	items map[string][]byte
}

func newMemoryStore() *memoryStore {
	return &memoryStore{items: map[string][]byte{}}
}

func (s *memoryStore) GetBytes(_ context.Context, key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.items[key]
	return data, ok
}

func (s *memoryStore) SetBytes(_ context.Context, key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = append([]byte(nil), data...)
}

func (s *memoryStore) Enabled() bool { return true }

func (s *memoryStore) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}

func newCachedApp(store ResponseStore, handler fiber.Handler) (*fiber.App, *int) {
	calls := 0

	app := fiber.New(fiber.Config{CaseSensitive: true})
	app.Use(ResponseCache(ResponseCacheConfig{
		Cache:  store,
		Routes: map[string]string{testPath: "houses"},
	}))
	app.Get(testPath, func(c fiber.Ctx) error {
		calls++
		return handler(c)
	})

	return app, &calls
}

func okHandler(body string) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.SendString(body)
	}
}

func get(t *testing.T, app *fiber.App, target string, headers map[string]string, cookies ...*http.Cookie) *http.Response {
	t.Helper()

	req := httptest.NewRequest(fiber.MethodGet, target, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	require.NoError(t, err)
	return resp
}

func TestResponseCache_ServesSecondRequestFromCache(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	first := get(t, app, testPath+"?page=1", nil)
	assert.Equal(t, http.StatusOK, first.StatusCode)

	second := get(t, app, testPath+"?page=1", nil)
	assert.Equal(t, http.StatusOK, second.StatusCode)

	assert.Equal(t, 1, *calls, "handler must run only on the miss")
	assert.Equal(t, 1, store.len())
}

func TestResponseCache_BypassesAuthenticatedRequests(t *testing.T) {
	tests := []struct {
		name   string
		cookie *http.Cookie
	}{
		{name: "access_token", cookie: &http.Cookie{Name: "access_token", Value: "token"}},
		{name: "jwt", cookie: &http.Cookie{Name: "jwt", Value: "token"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMemoryStore()
			app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

			get(t, app, testPath, nil, tt.cookie)
			get(t, app, testPath, nil, tt.cookie)

			assert.Equal(t, 2, *calls, "per-user responses must not be shared")
			assert.Equal(t, 0, store.len(), "nothing may be stored for a logged-in user")
		})
	}
}

func TestResponseCache_AnonymousEntryNotServedToAuthenticated(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	get(t, app, testPath, nil)
	assert.Equal(t, 1, *calls)

	get(t, app, testPath, nil, &http.Cookie{Name: "access_token", Value: "token"})
	assert.Equal(t, 2, *calls, "logged-in request must re-run the handler")
}

func TestResponseCache_KeyIsQueryOrderIndependent(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	get(t, app, testPath+"?page=1&min_price=100", nil)
	get(t, app, testPath+"?min_price=100&page=1", nil)

	assert.Equal(t, 1, *calls, "same filter in a different order is the same entry")
	assert.Equal(t, 1, store.len())
}

func TestResponseCache_KeySeparatesDifferentFilters(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	get(t, app, testPath+"?page=1", nil)
	get(t, app, testPath+"?page=2", nil)

	assert.Equal(t, 2, *calls)
	assert.Equal(t, 2, store.len())
}

func TestResponseCache_KeySeparatesByAcceptedEncodings(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	get(t, app, testPath, map[string]string{fiber.HeaderAcceptEncoding: "gzip"})
	get(t, app, testPath, map[string]string{fiber.HeaderAcceptEncoding: "br"})
	get(t, app, testPath, nil)

	assert.Equal(t, 3, *calls)
	assert.Equal(t, 3, store.len())
}

func TestResponseCache_ReplaysHeaders(t *testing.T) {
	store := newMemoryStore()
	app, _ := newCachedApp(store, func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		c.Set(fiber.HeaderETag, `"abc-123"`)
		c.Set(fiber.HeaderVary, fiber.HeaderAcceptEncoding)
		return c.SendString(`{"data":[]}`)
	})

	get(t, app, testPath, nil)
	cached := get(t, app, testPath, nil)

	assert.Equal(t, fiber.MIMEApplicationJSON, cached.Header.Get(fiber.HeaderContentType))
	assert.Equal(t, `"abc-123"`, cached.Header.Get(fiber.HeaderETag))
	assert.Equal(t, fiber.HeaderAcceptEncoding, cached.Header.Get(fiber.HeaderVary))
}

func TestResponseCache_AnswersNotModifiedForMatchingETag(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, func(c fiber.Ctx) error {
		c.Set(fiber.HeaderETag, `"abc-123"`)
		return c.SendString(`{"data":[]}`)
	})

	get(t, app, testPath, nil)

	resp := get(t, app, testPath, map[string]string{fiber.HeaderIfNoneMatch: `"abc-123"`})
	assert.Equal(t, http.StatusNotModified, resp.StatusCode)

	stale := get(t, app, testPath, map[string]string{fiber.HeaderIfNoneMatch: `"other"`})
	assert.Equal(t, http.StatusOK, stale.StatusCode)

	assert.Equal(t, 1, *calls, "both follow-ups are answered from the cache")
}

func TestResponseCache_SkipsUncacheableResponses(t *testing.T) {
	tests := []struct {
		name    string
		handler fiber.Handler
	}{
		{
			name: "non-200",
			handler: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusInternalServerError).SendString(`{"error":"boom"}`)
			},
		},
		{
			name:    "empty body",
			handler: func(c fiber.Ctx) error { return nil },
		},
		{
			name: "sets a cookie",
			handler: func(c fiber.Ctx) error {
				c.Cookie(&fiber.Cookie{Name: "session", Value: "s"})
				return c.SendString(`{"data":[]}`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMemoryStore()
			app, _ := newCachedApp(store, tt.handler)

			get(t, app, testPath, nil)

			assert.Equal(t, 0, store.len())
		})
	}
}

func TestResponseCache_IgnoresUnlistedRoutesAndMethods(t *testing.T) {
	store := newMemoryStore()

	app := fiber.New(fiber.Config{CaseSensitive: true})
	app.Use(ResponseCache(ResponseCacheConfig{
		Cache:  store,
		Routes: map[string]string{testPath: "houses"},
	}))
	app.Get("/api/v1/other", okHandler(`{"data":[]}`))
	app.Post(testPath, okHandler(`{"created":true}`))

	get(t, app, "/api/v1/other", nil)

	req := httptest.NewRequest(fiber.MethodPost, testPath, nil)
	_, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	require.NoError(t, err)

	assert.Equal(t, 0, store.len(), "only listed GET routes are cacheable")
}

func TestResponseCache_SkipsOversizedQuery(t *testing.T) {
	store := newMemoryStore()
	app, calls := newCachedApp(store, okHandler(`{"data":[]}`))

	long := testPath + "?junk=" + strings.Repeat("x", maxCacheableQueryLen)
	get(t, app, long, nil)
	get(t, app, long, nil)

	assert.Equal(t, 2, *calls)
	assert.Equal(t, 0, store.len(), "a long query must not mint an entry")
}

func TestCachedResponse_EncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry cachedResponse
	}{
		{
			name: "all headers set",
			entry: cachedResponse{
				contentType:     fiber.MIMEApplicationJSON,
				contentEncoding: "gzip",
				etag:            `"abc-123"`,
				vary:            fiber.HeaderAcceptEncoding,
				body:            []byte(`{"data":[]}`),
			},
		},
		{
			name:  "headers empty",
			entry: cachedResponse{body: []byte("plain")},
		},
		{
			name: "binary body with newlines",
			entry: cachedResponse{
				contentEncoding: "gzip",
				body:            []byte{0x1f, 0x8b, '\n', 0x00, '\n', '\n', 0xff},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, ok := decodeResponse(tt.entry.encode())

			require.True(t, ok)
			assert.Equal(t, tt.entry.contentType, decoded.contentType)
			assert.Equal(t, tt.entry.contentEncoding, decoded.contentEncoding)
			assert.Equal(t, tt.entry.etag, decoded.etag)
			assert.Equal(t, tt.entry.vary, decoded.vary)
			assert.Equal(t, tt.entry.body, decoded.body)
		})
	}
}

func TestDecodeResponse_RejectsTruncatedEntry(t *testing.T) {
	_, ok := decodeResponse([]byte("application/json\ngzip\n"))
	assert.False(t, ok, "an entry missing header lines must not be served")
}

func TestAcceptedEncodings(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "absent", header: "", want: "identity"},
		{name: "unknown only", header: "compress", want: "identity"},
		{name: "gzip", header: "gzip", want: "gzip"},
		{name: "sorted regardless of order", header: "gzip, br", want: "br.gzip"},
		{name: "same set same token", header: "br, gzip", want: "br.gzip"},
		{name: "with quality values", header: "gzip;q=1.0, deflate;q=0.5", want: "deflate.gzip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			var got string
			app.Get("/", func(c fiber.Ctx) error {
				got = acceptedEncodings(c)
				return nil
			})

			req := httptest.NewRequest(fiber.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set(fiber.HeaderAcceptEncoding, tt.header)
			}
			_, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCanonicalQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{name: "empty", query: "", want: ""},
		{name: "single param untouched", query: "page=1", want: "page=1"},
		{name: "sorted", query: "page=1&min_price=100", want: "min_price=100&page=1"},
		{name: "already sorted", query: "min_price=100&page=1", want: "min_price=100&page=1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, canonicalQuery(tt.query))
		})
	}
}
