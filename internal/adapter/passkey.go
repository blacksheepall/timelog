package adapter

import (
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

type passkeyCredentialRepo struct {
	db model.DBProvider
}

var _ ports.PasskeyCredentialRepository = (*passkeyCredentialRepo)(nil)

func newPasskeyCredentialRepo(db model.DBProvider) *passkeyCredentialRepo {
	return &passkeyCredentialRepo{db: db}
}

func (r *passkeyCredentialRepo) CreateWebAuthnCredential(credential *domain.PasskeyCredential) error {
	return model.CreateWebAuthnCredential(r.db.Db(), toModelPasskeyCredential(credential))
}

func (r *passkeyCredentialRepo) GetWebAuthnCredentialByCredentialID(credentialID []byte) (*domain.PasskeyCredential, error) {
	m, err := model.GetWebAuthnCredentialByCredentialID(r.db.Db(), credentialID)
	if err != nil {
		return nil, err
	}
	return toDomainPasskeyCredential(m), nil
}

func (r *passkeyCredentialRepo) ListWebAuthnCredentials() ([]domain.PasskeyCredential, error) {
	list, err := model.ListWebAuthnCredentials(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainPasskeyCredentials(list), nil
}

func (r *passkeyCredentialRepo) DeleteWebAuthnCredential(id int32) error {
	return model.DeleteWebAuthnCredential(r.db.Db(), uint(id))
}

func (r *passkeyCredentialRepo) UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error {
	return model.UpdateWebAuthnCredentialAuth(r.db.Db(), credentialID, credential)
}
