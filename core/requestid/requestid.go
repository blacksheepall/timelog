// Package requestid provides utilities for propagating a request/trace ID
// through HTTP requests and Go contexts. It lives in core so that both loggers
// and HTTP middleware can depend on it without import cycles.
package requestid

import (
	"context"
	"net/http"
)

const (
	// Header is the HTTP header used to carry the request/trace ID.
	Header = "X-Request-ID"
	key    = "request_id"
)

// FromContext extracts the request ID from a standard context.
func FromContext(ctx context.Context) string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// WithContext returns a new context with the given request ID attached.
func WithContext(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, key, reqID)
}

// FromRequest extracts the request ID from an http.Request.
func FromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	return FromContext(r.Context())
}
