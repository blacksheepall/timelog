package service

import (
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/core/errs"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type PasskeyUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u *PasskeyUser) WebAuthnID() []byte {
	return u.id
}

func (u *PasskeyUser) WebAuthnName() string {
	return u.name
}

func (u *PasskeyUser) WebAuthnDisplayName() string {
	return u.displayName
}

func (u *PasskeyUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func passkeyCredentialToWebAuthn(p *domain.PasskeyCredential) webauthn.Credential {
	if p == nil {
		return webauthn.Credential{}
	}
	transport := []protocol.AuthenticatorTransport{}
	if p.Transport != "" {
		transport = append(transport, protocol.AuthenticatorTransport(p.Transport))
	}
	return webauthn.Credential{
		ID:              p.CredentialID,
		PublicKey:       p.PublicKey,
		AttestationType: p.AttestationType,
		Transport:       transport,
		Flags: webauthn.CredentialFlags{
			UserPresent:    p.UserPresent,
			UserVerified:   p.UserVerified,
			BackupEligible: p.BackupEligible,
			BackupState:    p.BackupState,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:       p.AuthenticatorAaguid,
			SignCount:    uint32(p.AuthenticatorSignCount),
			CloneWarning: p.AuthenticatorCloneWarning,
			Attachment:   protocol.AuthenticatorAttachment(p.AuthenticatorAttachment),
		},
		Attestation: webauthn.CredentialAttestation{
			ClientDataJSON:     p.AttestationClientDataJSON,
			ClientDataHash:     p.AttestationClientDataHash,
			AuthenticatorData:  p.AttestationAuthenticatorData,
			PublicKeyAlgorithm: int64(p.AttestationPublicKeyAlgorithm),
			Object:             p.AttestationObject,
		},
	}
}

func (s *Service) LoadPasskeyUser() (*PasskeyUser, error) {
	records, err := s.ListPasskeyCredentials()
	if err != nil {
		return nil, err
	}

	credentials := make([]webauthn.Credential, 0, len(records))
	for i := range records {
		credentials = append(credentials, passkeyCredentialToWebAuthn(&records[i]))
	}

	return &PasskeyUser{
		id:          []byte("timelog-single-user"),
		name:        "timelog",
		displayName: "TimeLog",
		credentials: credentials,
	}, nil
}

func (s *Service) LoadPasskeyUserByHandle(_ []byte, userHandle []byte) (webauthn.User, error) {
	if string(userHandle) != "timelog-single-user" {
		return nil, errs.ErrPasskeyUserNotFound
	}

	return s.LoadPasskeyUser()
}
