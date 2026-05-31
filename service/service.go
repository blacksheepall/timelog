package service

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Service owns runtime dependencies for the business layer.
type Service struct {
	timelogRepo      ports.TimelogRepository
	categoryRepo     ports.CategoryRepository
	taskRepo         ports.TaskRepository
	constraintRepo   ports.ConstraintRepository
	passkeyRepo      ports.PasskeyCredentialRepository
	tempPasswordRepo ports.TempPasswordRepository
	cache            ports.CacheStore
	dbProvider       model.DBProvider
	log              logger.Logger
	webAuthn         *webauthn.WebAuthn
	cfg              *config.Config
}

// NewService creates a Service with injected repository ports.
func NewService(
	timelogRepo ports.TimelogRepository,
	categoryRepo ports.CategoryRepository,
	taskRepo ports.TaskRepository,
	constraintRepo ports.ConstraintRepository,
	passkeyRepo ports.PasskeyCredentialRepository,
	tempPasswordRepo ports.TempPasswordRepository,
	cache ports.CacheStore,
	dbProvider model.DBProvider,
	log logger.Logger,
	cfg *config.Config,
	webAuthn *webauthn.WebAuthn,
) *Service {
	return &Service{
		timelogRepo:      timelogRepo,
		categoryRepo:     categoryRepo,
		taskRepo:         taskRepo,
		constraintRepo:   constraintRepo,
		passkeyRepo:      passkeyRepo,
		tempPasswordRepo: tempPasswordRepo,
		cache:            cache,
		dbProvider:       dbProvider,
		log:              log,
		cfg:              cfg,
		webAuthn:         webAuthn,
	}
}

var defaultService *Service

func getDefaultService() *Service {
	return defaultService
}

// SetDefaultService wires the process-wide defaultService used by package-level helpers.
func SetDefaultService(svc *Service) {
	defaultService = svc
}

var webAuthnProvider func() *webauthn.WebAuthn

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

// GetCache implements ports.SessionTokenStore so Service can be passed to auth middleware.
func (s *Service) GetCache(key string) (any, bool) {
	return s.cache.GetCache(key)
}

// WriteCache exposes cache write for passkey/session management.
func (s *Service) WriteCache(key string, value any, seconds int64) {
	s.cache.WriteCache(key, value, seconds)
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
