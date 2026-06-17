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

func seedTaskCategory(t *testing.T, svc *service.Service) {
	t.Helper()
	if err := svc.CreateCategory(&domain.Category{Name: "tasks"}); err != nil {
		t.Fatalf("seed category: %v", err)
	}
}

func seedTask(t *testing.T, svc *service.Service) {
	t.Helper()
	if err := svc.CreateTask(&domain.Task{Title: "task", CategoryID: 1, DueDate: time.Now().UTC()}); err != nil {
		t.Fatalf("seed task: %v", err)
	}
}

func TestCreateTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)

	due := time.Now().UTC().Format("2006-01-02")
	reqBody, _ := json.Marshal(timelogv1.CreateTaskRequest{Title: "New", CategoryId: 1, DueDate: due})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateTaskHandlerBadJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListTasksHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListTasksHandlerByDate(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks?date="+time.Now().Format("2006-01-02"), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListTasksHandlerInvalidDate(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks?date=bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetTaskHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetTaskHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	newTitle := "Updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateTaskRequest{Title: &newTitle})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTaskHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	newTitle := "Updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateTaskRequest{Title: &newTitle})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/9999", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompleteTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/1/complete", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIncompleteTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/1/incomplete", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSuspendTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/1/suspend", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUnsuspendTaskHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	if err := svc.SuspendTask(1); err != nil {
		t.Fatalf("suspend: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/1/unsuspend", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetTaskStatsHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/stats/"+time.Now().Format("2006-01-02"), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetTaskStatsHandlerInvalidDate(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/stats/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
