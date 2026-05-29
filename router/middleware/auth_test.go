package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// FakeLogger for testing
type FakeLogger struct{}

func (l FakeLogger) Debug(fields ...interface{})                     {}
func (l FakeLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (l FakeLogger) Info(fields ...interface{})                      {}
func (l FakeLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (l FakeLogger) Warn(fields ...interface{})                      {}
func (l FakeLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (l FakeLogger) Error(fields ...interface{})                     {}
func (l FakeLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (l FakeLogger) Fatal(fields ...interface{})                     {}
func (l FakeLogger) Fatalw(msg string, keysAndValues ...interface{}) {}

func setupTestDAO() *model.Dao {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		panic(err)
	}
	service.InitService(FakeLogger{}, cfg, dao)
	return dao
}

func TestAuthMiddlewareRejectsPasskeySession(t *testing.T) {
	dao := setupTestDAO()

	sessionID := "test-session-abc123"
	dao.WriteCache("passkey_session:"+sessionID, map[string]string{"challenge": "test-challenge"}, 300)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+sessionID)

	Auth(dao)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	dao := setupTestDAO()

	token := "valid-auth-token-xyz789"
	err := service.StoreSessionToken(token, 300)
	if err != nil {
		t.Fatalf("Failed to store auth token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	Auth(dao)(c)

	if c.IsAborted() {
		t.Errorf("Expected request to pass, but middleware aborted it with status %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsUnprefixedCacheKey(t *testing.T) {
	dao := setupTestDAO()

	token := "unprefixed-token-123"
	dao.WriteCache(token, true, 300)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	Auth(dao)(c)

	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized for unprefixed key, got %d", w.Code)
	}
}
