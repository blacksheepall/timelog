package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestEnvelopeContract pins the req/resp contract that the frontend relies on
// (web/src/types/api.ts + err.response?.data?.message on non-2xx).
// Every API response — success or error, from handlers or from router-level
// fallbacks — must use the {data, message, status} envelope.
func TestEnvelopeContract(t *testing.T) {
	r, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		// success paths (200)
		{name: "list categories", method: http.MethodGet, path: "/api/categories", wantStatus: http.StatusOK},
		{name: "list timelogs", method: http.MethodGet, path: "/api/timelogs", wantStatus: http.StatusOK},
		{name: "list tasks", method: http.MethodGet, path: "/api/tasks", wantStatus: http.StatusOK},
		{name: "list metrics", method: http.MethodGet, path: "/api/metrics", wantStatus: http.StatusOK},
		{name: "list constraints", method: http.MethodGet, path: "/api/constraints", wantStatus: http.StatusOK},

		// bind/validation failures (400)
		{name: "create category bad json", method: http.MethodPost, path: "/api/categories", body: `{invalid`, wantStatus: http.StatusBadRequest},
		{name: "create timelog bad json", method: http.MethodPost, path: "/api/timelogs", body: `{invalid`, wantStatus: http.StatusBadRequest},
		{name: "create task bad json", method: http.MethodPost, path: "/api/tasks", body: `{invalid`, wantStatus: http.StatusBadRequest},

		// not found (404)
		{name: "get unknown category", method: http.MethodGet, path: "/api/categories/999999", wantStatus: http.StatusNotFound},
		{name: "get unknown timelog", method: http.MethodGet, path: "/api/timelogs/999999", wantStatus: http.StatusNotFound},

		// router-level fallback: unknown API endpoint must also use the envelope
		{name: "unknown api route", method: http.MethodGet, path: "/api/does-not-exist", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assertEnvelope(t, w, tt.wantStatus)
		})
	}
}
