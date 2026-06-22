package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/app"
	"github.com/blacksheepaul/timelog/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TimelogMCPServer is the main server struct
type TimelogMCPServer struct {
	service *service.Service
	config  *config.Config
}

// Tool parameter structs
type DateRangeParams struct {
	StartDate  string `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format"`
	EndDate    string `json:"end_date,omitempty" jsonschema:"End date in YYYY-MM-DD format"`
	ActiveOnly bool   `json:"active_only" jsonschema:"If true, return only ongoing time logs and ignore date range"`
}

type TaskStatusParams struct {
	Status string `json:"status" jsonschema:"Task status filter (completed/pending/all),required"`
}

type ConstraintParams struct{}

type CategoryListParams struct{}

type CreateTimeLogParams struct {
	CategoryID int32  `json:"category_id" jsonschema:"Category ID,required"`
	StartTime  string `json:"start_time,omitempty" jsonschema:"Start time in YYYY-MM-DD HH:MM:SS format (SGT), defaults to now"`
	EndTime    string `json:"end_time,omitempty" jsonschema:"End time in YYYY-MM-DD HH:MM:SS format (SGT), optional"`
	TaskID     int32  `json:"task_id,omitempty" jsonschema:"Associated task ID, optional"`
	Remark     string `json:"remark,omitempty" jsonschema:"Optional remark or description"`
}

type UpdateTimeLogParams struct {
	ID         int32  `json:"id" jsonschema:"Time log ID,required"`
	CategoryID int32  `json:"category_id,omitempty" jsonschema:"New category ID, optional"`
	StartTime  string `json:"start_time,omitempty" jsonschema:"New start time in YYYY-MM-DD HH:MM:SS format (SGT), optional"`
	EndTime    string `json:"end_time,omitempty" jsonschema:"New end time in YYYY-MM-DD HH:MM:SS format (SGT), optional"`
	TaskID     int32  `json:"task_id,omitempty" jsonschema:"Associated task ID, optional"`
	Remark     string `json:"remark,omitempty" jsonschema:"New remark, optional"`
}

type RecordMetricParams struct {
	MetricName string  `json:"metric_name" jsonschema:"Metric name,required"`
	Value      float64 `json:"value" jsonschema:"Metric value,required"`
	Source     string  `json:"source" jsonschema:"Source of the measurement,required"`
	RecordedAt string  `json:"recorded_at,omitempty" jsonschema:"Recording time in RFC3339 format"`
}

type IncrementMetricParams struct {
	MetricName string  `json:"metric_name" jsonschema:"Metric name,required"`
	Delta      float64 `json:"delta" jsonschema:"Amount to increment,required"`
	Source     string  `json:"source" jsonschema:"Source of the measurement,required"`
	RecordedAt string  `json:"recorded_at,omitempty" jsonschema:"Recording time in RFC3339 format"`
}

type GetMetricParams struct {
	Name string `json:"name" jsonschema:"Metric name,required"`
}

type ListMetricsParams struct{}

type EvaluateConstraintParams struct {
	ConstraintID int32 `json:"constraint_id" jsonschema:"Constraint ID,required"`
}

type CreateTaskParams struct {
	Title            string `json:"title" jsonschema:"Task title,required"`
	Description      string `json:"description,omitempty" jsonschema:"Task description,optional"`
	CategoryID       int32  `json:"category_id" jsonschema:"Associated category ID,required"`
	DueDate          string `json:"due_date,omitempty" jsonschema:"Due date in YYYY-MM-DD format,optional"`
	EstimatedMinutes int32  `json:"estimated_minutes,omitempty" jsonschema:"Estimated duration in minutes,optional"`
}

type UpdateTaskParams struct {
	ID               int32   `json:"id" jsonschema:"Task ID,required"`
	Title            *string `json:"title,omitempty" jsonschema:"New title,optional"`
	Description      *string `json:"description,omitempty" jsonschema:"New description,optional"`
	CategoryID       *int32  `json:"category_id,omitempty" jsonschema:"New category ID,optional"`
	DueDate          *string `json:"due_date,omitempty" jsonschema:"New due date in YYYY-MM-DD format,optional"`
	EstimatedMinutes *int32  `json:"estimated_minutes,omitempty" jsonschema:"New estimated duration in minutes,optional"`
}

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

	application, err := app.New(cfg, nil, nil)
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}

	LogMCPDebug("Database initialized", map[string]interface{}{
		"database_path": cfg.Database.Host,
	})

	return &TimelogMCPServer{
		service: application.Service,
		config:  cfg,
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

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all available categories for time logs",
	}, ListCategories)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_timelog",
		Description: "Create a new time log. If there is already an ongoing time log (without end_time), creation will be rejected",
	}, CreateTimeLog)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_timelog",
		Description: "Update an existing time log, such as setting end_time to stop the current activity",
	}, UpdateTimeLog)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "record_metric",
		Description: "Record a new value for a metric (write-only). The metric must already exist.",
	}, RecordMetric)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "increment_metric",
		Description: "Increment a counter metric by a delta (write-only). The metric must already exist.",
	}, IncrementMetric)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_metric",
		Description: "Get a metric by name",
	}, GetMetric)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_metrics",
		Description: "List all metrics",
	}, ListMetrics)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "evaluate_constraint",
		Description: "Evaluate whether a constraint's metric rule is currently met",
	}, EvaluateConstraintMCP)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_task",
		Description: "Create a new task",
	}, CreateTask)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_task",
		Description: "Update an existing task's metadata (title, description, category, due date, estimate). Does not change completion status.",
	}, UpdateTask)

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
