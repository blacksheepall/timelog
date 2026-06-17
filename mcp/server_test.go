package main

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

func TestNewTimelogMCPServer(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	cfg.MCP.Enabled = false
	cfg.Server.Port = 8080

	config.ResetConfig()
	config.SetConfig(cfg)
	defer config.ResetConfig()

	srv := NewTimelogMCPServer()
	if srv == nil {
		t.Fatal("expected server")
	}
	if srv.service == nil {
		t.Fatal("expected service")
	}
}
