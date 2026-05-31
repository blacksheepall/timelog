package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/pkg/errs"
	"github.com/blacksheepaul/timelog/pkg/temppassword"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func InitWebAuthnWithConfig(cfg *config.Config) error {
	if cfg == nil {
		return errs.ErrPasskeyConfigNotInitialized
	}

	if cfg.Passkey.RPID == "" || cfg.Passkey.RPName == "" || len(cfg.Passkey.RPOrigins) == 0 {
		return errs.ErrPasskeyConfigIncomplete
	}

	instance, err := webauthn.New(&webauthn.Config{
		RPID:          cfg.Passkey.RPID,
		RPDisplayName: cfg.Passkey.RPName,
		RPOrigins:     cfg.Passkey.RPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationPreferred,
		},
	})
	if err != nil {
		return err
	}

	webAuthnProvider = func() *webauthn.WebAuthn {
		return instance
	}
	setWebAuthn(instance)
	return nil
}

func (s *Service) GetWebAuthn() *webauthn.WebAuthn {
	return getWebAuthn()
}

func (s *Service) StorePasskeySession(sessionID string, session *webauthn.SessionData, ttlSeconds int64) error {
	if session == nil {
		return errs.ErrPasskeySessionNil
	}
	// Namespace the key to prevent confusion with auth tokens
	s.cache.WriteCache("passkey_session:"+sessionID, session, ttlSeconds)
	return nil
}

func (s *Service) LoadPasskeySession(sessionID string) (*webauthn.SessionData, error) {
	// Use namespaced key to retrieve passkey session
	raw, ok := s.cache.GetCache("passkey_session:" + sessionID)
	if !ok {
		return nil, errs.ErrPasskeySessionNotFound
	}

	session, ok := raw.(*webauthn.SessionData)
	if !ok || session == nil {
		return nil, errs.ErrPasskeySessionInvalid
	}

	return session, nil
}

func (s *Service) CreatePasskeyCredential(credential *webauthn.Credential, deviceName string) (*model.WebAuthnCredential, error) {
	record := model.WebAuthnCredentialFromCredential(credential)
	if record == nil {
		return nil, errs.ErrPasskeyCredentialNil
	}
	record.DeviceName = deviceName
	if err := s.passkeyRepo.CreateWebAuthnCredential(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) ListPasskeyCredentials() ([]model.WebAuthnCredential, error) {
	return s.passkeyRepo.ListWebAuthnCredentials()
}

func (s *Service) DeletePasskeyCredential(id uint) error {
	return s.passkeyRepo.DeleteWebAuthnCredential(id)
}

func (s *Service) LoadPasskeyCredentialByID(rawID []byte) (*model.WebAuthnCredential, error) {
	return s.passkeyRepo.GetWebAuthnCredentialByCredentialID(rawID)
}

func (s *Service) GenerateSessionToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func (s *Service) StoreSessionToken(token string, ttlSeconds int64) error {
	// Namespace the key to distinguish from passkey sessions
	s.cache.WriteCache("auth_token:"+token, true, ttlSeconds)
	return nil
}

func (s *Service) GenerateTempPassword() (string, string, error) {
	return temppassword.GeneratePassword()
}

func (s *Service) CreateTempPassword(ttlSeconds int) (*model.TempPassword, string, error) {
	password, hash, err := s.GenerateTempPassword()
	if err != nil {
		return nil, "", err
	}

	record := &model.TempPassword{
		PasswordHash: hash,
		ExpiresAt:    time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	if err := s.tempPasswordRepo.CreateTempPassword(record); err != nil {
		return nil, "", err
	}

	return record, password, nil
}

func (s *Service) ListTempPasswords() ([]model.TempPassword, error) {
	return s.tempPasswordRepo.ListTempPasswords()
}

func (s *Service) DeleteTempPassword(id uint) error {
	return s.tempPasswordRepo.DeleteTempPassword(id)
}

func (s *Service) CleanupExpiredTempPasswords() error {
	return s.tempPasswordRepo.DeleteExpiredTempPasswords(time.Now())
}

func (s *Service) ValidateTempPassword(password string) (*model.TempPassword, error) {
	hash := temppassword.HashPassword(password)
	return s.tempPasswordRepo.GetTempPasswordByHash(hash, time.Now())
}

// --- Package-level wrappers (transitional) ---

func StorePasskeySession(sessionID string, session *webauthn.SessionData, ttlSeconds int64) error {
	return defaultService.StorePasskeySession(sessionID, session, ttlSeconds)
}
func LoadPasskeySession(sessionID string) (*webauthn.SessionData, error) {
	return defaultService.LoadPasskeySession(sessionID)
}
func CreatePasskeyCredential(credential *webauthn.Credential, deviceName string) (*model.WebAuthnCredential, error) {
	return defaultService.CreatePasskeyCredential(credential, deviceName)
}
func ListPasskeyCredentials() ([]model.WebAuthnCredential, error) {
	return defaultService.ListPasskeyCredentials()
}
func DeletePasskeyCredential(id uint) error { return defaultService.DeletePasskeyCredential(id) }
func LoadPasskeyCredentialByID(rawID []byte) (*model.WebAuthnCredential, error) {
	return defaultService.LoadPasskeyCredentialByID(rawID)
}
func GenerateSessionToken() (string, error) { return defaultService.GenerateSessionToken() }
func StoreSessionToken(token string, ttlSeconds int64) error {
	return defaultService.StoreSessionToken(token, ttlSeconds)
}
func GenerateTempPassword() (string, string, error) { return defaultService.GenerateTempPassword() }
func CreateTempPassword(ttlSeconds int) (*model.TempPassword, string, error) {
	return defaultService.CreateTempPassword(ttlSeconds)
}
func ListTempPasswords() ([]model.TempPassword, error) { return defaultService.ListTempPasswords() }
func DeleteTempPassword(id uint) error                 { return defaultService.DeleteTempPassword(id) }
func CleanupExpiredTempPasswords() error               { return defaultService.CleanupExpiredTempPasswords() }
func ValidateTempPassword(password string) (*model.TempPassword, error) {
	return defaultService.ValidateTempPassword(password)
}
