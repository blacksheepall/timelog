package service

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

var log logger.Logger
var daoProvider func() *model.Dao = model.GetDao
var webAuthnProvider func() *webauthn.WebAuthn

func InitService(loggerInstance logger.Logger, _ *config.Config) {
	log = loggerInstance
}

func getDao() *model.Dao {
	return daoProvider()
}

func setDaoProviderForTest(provider func() *model.Dao) {
	if provider == nil {
		daoProvider = model.GetDao
		return
	}
	daoProvider = provider
}

func getWebAuthn() *webauthn.WebAuthn {
	if webAuthnProvider == nil {
		return nil
	}
	return webAuthnProvider()
}

func setWebAuthnProviderForTest(provider func() *webauthn.WebAuthn) {
	webAuthnProvider = provider
}

type Response struct {
	Items []any `json:"items"`
	Pages
}

type Pages struct {
	Page  int `form:"page" json:"page"`
	Size  int `form:"size" json:"size"`
	Total int `form:"total" json:"total"`
}
