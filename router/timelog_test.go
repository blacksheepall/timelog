package router

import (
	"context"
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

func seedTimelogCategory(t *testing.T, svc *service.Service) {
	t.Helper()
	if err := svc.CreateCategory(&domain.Category{Name: "timelog"}); err != nil {
		t.Fatalf("seed category: %v", err)
	}
}

func TestCreateTimelogHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)

	start := time.Now().UTC().Format(time.RFC3339)
	reqBody, _ := json.Marshal(timelogv1.CreateTimelogRequest{StartTime: start, CategoryId: 1})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/timelogs", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateTimelogHandlerBadJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/timelogs", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListTimelogsHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("seed timelog: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/timelogs", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListTimelogsHandlerWithOptions(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)
	if err := svc.CreateTimeLog(context.Background(), &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("seed timelog: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/timelogs?limit=10&order=start_time+DESC", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetTimelogHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)
	tl := &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := svc.CreateTimeLog(context.Background(), tl); err != nil {
		t.Fatalf("seed timelog: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/timelogs/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetTimelogHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/timelogs/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTimelogHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)
	tl := &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := svc.CreateTimeLog(context.Background(), tl); err != nil {
		t.Fatalf("seed timelog: %v", err)
	}

	remark := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateTimelogRequest{Remark: &remark})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/timelogs/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTimelogHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	remark := "updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateTimelogRequest{Remark: &remark})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/timelogs/9999", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTimelogHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTimelogCategory(t, svc)
	tl := &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := svc.CreateTimeLog(context.Background(), tl); err != nil {
		t.Fatalf("seed timelog: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/timelogs/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
