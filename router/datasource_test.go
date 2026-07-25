package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSyncDatasourceHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	_ = svc // future seeding if needed

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/datasources/maimemo/sync", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSyncDatasourceHandler_NotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/datasources/missing/sync", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}
