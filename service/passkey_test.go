package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
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

// setupTestModel initializes the model with an in-memory database for testing
func setupTestModel() *model.Dao {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		panic(err)
	}
	repos := adapter.NewRepositories(dao)
	svc := NewService(repos, repos, repos, repos, repos, repos, repos, repos, FakeLogger{}, cfg, nil)
	SetDefaultService(svc)
	return dao
}

func TestPasskeySessionKeyNamespacing(t *testing.T) {
	dao := setupTestModel()

	// Test 1: Verify passkey session uses namespaced key
	sessionID := "test-session-123"
	sessionData := &webauthn.SessionData{
		Challenge: "test-challenge",
	}

	err := StorePasskeySession(sessionID, sessionData, 300)
	if err != nil {
		t.Fatalf("Failed to store passkey session: %v", err)
	}

	// The raw sessionID should NOT be in cache
	_, foundRaw := dao.GetCache(sessionID)
	if foundRaw {
		t.Error("Security vulnerability: raw session ID found in cache without namespace")
	}

	// The namespaced key should be in cache
	_, foundNamespaced := dao.GetCache("passkey_session:" + sessionID)
	if !foundNamespaced {
		t.Error("Namespaced passkey session not found in cache")
	}

	// LoadPasskeySession should work with the sessionID (it adds the namespace internally)
	loadedSession, err := LoadPasskeySession(sessionID)
	if err != nil {
		t.Errorf("Failed to load passkey session: %v", err)
	}
	if loadedSession == nil {
		t.Error("Loaded session is nil")
	}

	t.Log("✓ Passkey sessions are correctly namespaced with 'passkey_session:' prefix")
}

func TestAuthTokenKeyNamespacing(t *testing.T) {
	dao := setupTestModel()

	// Test 2: Verify auth token uses namespaced key
	token := "test-auth-token-456"

	err := StoreSessionToken(token, 300)
	if err != nil {
		t.Fatalf("Failed to store auth token: %v", err)
	}

	// The raw token should NOT be in cache
	_, foundRaw := dao.GetCache(token)
	if foundRaw {
		t.Error("Security vulnerability: raw token found in cache without namespace")
	}

	// The namespaced key should be in cache
	_, foundNamespaced := dao.GetCache("auth_token:" + token)
	if !foundNamespaced {
		t.Error("Namespaced auth token not found in cache")
	}

	t.Log("✓ Auth tokens are correctly namespaced with 'auth_token:' prefix")
}

func TestSessionAndTokenIsolation(t *testing.T) {
	dao := setupTestModel()

	// Test 3: Verify that using the same ID for both doesn't cause collision
	sharedID := "shared-id-789"

	// Store a passkey session
	sessionData := &webauthn.SessionData{
		Challenge: "test-challenge",
	}
	err := StorePasskeySession(sharedID, sessionData, 300)
	if err != nil {
		t.Fatalf("Failed to store passkey session: %v", err)
	}

	// Store an auth token with the same ID
	err = StoreSessionToken(sharedID, 300)
	if err != nil {
		t.Fatalf("Failed to store auth token: %v", err)
	}

	// Both should exist independently
	passkeyVal, passkeyFound := dao.GetCache("passkey_session:" + sharedID)
	tokenVal, tokenFound := dao.GetCache("auth_token:" + sharedID)

	if !passkeyFound {
		t.Error("Passkey session not found")
	}
	if !tokenFound {
		t.Error("Auth token not found")
	}

	// They should be different values
	if passkeyVal == tokenVal {
		t.Error("Passkey session and auth token should be different but are the same")
	}

	t.Log("✓ Passkey sessions and auth tokens are properly isolated even with same ID")
}
