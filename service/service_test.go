package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/model"
)

func TestNewServiceWiresInjectedRepositories(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1

	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}

	repos := adapter.NewRepositories(dao)
	svc := NewService(repos, repos, repos, repos, repos, repos, repos, repos, FakeLogger{}, cfg, nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
