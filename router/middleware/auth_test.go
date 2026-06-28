package middleware

import (
	"net/http/httptest"
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
	return service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, cfg, nil)
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
