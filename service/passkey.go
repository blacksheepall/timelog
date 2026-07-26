package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/errs"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/temppassword"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// InitWebAuthnWithConfig creates a WebAuthn instance from the given config.
// Callers must inject the returned instance into the Service via SetWebAuthn.
func InitWebAuthnWithConfig(cfg *config.Config) (*webauthn.WebAuthn, error) {
	if cfg == nil {
		return nil, errs.ErrPasskeyConfigNotInitialized
	}

	if cfg.Passkey.RPID == "" || cfg.Passkey.RPName == "" || len(cfg.Passkey.RPOrigins) == 0 {
		return nil, errs.ErrPasskeyConfigIncomplete
	}

	return webauthn.New(&webauthn.Config{
		RPID:          cfg.Passkey.RPID,
		RPDisplayName: cfg.Passkey.RPName,
		RPOrigins:     cfg.Passkey.RPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationPreferred,
		},
	})
}

// SetWebAuthn injects the WebAuthn instance used by passkey operations.
func (s *Service) SetWebAuthn(instance *webauthn.WebAuthn) {
	s.webAuthn = instance
}

func (s *Service) GetWebAuthn() *webauthn.WebAuthn {
	return s.webAuthn
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

func passkeyCredentialFromWebAuthn(credential *webauthn.Credential) *domain.PasskeyCredential {
	if credential == nil {
		return nil
	}
	transport := ""
	if len(credential.Transport) > 0 {
		transport = string(credential.Transport[0])
	}
	return &domain.PasskeyCredential{
		CredentialID:                  credential.ID,
		PublicKey:                     credential.PublicKey,
		AttestationType:               credential.AttestationType,
		Transport:                     transport,
		UserPresent:                   credential.Flags.UserPresent,
		UserVerified:                  credential.Flags.UserVerified,
		BackupEligible:                credential.Flags.BackupEligible,
		BackupState:                   credential.Flags.BackupState,
		AuthenticatorAaguid:           credential.Authenticator.AAGUID,
		AuthenticatorSignCount:        int32(credential.Authenticator.SignCount),
		AuthenticatorCloneWarning:     credential.Authenticator.CloneWarning,
		AuthenticatorAttachment:       string(credential.Authenticator.Attachment),
		AttestationClientDataJSON:     credential.Attestation.ClientDataJSON,
		AttestationClientDataHash:     credential.Attestation.ClientDataHash,
		AttestationAuthenticatorData:  credential.Attestation.AuthenticatorData,
		AttestationPublicKeyAlgorithm: int32(credential.Attestation.PublicKeyAlgorithm),
		AttestationObject:             credential.Attestation.Object,
	}
}

func (s *Service) CreatePasskeyCredential(credential *webauthn.Credential, deviceName string) (*domain.PasskeyCredential, error) {
	record := passkeyCredentialFromWebAuthn(credential)
	if record == nil {
		return nil, errs.ErrPasskeyCredentialNil
	}
	record.DeviceName = deviceName
	if err := s.passkeyRepo.CreateWebAuthnCredential(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) ListPasskeyCredentials() ([]domain.PasskeyCredential, error) {
	return s.passkeyRepo.ListWebAuthnCredentials()
}

func (s *Service) DeletePasskeyCredential(id int32) error {
	return s.passkeyRepo.DeleteWebAuthnCredential(id)
}

func (s *Service) LoadPasskeyCredentialByID(rawID []byte) (*domain.PasskeyCredential, error) {
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

func (s *Service) CreateTempPassword(ttlSeconds int) (*domain.TempPassword, string, error) {
	password, hash, err := s.GenerateTempPassword()
	if err != nil {
		return nil, "", err
	}

	record := &domain.TempPassword{
		PasswordHash: hash,
		ExpiresAt:    time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	if err := s.tempPasswordRepo.CreateTempPassword(record); err != nil {
		return nil, "", err
	}

	return record, password, nil
}

func (s *Service) ListTempPasswords() ([]domain.TempPassword, error) {
	return s.tempPasswordRepo.ListTempPasswords()
}

func (s *Service) DeleteTempPassword(id int32) error {
	return s.tempPasswordRepo.DeleteTempPassword(id)
}

func (s *Service) CleanupExpiredTempPasswords() error {
	return s.tempPasswordRepo.DeleteExpiredTempPasswords(time.Now())
}

func (s *Service) ValidateTempPassword(password string) (*domain.TempPassword, error) {
	hash := temppassword.HashPassword(password)
	return s.tempPasswordRepo.GetTempPasswordByHash(hash, time.Now())
}
