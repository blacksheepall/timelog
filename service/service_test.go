package service

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
)

func TestInitServiceWiresInjectedDAO(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1

	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}

	svc := InitService(FakeLogger{}, cfg, dao)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}

	if getDao() != dao {
		t.Fatal("getDao should return the injected DAO instance")
	}
}
