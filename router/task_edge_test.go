package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
)

func TestCreateTaskHandlerInvalidCategory(t *testing.T) {
	r, _ := setupTestRouter(t)

	body, _ := json.Marshal(timelogv1.CreateTaskRequest{Title: "t", CategoryId: 9999, DueDate: "2026-06-17"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTaskHandlerInvalidDate(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	badDate := "not-a-date"
	body, _ := json.Marshal(timelogv1.UpdateTaskRequest{DueDate: &badDate})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTaskHandlerInvalidCategory(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedTaskCategory(t, svc)
	seedTask(t, svc)

	catID := int32(9999)
	body, _ := json.Marshal(timelogv1.UpdateTaskRequest{CategoryId: &catID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompleteTaskHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/bad/complete", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIncompleteTaskHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/bad/incomplete", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSuspendTaskHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/9999/suspend", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUnsuspendTaskHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/9999/unsuspend", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}
