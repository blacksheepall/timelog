package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/blacksheepaul/timelog/pkg/temppassword"
)

func InitWebAuthnWithConfig(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("config not initialized")
	}

	if cfg.Passkey.RPID == "" || cfg.Passkey.RPName == "" || len(cfg.Passkey.RPOrigins) == 0 {
		return errors.New("passkey config missing rp_id/rp_name/rp_origins")
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
	return nil
}

func GetWebAuthn() *webauthn.WebAuthn {
	return getWebAuthn()
}

func StorePasskeySession(sessionID string, session *webauthn.SessionData, ttlSeconds int64) error {
	if session == nil {
		return errors.New("session is nil")
	}
	dao := getDao()
	// Namespace the key to prevent confusion with auth tokens
	dao.WriteCache("passkey_session:"+sessionID, session, ttlSeconds)
	return nil
}

func LoadPasskeySession(sessionID string) (*webauthn.SessionData, error) {
	dao := getDao()
	// Use namespaced key to retrieve passkey session
	raw, ok := dao.GetCache("passkey_session:" + sessionID)
	if !ok {
		return nil, errors.New("session not found")
	}

	session, ok := raw.(*webauthn.SessionData)
	if !ok || session == nil {
		return nil, errors.New("invalid session data")
	}

	return session, nil
}

func CreatePasskeyCredential(credential *webauthn.Credential, deviceName string) (*model.WebAuthnCredential, error) {
	dao := getDao()
	record := model.WebAuthnCredentialFromCredential(credential)
	if record == nil {
		return nil, errors.New("credential is nil")
	}
	record.DeviceName = deviceName
	if err := model.CreateWebAuthnCredential(dao.Db(), record); err != nil {
		return nil, err
	}
	return record, nil
}

func ListPasskeyCredentials() ([]model.WebAuthnCredential, error) {
	dao := getDao()
	return model.ListWebAuthnCredentials(dao.Db())
}

func DeletePasskeyCredential(id uint) error {
	dao := getDao()
	return model.DeleteWebAuthnCredential(dao.Db(), id)
}

func LoadPasskeyCredentialByID(rawID []byte) (*model.WebAuthnCredential, error) {
	dao := getDao()
	return model.GetWebAuthnCredentialByCredentialID(dao.Db(), rawID)
}

func GenerateSessionToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func StoreSessionToken(token string, ttlSeconds int64) error {
	dao := getDao()
	// Namespace the key to distinguish from passkey sessions
	dao.WriteCache("auth_token:"+token, true, ttlSeconds)
	return nil
}

func GenerateTempPassword() (string, string, error) {
	return temppassword.GeneratePassword()
}

func CreateTempPassword(ttlSeconds int) (*model.TempPassword, string, error) {
	password, hash, err := GenerateTempPassword()
	if err != nil {
		return nil, "", err
	}

	dao := getDao()
	record := &model.TempPassword{
		PasswordHash: hash,
		ExpiresAt:    time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	if err := model.CreateTempPassword(dao.Db(), record); err != nil {
		return nil, "", err
	}

	return record, password, nil
}

func ListTempPasswords() ([]model.TempPassword, error) {
	dao := getDao()
	return model.ListTempPasswords(dao.Db())
}

func DeleteTempPassword(id uint) error {
	dao := getDao()
	return model.DeleteTempPassword(dao.Db(), id)
}

func CleanupExpiredTempPasswords() error {
	dao := getDao()
	return model.DeleteExpiredTempPasswords(dao.Db(), time.Now())
}

func ValidateTempPassword(password string) (*model.TempPassword, error) {
	hash := temppassword.HashPassword(password)

	dao := getDao()
	return model.GetTempPasswordByHash(dao.Db(), hash, time.Now())
}
