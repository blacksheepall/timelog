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

// setupTestEnvironment initializes model and service for testing
func setupTestEnvironment() {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	model.InitDao(cfg, FakeLogger{})
	service.InitService(FakeLogger{}, cfg)
	setMiddlewareDAOProviderForTest(model.GetDao)
}

func TestAuthMiddlewareRejectsPasskeySession(t *testing.T) {
	setupTestEnvironment()

	// Create a passkey session (stored with passkey_session: prefix)
	sessionID := "test-session-abc123"

	// Simulate storing a passkey session directly in cache
	dao := model.GetDao()
	dao.WriteCache("passkey_session:"+sessionID, map[string]string{"challenge": "test-challenge"}, 300)

	// Try to use the sessionID as an auth token
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+sessionID)

	// Call the auth middleware
	Auth()(c)

	// Should be rejected (401)
	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}

	t.Log("✓ Auth middleware correctly rejects passkey session IDs")
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	setupTestEnvironment()

	// Create a valid auth token (stored with auth_token: prefix)
	token := "valid-auth-token-xyz789"
	err := service.StoreSessionToken(token, 300)
	if err != nil {
		t.Fatalf("Failed to store auth token: %v", err)
	}

	// Use the token in a request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Create a handler that sets status 200 if middleware passes
	Auth()(c)

	if c.IsAborted() {
		t.Errorf("Expected request to pass, but middleware aborted it with status %d", w.Code)
	}

	t.Log("✓ Auth middleware correctly accepts valid auth tokens")
}

func TestAuthMiddlewareRejectsUnprefixedCacheKey(t *testing.T) {
	setupTestEnvironment()

	// Simulate an old-style cache entry without prefix (should not work)
	token := "unprefixed-token-123"
	dao := model.GetDao()
	dao.WriteCache(token, true, 300) // Store without prefix

	// Try to use it
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Call the auth middleware
	Auth()(c)

	// Should be rejected (401) because middleware only accepts auth_token: prefix
	if w.Code != 401 {
		t.Errorf("Expected 401 Unauthorized for unprefixed key, got %d", w.Code)
	}

	t.Log("✓ Auth middleware correctly rejects unprefixed cache keys")
}
