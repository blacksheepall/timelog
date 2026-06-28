package temppassword

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GeneratePassword generates a random hex-encoded password and its SHA256 hash.
// No external dependencies - pure password generation logic.
func GeneratePassword() (password string, hash string, err error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	password = hex.EncodeToString(raw)
	hashBytes := sha256.Sum256([]byte(password))
	hash = hex.EncodeToString(hashBytes[:])

	return password, hash, nil
}

// HashPassword computes the SHA256 hash of a password string.
func HashPassword(password string) string {
	hashBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashBytes[:])
}

// PasswordWithExpiry contains generated password and its expiry time.
type PasswordWithExpiry struct {
	Password  string    `json:"password"`
	Hash      string    `json:"hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GeneratePasswordWithExpiry generates a password with expiry time in seconds from now.
func GeneratePasswordWithExpiry(ttlSeconds int) (*PasswordWithExpiry, error) {
	password, hash, err := GeneratePassword()
	if err != nil {
		return nil, err
	}

	return &PasswordWithExpiry{
		Password:  password,
		Hash:      hash,
		ExpiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}, nil
}
