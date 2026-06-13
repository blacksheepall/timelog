package model

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

type noopLogger struct{}

func (noopLogger) Debug(...interface{})          {}
func (noopLogger) Debugw(string, ...interface{}) {}
func (noopLogger) Info(...interface{})           {}
func (noopLogger) Infow(string, ...interface{})  {}
func (noopLogger) Warn(...interface{})           {}
func (noopLogger) Warnw(string, ...interface{})  {}
func (noopLogger) Error(...interface{})          {}
func (noopLogger) Errorw(string, ...interface{}) {}
func (noopLogger) Fatal(...interface{})          {}
func (noopLogger) Fatalw(string, ...interface{}) {}

func TestNewDaoOpensInMemoryDatabase(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1

	dao, err := NewDao(cfg, noopLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	if dao.Db() == nil {
		t.Fatal("expected non-nil gorm DB")
	}
	if dao.RawDB == nil {
		t.Fatal("expected non-nil raw DB")
	}
}
