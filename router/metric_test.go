package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/service"
)

func seedMetric(t *testing.T, svc *service.Service) {
	t.Helper()
	if err := svc.CreateMetric(&domain.Metric{Name: "pages", MetricType: "counter", Unit: "count"}); err != nil {
		t.Fatalf("seed metric: %v", err)
	}
}

func TestCreateMetricHandler(t *testing.T) {
	r, _ := setupTestRouter(t)

	reqBody, _ := json.Marshal(timelogv1.CreateMetricRequest{Name: "pages", MetricType: "counter", Unit: "count"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/metrics", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateMetricHandlerBadJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/metrics", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListMetricsHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedMetric(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetMetricHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedMetric(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetMetricHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetMetricHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateMetricHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedMetric(t, svc)

	newUnit := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateMetricRequest{Unit: &newUnit})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/metrics/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateMetricHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	newUnit := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateMetricRequest{Unit: &newUnit})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/metrics/9999", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteMetricHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedMetric(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/metrics/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListMetricRecordsHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedMetric(t, svc)
	if _, err := svc.RecordMetric(service.RecordMetricInput{MetricName: "pages", Value: 1, Source: "test"}); err != nil {
		t.Fatalf("seed record: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/1/records", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
