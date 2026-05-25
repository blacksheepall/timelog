package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

// TimelogMCPServer is the main server struct
type TimelogMCPServer struct {
	db     *gorm.DB
	config *config.Config
}

// Tool parameter structs
type DateRangeParams struct {
	StartDate  string `json:"start_date" jsonschema:"Start date in YYYY-MM-DD format"`
	EndDate    string `json:"end_date" jsonschema:"End date in YYYY-MM-DD format"`
	ActiveOnly bool   `json:"active_only" jsonschema:"If true, return only ongoing time logs and ignore date range"`
}

type TaskStatusParams struct {
	Status string `json:"status" jsonschema:"Task status filter (completed/pending/all),required"`
}

type ConstraintParams struct{}

var server *TimelogMCPServer

// NewTimelogMCPServer creates and initializes a new TimelogMCPServer instance
func NewTimelogMCPServer() *TimelogMCPServer {
	configPath := config.ResolveConfigPath("config.yml")
	cfg := config.GetConfig(configPath)

	InitMCPLogger(cfg)
	LogMCPDebug("MCP server initializing", map[string]interface{}{
		"config_path":         configPath,
		"mcp_logging_enabled": cfg.MCP.Enabled,
	})

	cfg.Log.ORMLogLevel = 1

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

func main() {
	server = NewTimelogMCPServer()

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "timelog",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_timelogs_by_date_range",
		Description: "Get time logs within a specific date range",
	}, GetTimeLogsByDateRange)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_tasks_by_status",
		Description: "Get tasks filtered by completion status",
	}, GetTasksByStatus)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_active_constraints",
		Description: "To know self discipline and external conditions",
	}, GetActiveConstraints)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_date_info",
		Description: "Get current date, time, today, yesterday, and this week's date range",
	}, GetDateInfo)

	transportMode := server.config.MCP.Transport
	switch transportMode {
	case "http":
		listenAddr := server.config.MCP.ListenAddr
		token := server.config.MCP.Token

		if token == "" {
			fmt.Fprintln(os.Stderr, "FATAL: HTTP transport requires authentication token for security.")
			fmt.Fprintln(os.Stderr, "Set token via: MCP.Token in config.yml OR MCP_TOKEN environment variable")
			os.Exit(1)
		}

		handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
			return mcpServer
		}, nil)

		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
				return
			}

			if token != "" {
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				provided := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			handler.ServeHTTP(w, r)
		})

		httpServer := &http.Server{
			Addr:    listenAddr,
			Handler: wrappedHandler,
		}

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			LogMCPError("http_server", err, map[string]interface{}{"addr": listenAddr})
		}
	default:
		ctx := context.Background()
		transport := &mcp.StdioTransport{}
		if err := mcpServer.Run(ctx, transport); err != nil {
			LogMCPError("stdio_server", err, nil)
		}
	}
}
