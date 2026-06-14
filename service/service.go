package service

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Service owns runtime dependencies for the business layer.
type Service struct {
	timelogRepo      ports.TimelogRepository
	categoryRepo     ports.CategoryRepository
	taskRepo         ports.TaskRepository
	constraintRepo   ports.ConstraintRepository
	metricRepo       ports.MetricRepository
	passkeyRepo      ports.PasskeyCredentialRepository
	tempPasswordRepo ports.TempPasswordRepository
	cache            ports.CacheStore
	unitOfWork       ports.UnitOfWork
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
	metricRepo ports.MetricRepository,
	passkeyRepo ports.PasskeyCredentialRepository,
	tempPasswordRepo ports.TempPasswordRepository,
	cache ports.CacheStore,
	unitOfWork ports.UnitOfWork,
	log logger.Logger,
	cfg *config.Config,
	webAuthn *webauthn.WebAuthn,
) *Service {
	return &Service{
		timelogRepo:      timelogRepo,
		categoryRepo:     categoryRepo,
		taskRepo:         taskRepo,
		constraintRepo:   constraintRepo,
		metricRepo:       metricRepo,
		passkeyRepo:      passkeyRepo,
		tempPasswordRepo: tempPasswordRepo,
		cache:            cache,
		unitOfWork:       unitOfWork,
		log:              log,
		webAuthn:         webAuthn,
		cfg:              cfg,
	}
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
