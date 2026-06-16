// Package servicetest provides test helpers that depend on the service package.
package servicetest

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
)

// NewTestService creates a migrated in-memory DAO and assembles a Service.
func NewTestService(t *testing.T) (*service.Service, *model.Dao) {
	t.Helper()
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	repos := adapter.NewRepositories(dao)
	cfg := testutil.NewTestConfig()
	svc := service.NewService(
		repos, repos, repos, repos, repos, repos, repos, repos, repos,
		testutil.FakeLogger{}, cfg, nil,
	)
	return svc, dao
}
