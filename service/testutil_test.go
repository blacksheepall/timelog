package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/domain"
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
	repos := adapter.NewRepositories(dao, testutil.FakeLogger{})
	cfg := testutil.NewTestConfig()
	svc := NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, cfg, nil)
	return svc, dao
}

func strPtr(s string) *string       { return &s }
func int32Ptr(v int32) *int32       { return &v }
func boolPtr(v bool) *bool          { return &v }
func float64Ptr(v float64) *float64 { return &v }

func seedServiceConstraint(t *testing.T, svc *Service) *domain.Constraint {
	t.Helper()
	c := &domain.Constraint{
		Description:     "seed",
		PunishmentQuote: "seed",
		StartDate:       time.Now().UTC(),
		IsActive:        true,
	}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("seed constraint: %v", err)
	}
	return c
}
