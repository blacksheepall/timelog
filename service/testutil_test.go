package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
)

// newTestService creates a migrated in-memory DAO and assembles a Service for
// use by tests in the service package. It lives here rather than in
// internal/testutil to avoid an import cycle (internal/testutil must not import
// service, and service tests must not import a package that imports service).
func newTestService(t *testing.T) (*Service, *model.Dao) {
	t.Helper()
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	repos := adapter.NewRepositories(dao)
	cfg := testutil.NewTestConfig()
	svc := NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, cfg, nil)
	return svc, dao
}

func strPtr(s string) *string    { return &s }
func int32Ptr(v int32) *int32    { return &v }
func boolPtr(v bool) *bool       { return &v }
func float64Ptr(v float64) *float64 { return &v }
