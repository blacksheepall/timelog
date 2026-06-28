package app

import (
	"fmt"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/go-webauthn/webauthn/webauthn"
)

// App encapsulates the composed application layer for all entrypoints.
type App struct {
	DAO     *model.Dao
	Service *service.Service
}

// New bootstraps the database, per-domain repository adapters, and service layer.
// It returns an App whose Service is ready for handlers or MCP tools.
func New(cfg *config.Config, log logger.Logger, webAuthn *webauthn.WebAuthn) (*App, error) {
	dao, err := model.NewDao(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	repos := adapter.NewRepositories(dao, log)
	svc := service.NewService(
		repos, // timelogRepo
		repos, // categoryRepo
		repos, // taskRepo
		repos, // constraintRepo
		repos, // metricRepo
		repos, // passkeyRepo
		repos, // tempPasswordRepo
		repos, // cacheStore
		repos, // unitOfWork
		repos, // auditLogger
		log,
		cfg,
		webAuthn,
	)

	return &App{DAO: dao, Service: svc}, nil
}
