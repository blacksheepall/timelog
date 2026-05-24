package main

import (
	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
)

// NewTimelogMCPServer creates and initializes a new TimelogMCPServer instance
func NewTimelogMCPServer() *TimelogMCPServer {
	// Initialize configuration
	configPath := config.ResolveConfigPath("config.yml")
	cfg := config.GetConfig(configPath)

	// Initialize MCP logger (file-only, configurable)
	InitMCPLogger(cfg)
	LogMCPDebug("MCP server initializing", map[string]interface{}{
		"config_path":         configPath,
		"mcp_logging_enabled": cfg.MCP.Enabled,
	})

	// Disable ORM logging to prevent ANSI escape codes in MCP output
	cfg.Log.ORMLogLevel = 1 // Silent mode

	// Initialize database connection using existing DAO pattern
	// Pass nil for logger to avoid stdout output that interferes with MCP protocol
	model.InitDao(cfg, nil)
	dao := model.GetDao()

	LogMCPDebug("Database initialized", map[string]interface{}{
		"database_path": cfg.Database.Host,
	})

	return &TimelogMCPServer{
		db:     dao.Db(),
		config: cfg,
	}
}
