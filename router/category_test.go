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

func TestCreateCategoryHandler(t *testing.T) {
	r, _ := setupTestRouter(t)

	reqBody, _ := json.Marshal(timelogv1.CreateCategoryRequest{Name: "Work", Color: strPtr("#FF0000")})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp ApiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != http.StatusOK {
		t.Fatalf("expected status 200 in envelope, got %v", resp.Status)
	}
}

func TestCreateCategoryHandlerInvalidParent(t *testing.T) {
	r, _ := setupTestRouter(t)

	parentID := int32(0)
	reqBody, _ := json.Marshal(timelogv1.CreateCategoryRequest{Name: "Bad", ParentId: &parentID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCategoryHandlerBadJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListCategoriesHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListCategoriesHandlerByLevel(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories?level=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListCategoriesHandlerInvalidLevel(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories?level=bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetCategoryHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetCategoryHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/9999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetCategoryHandlerInvalidID(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateCategoryHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)

	newName := "Updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateCategoryRequest{Name: &newName})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCategoryHandlerNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	newName := "Updated"
	reqBody, _ := json.Marshal(timelogv1.UpdateCategoryRequest{Name: &newName})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/categories/9999", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetCategoryTreeHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	seedCategory(t, svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/tree", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMoveCategoryHandler(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "Root"}); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if err := svc.CreateCategory(&domain.Category{Name: "Child", ParentID: int32Ptr(1)}); err != nil {
		t.Fatalf("create child: %v", err)
	}

	newParent := int32(1)
	reqBody, _ := json.Marshal(timelogv1.MoveCategoryRequest{ParentId: &newParent})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories/2/move", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMoveCategoryHandlerInvalidMove(t *testing.T) {
	r, svc := setupTestRouter(t)
	if err := svc.CreateCategory(&domain.Category{Name: "Root"}); err != nil {
		t.Fatalf("create root: %v", err)
	}

	parentID := int32(1)
	reqBody, _ := json.Marshal(timelogv1.MoveCategoryRequest{ParentId: &parentID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categories/1/move", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
