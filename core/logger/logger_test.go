package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

func TestGetLoggerBeforeSet(t *testing.T) {
	ilogger = nil
	if got := GetLogger(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestSetZapLogger(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{}
	cfg.Log.Level = "info"
	cfg.Log.Path = filepath.Join(dir, "test.log")
	cfg.Log.Rotation.MaxSize = 1
	cfg.Log.Rotation.MaxBackups = 1
	cfg.Log.Rotation.MaxAge = 1

	log := SetZapLogger(cfg)
	if log == nil {
		t.Fatal("expected logger")
	}

	// Ensure log directory exists (file may be created lazily on first write)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("log directory missing: %v", err)
	}
}

func TestSetZapLoggerCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{}
	cfg.Log.Level = "debug"
	cfg.Log.Path = filepath.Join(dir, "sub", "dir", "test.log")
	cfg.Log.Rotation.MaxSize = 1
	cfg.Log.Rotation.MaxBackups = 1
	cfg.Log.Rotation.MaxAge = 1

	log := SetZapLogger(cfg)
	if log == nil {
		t.Fatal("expected logger")
	}
}
