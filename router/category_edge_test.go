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

func TestListCategoriesHandlerInvalidParentID(t *testing.T) {
	r, _ := setupTestRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories?parent_id=bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateCategoryHandlerValidationError(t *testing.T) {
	r, _ := setupTestRouter(t)
	level := int32(1)
	body, _ := json.Marshal(timelogv1.CreateCategoryRequest{Name: "x", Level: &level})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCategoryHandlerDuplicateName(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "dup"}); err != nil {
		t.Fatalf("seed parent: %v", err)
	}
	parentID := int32(1)
	if err := svc.CreateCategory(&domain.Category{Name: "dup", ParentID: &parentID}); err != nil {
		t.Fatalf("seed child: %v", err)
	}

	body, _ := json.Marshal(timelogv1.CreateCategoryRequest{Name: "dup", ParentId: &parentID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCategoryHandlerInvalidParent(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "up"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	pid := int32(1)
	body, _ := json.Marshal(timelogv1.UpdateCategoryRequest{ParentId: &pid})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMoveCategoryHandlerInvalidRequest(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "mv"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	pid := int32(-1)
	body, _ := json.Marshal(timelogv1.MoveCategoryRequest{ParentId: &pid})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories/1/move", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMoveCategoryHandlerParentNotFound(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "mv"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	pid := int32(9999)
	body, _ := json.Marshal(timelogv1.MoveCategoryRequest{ParentId: &pid})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories/1/move", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
