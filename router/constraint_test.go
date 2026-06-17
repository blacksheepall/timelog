package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/service"
)

func seedConstraintForRouter(t *testing.T, svc *service.Service) {
	t.Helper()
	c := &domain.Constraint{
		Description:     "focus",
		PunishmentQuote: "punish",
		StartDate:       time.Now().UTC(),
		IsActive:        true,
	}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("seed constraint: %v", err)
	}
}

func TestCreateConstraintHandler(t *testing.T) {
	r, _ := setupTestRouter(t)

	start := time.Now().UTC().Format("2006-01-02")
	reqBody, _ := json.Marshal(timelogv1.CreateConstraintRequest{Description: "focus", PunishmentQuote: "punish", StartDate: start})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateConstraintHandlerBadJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListConstraintsHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListConstraintsHandlerActive(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints?active=true", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetConstraintHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetConstraintHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	newDesc := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateConstraintRequest{Description: &newDesc})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/constraints/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateConstraintHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	newDesc := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateConstraintRequest{Description: &newDesc})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/constraints/9999", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/constraints/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompleteConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	reqBody, _ := json.Marshal(timelogv1.CompleteConstraintRequest{EndReason: "done"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints/1/complete", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestReactivateConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)
	if err := svc.MarkConstraintAsCompleted(1, "done"); err != nil {
		t.Fatalf("complete: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/constraints/1/reactivate", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEvaluateConstraintHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedConstraintForRouter(t, svc)

	// create metric and link constraint to it
	if err := svc.CreateMetric(&domain.Metric{Name: "pages", MetricType: "counter", Unit: "count", CurrentValue: float64Ptr(10)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}
	c, err := svc.GetConstraintByID(1)
	if err != nil {
		t.Fatalf("get constraint: %v", err)
	}
	c.MetricID = int32Ptr(1)
	c.MetricOperator = strPtr("gte")
	c.MetricTargetValue = float64Ptr(5)
	if err := svc.UpdateConstraint(c); err != nil {
		t.Fatalf("update constraint: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/constraints/1/evaluation", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
