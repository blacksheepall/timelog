package service

import (
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
)

func TestPasskeySessionKeyNamespacing(t *testing.T) {
	svc, _ := newTestService(t)

	sessionID := "test-session-123"
	sessionData := &webauthn.SessionData{Challenge: "test-challenge"}

	if err := svc.StorePasskeySession(sessionID, sessionData, 300); err != nil {
		t.Fatalf("StorePasskeySession: %v", err)
	}

	if _, foundRaw := svc.GetCache(sessionID); foundRaw {
		t.Error("raw session ID found in cache without namespace")
	}
	if _, foundNamespaced := svc.GetCache("passkey_session:" + sessionID); !foundNamespaced {
		t.Error("namespaced passkey session not found in cache")
	}
	loaded, err := svc.LoadPasskeySession(sessionID)
	if err != nil || loaded == nil {
		t.Fatalf("LoadPasskeySession: (%v, %v)", loaded, err)
	}
}

func TestAuthTokenKeyNamespacing(t *testing.T) {
	svc, _ := newTestService(t)

	token := "test-auth-token-456"
	if err := svc.StoreSessionToken(token, 300); err != nil {
		t.Fatalf("StoreSessionToken: %v", err)
	}

	if _, foundRaw := svc.GetCache(token); foundRaw {
		t.Error("raw token found in cache without namespace")
	}
	if _, foundNamespaced := svc.GetCache("auth_token:" + token); !foundNamespaced {
		t.Error("namespaced auth token not found in cache")
	}
}

func TestSessionAndTokenIsolation(t *testing.T) {
	svc, _ := newTestService(t)

	sharedID := "shared-id-789"
	sessionData := &webauthn.SessionData{Challenge: "test-challenge"}
	if err := svc.StorePasskeySession(sharedID, sessionData, 300); err != nil {
		t.Fatalf("StorePasskeySession: %v", err)
	}
	if err := svc.StoreSessionToken(sharedID, 300); err != nil {
		t.Fatalf("StoreSessionToken: %v", err)
	}

	passkeyVal, passkeyFound := svc.GetCache("passkey_session:" + sharedID)
	tokenVal, tokenFound := svc.GetCache("auth_token:" + sharedID)
	if !passkeyFound || !tokenFound {
		t.Fatalf("expected both entries, got passkey=%v token=%v", passkeyFound, tokenFound)
	}
	if passkeyVal == tokenVal {
		t.Error("passkey session and auth token should be different")
	}
}
