package service

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

func (s *Service) UpdatePasskeyCredentialAuth(credential *webauthn.Credential) error {
	if credential == nil {
		return nil
	}

	return s.passkeyRepo.UpdateWebAuthnCredentialAuth(credential.ID, credential)
}

// --- Package-level wrapper (transitional) ---

func UpdatePasskeyCredentialAuth(credential *webauthn.Credential) error {
	return defaultService.UpdatePasskeyCredentialAuth(credential)
}
