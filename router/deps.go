package router

import "github.com/blacksheepaul/timelog/service"

// Dependencies groups HTTP-layer runtime dependencies injected at bootstrap.
type Dependencies struct {
	Service *service.Service
}
