package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
)

func TestCreateConstraintHandlerInvalidDate(t *testing.T) {
	r, _ := setupTestRouter(t)

	body, _ := json.Marshal(timelogv1.CreateConstraintRequest{Description: "d", PunishmentQuote: "p", StartDate: "bad"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateConstraintHandlerInvalidDate(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	badDate := "bad"
	body, _ := json.Marshal(timelogv1.UpdateConstraintRequest{StartDate: &badDate})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/constraints/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompleteConstraintHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(timelogv1.CompleteConstraintRequest{EndReason: "done"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints/bad/complete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestReactivateConstraintHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints/bad/reactivate", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEvaluateConstraintHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/9999/evaluation", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEvaluateConstraintHandlerMissingMetricRule(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/1/evaluation", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
