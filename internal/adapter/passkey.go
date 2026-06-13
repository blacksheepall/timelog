package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

// passkeyCredentialRepo implements ports.PasskeyCredentialRepository using the model layer.
type passkeyCredentialRepo struct {
	db model.DBProvider
}

var _ ports.PasskeyCredentialRepository = (*passkeyCredentialRepo)(nil)

func newPasskeyCredentialRepo(db model.DBProvider) *passkeyCredentialRepo {
	return &passkeyCredentialRepo{db: db}
}

func (r *passkeyCredentialRepo) CreateWebAuthnCredential(credential *model.WebAuthnCredential) error {
	return model.CreateWebAuthnCredential(r.db.Db(), credential)
}

func (r *passkeyCredentialRepo) GetWebAuthnCredentialByCredentialID(credentialID []byte) (*model.WebAuthnCredential, error) {
	return model.GetWebAuthnCredentialByCredentialID(r.db.Db(), credentialID)
}

func (r *passkeyCredentialRepo) ListWebAuthnCredentials() ([]model.WebAuthnCredential, error) {
	return model.ListWebAuthnCredentials(r.db.Db())
}

func (r *passkeyCredentialRepo) DeleteWebAuthnCredential(id uint) error {
	return model.DeleteWebAuthnCredential(r.db.Db(), id)
}

func (r *passkeyCredentialRepo) UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error {
	return model.UpdateWebAuthnCredentialAuth(r.db.Db(), credentialID, credential)
}
