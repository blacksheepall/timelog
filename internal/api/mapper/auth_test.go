package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestPasskeyCredentialToProto(t *testing.T) {
	credential := &domain.PasskeyCredential{
		ID:         42,
		DeviceName: "Laptop",
		CreatedAt:  time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
	}
	got := PasskeyCredentialToProto(credential)
	if got.GetId() != 42 || got.GetDeviceName() != "Laptop" || got.GetCreatedAt() != "2026-05-30T12:00:00Z" {
		t.Fatalf("unexpected mapped passkey credential: %#v", got)
	}
}

func TestPasskeyCredentialToProtoNil(t *testing.T) {
	if PasskeyCredentialToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestPasskeyCredentialsToProto(t *testing.T) {
	credentials := []domain.PasskeyCredential{
		{ID: 1, DeviceName: "A"},
		{ID: 2, DeviceName: "B"},
	}
	got := PasskeyCredentialsToProto(credentials)
	if len(got) != 2 || got[0].GetDeviceName() != "A" || got[1].GetDeviceName() != "B" {
		t.Fatalf("unexpected mapped credentials: %#v", got)
	}
}

func TestLoginResponse(t *testing.T) {
	got := LoginResponse("token", "Bearer", 3600)
	if got.GetToken() != "token" || got.GetTokenType() != "Bearer" || got.GetExpiresIn() != 3600 {
		t.Fatalf("unexpected login response: %#v", got)
	}
}
