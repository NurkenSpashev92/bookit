package logger

import (
	"context"
	"log/slog"
)

const (
	KeyRequestID = "request_id"
	KeyError     = "error"
)

type contextKey struct{}

var loggerKey contextKey

func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	if log == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	if log, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	id, _ := ctx.Value(requestIDKey).(string)

	return id
}

func WithRequestID(ctx context.Context, id string) context.Context {
	if id == "" {
		return ctx
	}

	return context.WithValue(ctx, requestIDKey, id)
}

type requestIDContextKey struct{}

var requestIDKey requestIDContextKey

func Err(err error) slog.Attr {
	if err == nil {
		return slog.String(KeyError, "")
	}

	return slog.String(KeyError, err.Error())
}
