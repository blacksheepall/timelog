package middleware

import (
	"net/http"

	"github.com/blacksheepaul/timelog/core/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the HTTP header used to carry the request/trace ID.
	RequestIDHeader = requestid.Header
)

// RequestID returns a middleware that ensures every request has an ID.
// It first checks the incoming X-Request-ID header; if absent, it generates a new UUID.
// The ID is stored in the gin context (and propagated to the request context),
// and echoed back in the response header so clients can correlate logs.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Set("request_id", reqID)
		c.Request = c.Request.WithContext(requestid.WithContext(c.Request.Context(), reqID))
		c.Header(RequestIDHeader, reqID)

		c.Next()
	}
}

// GetRequestID extracts the request ID from a gin context.
// It returns an empty string when no request ID has been set.
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RequestIDFromRequest extracts the request ID from an http.Request.
func RequestIDFromRequest(r *http.Request) string {
	return requestid.FromRequest(r)
}
