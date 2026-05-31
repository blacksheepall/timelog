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

// New bootstraps the database, repositories, and service layer.
// It returns an App whose Service is ready for handlers or MCP tools.
// Callers that still rely on package-level service helpers should also call
// service.SetDefaultService(app.Service) after this function returns.
func New(cfg *config.Config, log logger.Logger, webAuthn *webauthn.WebAuthn) (*App, error) {
	dao, err := model.NewDao(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	repos := adapter.NewRepositories(dao)
	svc := service.NewService(
		repos, // timelogRepo
		repos, // categoryRepo
		repos, // taskRepo
		repos, // constraintRepo
		repos, // passkeyRepo
		repos, // tempPasswordRepo
		repos, // cacheStore
		dao,   // dbProvider (transaction support)
		log,
		cfg,
		webAuthn,
	)

	return &App{DAO: dao, Service: svc}, nil
}
