package temppassword

import (
	"testing"
	"time"
)

func TestGeneratePassword(t *testing.T) {
	password, hash, err := GeneratePassword()
	if err != nil {
		t.Fatalf("GeneratePassword() error = %v", err)
	}

	if len(password) != 32 {
		t.Fatalf("password length = %d, want 32", len(password))
	}
	if len(hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(hash))
	}
	if got := HashPassword(password); got != hash {
		t.Fatalf("HashPassword() = %s, want %s", got, hash)
	}
}

func TestGeneratePasswordWithExpiry(t *testing.T) {
	start := time.Now()
	result, err := GeneratePasswordWithExpiry(90)
	if err != nil {
		t.Fatalf("GeneratePasswordWithExpiry() error = %v", err)
	}

	if result == nil {
		t.Fatal("GeneratePasswordWithExpiry() returned nil result")
	}
	if result.ExpiresAt.Before(start.Add(89*time.Second)) || result.ExpiresAt.After(start.Add(91*time.Second)) {
		t.Fatalf("ExpiresAt = %s, want about %s", result.ExpiresAt, start.Add(90*time.Second))
	}
}
