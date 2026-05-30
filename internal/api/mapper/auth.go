package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/model"
)

func PasskeyCredentialToProto(credential *model.WebAuthnCredential) *timelogv1.PasskeyCredential {
	if credential == nil {
		return nil
	}
	return &timelogv1.PasskeyCredential{
		Id:         uint32(credential.ID),
		DeviceName: credential.DeviceName,
		CreatedAt:  FormatTimeUTC(credential.CreatedAt),
	}
}

func PasskeyCredentialsToProto(credentials []model.WebAuthnCredential) []*timelogv1.PasskeyCredential {
	out := make([]*timelogv1.PasskeyCredential, 0, len(credentials))
	for i := range credentials {
		out = append(out, PasskeyCredentialToProto(&credentials[i]))
	}
	return out
}

func LoginResponse(token string, tokenType string, expiresIn int64) *timelogv1.PasskeyLoginResponse {
	return &timelogv1.PasskeyLoginResponse{
		Token:     token,
		TokenType: tokenType,
		ExpiresIn: expiresIn,
	}
}
