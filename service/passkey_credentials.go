package service

import (
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/blacksheepaul/timelog/model"
)

func UpdatePasskeyCredentialAuth(credential *webauthn.Credential) error {
	if credential == nil {
		return nil
	}

	dao := getDao()
	return model.UpdateWebAuthnCredentialAuth(dao.Db(), credential.ID, credential)
}
