package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCreateTimelogHandlerInvalidCategory(t *testing.T) {
	r, _ := setupTestRouter(t)

	start := time.Now().UTC().Format(time.RFC3339)
	body, _ := json.Marshal(timelogv1.CreateTimelogRequest{StartTime: start, CategoryId: 9999})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/timelogs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTimelogHandlerInvalidStartTime(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{CategoryID: 1, StartTime: time.Now().UTC()}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	bad := "bad"
	body, _ := json.Marshal(timelogv1.UpdateTimelogRequest{StartTime: &bad})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/timelogs/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTimelogHandlerInvalidEndTime(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{CategoryID: 1, StartTime: time.Now().UTC()}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	bad := "bad"
	body, _ := json.Marshal(timelogv1.UpdateTimelogRequest{EndTime: &bad})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/timelogs/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTimelogHandlerInvalidCategory(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{CategoryID: 1, StartTime: time.Now().UTC()}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	catID := int32(9999)
	body, _ := json.Marshal(timelogv1.UpdateTimelogRequest{CategoryId: &catID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/timelogs/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
