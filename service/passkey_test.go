package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/domain"
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

func TestInitWebAuthnWithConfigNil(t *testing.T) {
	if _, err := InitWebAuthnWithConfig(nil); err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestInitWebAuthnWithConfigIncomplete(t *testing.T) {
	cfg := &config.Config{}
	if _, err := InitWebAuthnWithConfig(cfg); err == nil {
		t.Fatal("expected error for incomplete config")
	}
}

func TestInitWebAuthnWithConfigOK(t *testing.T) {
	cfg := &config.Config{}
	cfg.Passkey.RPID = "localhost"
	cfg.Passkey.RPName = "test"
	cfg.Passkey.RPOrigins = []string{"http://localhost:8080"}
	w, err := InitWebAuthnWithConfig(cfg)
	if err != nil {
		t.Fatalf("InitWebAuthnWithConfig: %v", err)
	}
	if w == nil {
		t.Fatal("expected webauthn instance")
	}
}

func TestLoadPasskeySessionMissing(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.LoadPasskeySession("missing"); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadPasskeySessionInvalidType(t *testing.T) {
	svc, _ := newTestService(t)
	svc.WriteCache("passkey_session:bad", "not-a-session", 60)
	if _, err := svc.LoadPasskeySession("bad"); err == nil {
		t.Fatal("expected error")
	}
}

func TestGenerateSessionToken(t *testing.T) {
	svc, _ := newTestService(t)
	token, err := svc.GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
}

func TestTempPasswordLifecycle(t *testing.T) {
	svc, _ := newTestService(t)
	record, password, err := svc.CreateTempPassword(60)
	if err != nil {
		t.Fatalf("CreateTempPassword: %v", err)
	}
	if password == "" {
		t.Fatal("expected password")
	}

	validated, err := svc.ValidateTempPassword(password)
	if err != nil {
		t.Fatalf("ValidateTempPassword: %v", err)
	}
	if validated == nil || validated.PasswordHash != record.PasswordHash {
		t.Fatal("validated record mismatch")
	}

	list, err := svc.ListTempPasswords()
	if err != nil {
		t.Fatalf("ListTempPasswords: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 temp password, got %d", len(list))
	}

	if err := svc.CleanupExpiredTempPasswords(); err != nil {
		t.Fatalf("CleanupExpiredTempPasswords: %v", err)
	}

	if err := svc.DeleteTempPassword(record.ID); err != nil {
		t.Fatalf("DeleteTempPassword: %v", err)
	}
}

func TestCreatePasskeyCredentialAndList(t *testing.T) {
	svc, _ := newTestService(t)
	cred := &webauthn.Credential{
		ID:              []byte("id"),
		PublicKey:       []byte("pk"),
		AttestationType: "none",
	}
	record, err := svc.CreatePasskeyCredential(cred, "device")
	if err != nil {
		t.Fatalf("CreatePasskeyCredential: %v", err)
	}
	if record == nil || record.DeviceName != "device" {
		t.Fatal("unexpected record")
	}

	list, err := svc.ListPasskeyCredentials()
	if err != nil {
		t.Fatalf("ListPasskeyCredentials: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(list))
	}

	loaded, err := svc.LoadPasskeyCredentialByID([]byte("id"))
	if err != nil {
		t.Fatalf("LoadPasskeyCredentialByID: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected credential")
	}

	if err := svc.DeletePasskeyCredential(record.ID); err != nil {
		t.Fatalf("DeletePasskeyCredential: %v", err)
	}
}

func TestCreatePasskeyCredentialNil(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.CreatePasskeyCredential(nil, "device"); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadPasskeyUser(t *testing.T) {
	svc, _ := newTestService(t)
	user, err := svc.LoadPasskeyUser()
	if err != nil {
		t.Fatalf("LoadPasskeyUser: %v", err)
	}
	if user == nil {
		t.Fatal("expected user")
	}
}

func TestLoadPasskeyUserByHandle(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.LoadPasskeyUserByHandle(nil, []byte("timelog-single-user")); err != nil {
		t.Fatalf("LoadPasskeyUserByHandle: %v", err)
	}
	if _, err := svc.LoadPasskeyUserByHandle(nil, []byte("unknown")); err == nil {
		t.Fatal("expected error for unknown handle")
	}
}

func TestUpdatePasskeyCredentialAuth(t *testing.T) {
	svc, _ := newTestService(t)
	cred := &webauthn.Credential{
		ID:              []byte("id"),
		PublicKey:       []byte("pk"),
		AttestationType: "none",
	}
	record, err := svc.CreatePasskeyCredential(cred, "device")
	if err != nil {
		t.Fatalf("CreatePasskeyCredential: %v", err)
	}

	updated := &webauthn.Credential{
		ID: record.CredentialID,
		Authenticator: webauthn.Authenticator{
			SignCount: 5,
		},
	}
	if err := svc.UpdatePasskeyCredentialAuth(updated); err != nil {
		t.Fatalf("UpdatePasskeyCredentialAuth: %v", err)
	}
}

func TestPasskeyCredentialFromWebAuthnNil(t *testing.T) {
	if got := passkeyCredentialFromWebAuthn(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestDomainPasskeyCredentialToWebAuthnNil(t *testing.T) {
	if got := passkeyCredentialToWebAuthn(nil); got.ID != nil {
		t.Fatal("expected empty credential")
	}
}

func TestDomainPasskeyCredentialToWebAuthnWithTransport(t *testing.T) {
	p := &domain.PasskeyCredential{
		CredentialID: []byte("id"),
		Transport:    "usb",
	}
	got := passkeyCredentialToWebAuthn(p)
	if len(got.Transport) != 1 {
		t.Fatalf("expected 1 transport, got %d", len(got.Transport))
	}
}
