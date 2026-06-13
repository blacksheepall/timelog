package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

// passkeyCredentialRepo implements ports.PasskeyCredentialRepository using the model layer.
type passkeyCredentialRepo struct {
	dao *model.Dao
}

var _ ports.PasskeyCredentialRepository = (*passkeyCredentialRepo)(nil)

func newPasskeyCredentialRepo(dao *model.Dao) *passkeyCredentialRepo {
	return &passkeyCredentialRepo{dao: dao}
}

func (r *passkeyCredentialRepo) CreateWebAuthnCredential(credential *model.WebAuthnCredential) error {
	return model.CreateWebAuthnCredential(r.dao.Db(), credential)
}

func (r *passkeyCredentialRepo) GetWebAuthnCredentialByCredentialID(credentialID []byte) (*model.WebAuthnCredential, error) {
	return model.GetWebAuthnCredentialByCredentialID(r.dao.Db(), credentialID)
}

func (r *passkeyCredentialRepo) ListWebAuthnCredentials() ([]model.WebAuthnCredential, error) {
	return model.ListWebAuthnCredentials(r.dao.Db())
}

func (r *passkeyCredentialRepo) DeleteWebAuthnCredential(id uint) error {
	return model.DeleteWebAuthnCredential(r.dao.Db(), id)
}

func (r *passkeyCredentialRepo) UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error {
	return model.UpdateWebAuthnCredentialAuth(r.dao.Db(), credentialID, credential)
}
