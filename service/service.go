package service

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Service owns runtime dependencies for the business layer.
type Service struct {
	dao      *model.Dao
	log      logger.Logger
	webAuthn *webauthn.WebAuthn
}

var defaultService *Service
var daoProvider func() *model.Dao = model.GetDao
var webAuthnProvider func() *webauthn.WebAuthn

func New(dao *model.Dao, loggerInstance logger.Logger) *Service {
	return &Service{dao: dao, log: loggerInstance}
}

// InitService wires the process-wide service instance used by package-level helpers.
func InitService(loggerInstance logger.Logger, _ *config.Config, dao *model.Dao) *Service {
	defaultService = New(dao, loggerInstance)
	daoProvider = func() *model.Dao { return dao }
	return defaultService
}

func getDao() *model.Dao {
	if defaultService != nil {
		return defaultService.dao
	}
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
	if defaultService != nil && defaultService.webAuthn != nil {
		return defaultService.webAuthn
	}
	if webAuthnProvider == nil {
		return nil
	}
	return webAuthnProvider()
}

func setWebAuthnProviderForTest(provider func() *webauthn.WebAuthn) {
	webAuthnProvider = provider
}

func setWebAuthn(instance *webauthn.WebAuthn) {
	if defaultService != nil {
		defaultService.webAuthn = instance
	}
	webAuthnProvider = func() *webauthn.WebAuthn { return instance }
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
