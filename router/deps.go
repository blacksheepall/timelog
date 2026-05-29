package router

import (
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
)

// Dependencies groups HTTP-layer runtime dependencies injected at bootstrap.
// Handlers still use package-level service helpers during the transitional refactor.
type Dependencies struct {
	Service *service.Service
	DAO     *model.Dao
}
