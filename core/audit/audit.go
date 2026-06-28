// Package audit provides source constants and a Logger interface for recording
// auditable changes to entities. It is intentionally dependency-free so that
// both HTTP handlers and MCP handlers can tag requests with a source before
// calling into the service layer.
package audit

import (
	"context"

	"github.com/blacksheepaul/timelog/core/requestid"
)

// Source values identify who/what initiated a mutating operation.
const (
	SourceHuman = "human"
	SourceAI    = "ai"
)

// Action values describe the type of mutation performed.
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Entity values identify the domain object that changed.
const (
	EntityTimeLog = "timelog"
)

// Logger records an auditable event both durably (database) and to the
// application log for easy searching.
type Logger interface {
	Log(ctx context.Context, action, entityType string, entityID int32, payload map[string]any) error
}

type sourceKey struct{}

// WithSource returns a new context annotated with the given audit source.
func WithSource(ctx context.Context, source string) context.Context {
	return context.WithValue(ctx, sourceKey{}, source)
}

// SourceFromContext extracts the audit source from a context.
// It returns an empty string when no source has been set.
func SourceFromContext(ctx context.Context) string {
	if v := ctx.Value(sourceKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RequestIDFromContext extracts the request/trace ID from the context if present.
func RequestIDFromContext(ctx context.Context) string {
	return requestid.FromContext(ctx)
}
