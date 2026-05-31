package service

import (
	"github.com/blacksheepaul/timelog/pkg/errs"
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

func (s *Service) LoadPasskeyUser() (*PasskeyUser, error) {
	records, err := s.ListPasskeyCredentials()
	if err != nil {
		return nil, err
	}

	credentials := make([]webauthn.Credential, 0, len(records))
	for _, record := range records {
		credentials = append(credentials, record.ToCredential())
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

// --- Package-level wrappers (transitional) ---

func LoadPasskeyUser() (*PasskeyUser, error) {
	return defaultService.LoadPasskeyUser()
}

func LoadPasskeyUserByHandle(_ []byte, userHandle []byte) (webauthn.User, error) {
	return defaultService.LoadPasskeyUserByHandle(nil, userHandle)
}
