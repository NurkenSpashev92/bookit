package middleware

import (
	"bytes"
	"context"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/pkg/cache"
)

const maxCacheableQueryLen = 512

type ResponseStore interface {
	GetBytes(ctx context.Context, key string) ([]byte, bool)
	SetBytesTTL(ctx context.Context, key string, data []byte, ttl time.Duration)
	Enabled() bool
}

type CachedRoute struct {
	Namespace string
	TTL       time.Duration
}

type ResponseCacheConfig struct {
	Cache  ResponseStore
	Routes map[string]CachedRoute
}

func ResponseCache(cfg ResponseCacheConfig) fiber.Handler {
	if cfg.Cache == nil || !cfg.Cache.Enabled() || len(cfg.Routes) == 0 {
		return func(c fiber.Ctx) error { return c.Next() }
	}

	return func(c fiber.Ctx) error {
		route, ok := cfg.Routes[c.Path()]
		if !ok || c.Method() != fiber.MethodGet || isAuthenticated(c) {
			return c.Next()
		}

		query := string(c.Request().URI().QueryString())
		if len(query) > maxCacheableQueryLen {
			return c.Next()
		}

		key := responseCacheKey(route.Namespace, c.Path(), query, acceptedEncodings(c))

		if data, found := cfg.Cache.GetBytes(c.Context(), key); found {
			if entry, valid := decodeResponse(data); valid {
				return entry.send(c)
			}
		}

		if err := c.Next(); err != nil {
			return err
		}

		if storable(c) {
			entry := cachedResponse{
				contentType:     c.GetRespHeader(fiber.HeaderContentType),
				contentEncoding: c.GetRespHeader(fiber.HeaderContentEncoding),
				etag:            c.GetRespHeader(fiber.HeaderETag),
				vary:            c.GetRespHeader(fiber.HeaderVary),
				body:            c.Response().Body(),
			}
			cfg.Cache.SetBytesTTL(c.Context(), key, entry.encode(), route.TTL)
		}

		return nil
	}
}

func isAuthenticated(c fiber.Ctx) bool {
	return c.Cookies("access_token") != "" || c.Cookies("jwt") != ""
}

func storable(c fiber.Ctx) bool {
	return c.Response().StatusCode() == fiber.StatusOK &&
		len(c.Response().Body()) > 0 &&
		c.GetRespHeader(fiber.HeaderSetCookie) == ""
}

func responseCacheKey(ns, path, query, encodings string) string {
	var b strings.Builder
	b.Grow(len(cache.ResponsePrefix) + len(ns) + len(path) + len(query) + len(encodings) + 4)
	b.WriteString(cache.ResponsePrefix)
	b.WriteString(ns)
	b.WriteByte(':')
	b.WriteString(encodings)
	b.WriteByte(':')
	b.WriteString(path)
	b.WriteByte('?')
	b.WriteString(canonicalQuery(query))
	return b.String()
}

func canonicalQuery(query string) string {
	if query == "" || !strings.Contains(query, "&") {
		return query
	}
	parts := strings.Split(query, "&")
	sort.Strings(parts)
	return strings.Join(parts, "&")
}

func acceptedEncodings(c fiber.Ctx) string {
	header := c.Get(fiber.HeaderAcceptEncoding)
	if header == "" {
		return "identity"
	}

	var accepted []string
	for _, candidate := range []string{"br", "deflate", "gzip", "zstd"} {
		if strings.Contains(header, candidate) {
			accepted = append(accepted, candidate)
		}
	}
	if len(accepted) == 0 {
		return "identity"
	}
	return strings.Join(accepted, ".")
}

type cachedResponse struct {
	contentType     string
	contentEncoding string
	etag            string
	vary            string
	body            []byte
}

const cachedResponseFields = 4

func (r cachedResponse) encode() []byte {
	buf := make([]byte, 0, len(r.contentType)+len(r.contentEncoding)+len(r.etag)+len(r.vary)+len(r.body)+cachedResponseFields)
	for _, header := range []string{r.contentType, r.contentEncoding, r.etag, r.vary} {
		buf = append(buf, header...)
		buf = append(buf, '\n')
	}
	return append(buf, r.body...)
}

func decodeResponse(data []byte) (cachedResponse, bool) {
	var (
		entry   cachedResponse
		headers [cachedResponseFields]string
	)

	for i := range cachedResponseFields {
		idx := bytes.IndexByte(data, '\n')
		if idx < 0 {
			return entry, false
		}
		headers[i] = string(data[:idx])
		data = data[idx+1:]
	}

	entry.contentType = headers[0]
	entry.contentEncoding = headers[1]
	entry.etag = headers[2]
	entry.vary = headers[3]
	entry.body = data

	return entry, true
}

func (r cachedResponse) send(c fiber.Ctx) error {
	if r.contentType != "" {
		c.Set(fiber.HeaderContentType, r.contentType)
	}
	if r.contentEncoding != "" {
		c.Set(fiber.HeaderContentEncoding, r.contentEncoding)
	}
	if r.vary != "" {
		c.Set(fiber.HeaderVary, r.vary)
	}

	if r.etag != "" {
		c.Set(fiber.HeaderETag, r.etag)

		if match := c.Get(fiber.HeaderIfNoneMatch); match != "" && strings.Contains(match, r.etag) {
			c.RequestCtx().ResetBody()
			return c.SendStatus(fiber.StatusNotModified)
		}
	}

	return c.Status(fiber.StatusOK).Send(r.body)
}
