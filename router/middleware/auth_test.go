package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

func setupTestService(t *testing.T) *service.Service {
	t.Helper()
	cfg := testutil.NewTestConfig()
	dao, err := model.NewDao(cfg, testutil.FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	repos := adapter.NewRepositories(dao, testutil.FakeLogger{})
	return service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, repos, nil, testutil.FakeLogger{}, cfg, nil)
}

func TestAuthMiddlewareRejectsPasskeySession(t *testing.T) {
	svc := setupTestService(t)

	sessionID := "test-session-abc123"
	svc.WriteCache("passkey_session:"+sessionID, map[string]string{"challenge": "test-challenge"}, 300)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+sessionID)

	Auth(svc)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	svc := setupTestService(t)

	token := "valid-auth-token-xyz789"
	err := svc.StoreSessionToken(token, 300)
	if err != nil {
		t.Fatalf("Failed to store auth token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	Auth(svc)(c)

	if c.IsAborted() {
		t.Errorf("Expected request to pass, but middleware aborted it with status %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsUnprefixedCacheKey(t *testing.T) {
	svc := setupTestService(t)

	token := "unprefixed-token-123"
	svc.WriteCache(token, true, 300)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	Auth(svc)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized for unprefixed key, got %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsMissingHeader(t *testing.T) {
	svc := setupTestService(t)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	Auth(svc)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 for missing header, got %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidHeaderFormat(t *testing.T) {
	svc := setupTestService(t)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "token")

	Auth(svc)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 for invalid header, got %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsWrongScheme(t *testing.T) {
	svc := setupTestService(t)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Basic abc")

	Auth(svc)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 for wrong scheme, got %d", w.Code)
	}
}

func TestAuthMiddlewareWithNilStore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer token")

	Auth(nil)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 for nil store, got %d", w.Code)
	}
}

// TestAuthMiddleware401Envelope pins the req/resp contract for middleware-level
// auth failures: the body must use the same {data, message, status} envelope as
// handler-level errors (see router/apiresponse.go and web/src/types/api.ts).
func TestAuthMiddleware401Envelope(t *testing.T) {
	svc := setupTestService(t)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	Auth(svc)(c)

	if w.Code != 401 {
		t.Fatalf("Expected 401 for missing header, got %d", w.Code)
	}

	var body struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Status  int         `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("401 body is not the ApiResponse envelope: %v (body: %s)", err, w.Body.String())
	}
	if body.Status != 401 {
		t.Errorf("envelope status = %d, want 401", body.Status)
	}
	if body.Message == "" {
		t.Errorf("envelope message must be non-empty (body: %s)", w.Body.String())
	}
	if body.Data != nil {
		t.Errorf("error envelope data must be null, got %v", body.Data)
	}
	if !strings.Contains(w.Body.String(), `"data"`) {
		t.Errorf("envelope must always carry a \"data\" key (body: %s)", w.Body.String())
	}
}
