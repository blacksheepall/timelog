package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestID_GeneratesWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, GetRequestID(c))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Fatal("expected generated request ID in body")
	}

	headerID := w.Header().Get(RequestIDHeader)
	if headerID == "" {
		t.Fatal("expected response header " + RequestIDHeader)
	}
	if headerID != body {
		t.Fatalf("response header %s (%s) does not match body (%s)", RequestIDHeader, headerID, body)
	}
}

func TestRequestID_PreservesIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, GetRequestID(c))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	incomingID := "trace-12345"
	req.Header.Set(RequestIDHeader, incomingID)
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if body != incomingID {
		t.Fatalf("expected incoming request ID %q, got %q", incomingID, body)
	}

	if got := w.Header().Get(RequestIDHeader); got != incomingID {
		t.Fatalf("expected response header %q, got %q", incomingID, got)
	}
}

func TestRequestID_PropagatesToRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		id := RequestIDFromRequest(c.Request)
		c.String(http.StatusOK, id)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Body.String() == "" {
		t.Fatal("expected request ID in request context")
	}
}
