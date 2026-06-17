package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCreateMetricHandlerDuplicateName(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "dup", MetricType: "counter", Unit: "c", CurrentValue: float64Ptr(0)}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	body, _ := json.Marshal(timelogv1.CreateMetricRequest{Name: "dup", MetricType: "counter", Unit: "c"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/metrics", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateMetricHandlerDuplicateName(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "a", MetricType: "counter", Unit: "c", CurrentValue: float64Ptr(0)}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := svc.CreateMetric(&domain.Metric{Name: "b", MetricType: "counter", Unit: "c", CurrentValue: float64Ptr(0)}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	name := "a"
	body, _ := json.Marshal(timelogv1.UpdateMetricRequest{Name: &name})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/metrics/2", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteMetricHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/metrics/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
