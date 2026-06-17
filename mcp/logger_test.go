package main

import (
	"errors"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

func TestInitMCPLoggerDisabled(t *testing.T) {
	old := mcpLogger
	t.Cleanup(func() { mcpLogger = old })

	cfg := &config.Config{}
	cfg.MCP.Enabled = false
	log := InitMCPLogger(cfg)
	if log == nil {
		t.Fatal("expected nop logger")
	}
}

func TestInitMCPLoggerDebugFalse(t *testing.T) {
	old := mcpLogger
	t.Cleanup(func() { mcpLogger = old })

	cfg := &config.Config{}
	cfg.MCP.Enabled = true
	t.Setenv("MCP_DEBUG", "false")
	log := InitMCPLogger(cfg)
	if log == nil {
		t.Fatal("expected nop logger")
	}
}

func TestLogMCPErrorBeforeInit(t *testing.T) {
	old := mcpLogger
	t.Cleanup(func() { mcpLogger = old })

	mcpLogger = nil
	LogMCPError("op", errors.New("boom"), map[string]interface{}{"k": "v"})
}

func TestLogMCPDebugBeforeInit(t *testing.T) {
	old := mcpLogger
	t.Cleanup(func() { mcpLogger = old })

	mcpLogger = nil
	LogMCPDebug("msg", map[string]interface{}{"k": "v"})
}
