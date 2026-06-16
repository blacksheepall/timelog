package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/testutil"
)

func TestNewServiceWiresInjectedRepositories(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	repos := adapter.NewRepositories(dao)
	svc := NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, testutil.NewTestConfig(), nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
