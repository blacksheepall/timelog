package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupMCPTestServer(t *testing.T) *service.Service {
	t.Helper()
	cfg := testutil.NewTestConfig()
	dao, err := model.NewDao(cfg, testutil.FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	testutil.ApplyMigrations(t, dao)
	repos := adapter.NewRepositories(dao)
	svc := service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, cfg, nil)

	oldServer := server
	t.Cleanup(func() { server = oldServer })
	server = &TimelogMCPServer{service: svc, config: cfg}
	return svc
}

func TestGetDateInfo(t *testing.T) {
	setupMCPTestServer(t)
	res, _, err := GetDateInfo(context.Background(), nil, DateInfoParams{})
	if err != nil {
		t.Fatalf("GetDateInfo: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestGetTimeLogsByDateRangeActiveOnly(t *testing.T) {
	svc := setupMCPTestServer(t)
	if _, err := svc.CreateTimeLogFromMCPInput(service.CreateTimeLogMCPInput{CategoryID: 1, Remark: "active"}); err != nil {
		t.Fatalf("create timelog: %v", err)
	}

	res, _, err := GetTimeLogsByDateRange(context.Background(), nil, DateRangeParams{ActiveOnly: true})
	if err != nil {
		t.Fatalf("GetTimeLogsByDateRange: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}

	text := res.Content[0].(*mcp.TextContent).Text
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	logs := payload["active_logs"].([]interface{})
	if len(logs) != 1 {
		t.Fatalf("expected 1 active log, got %d", len(logs))
	}
	entry := logs[0].(map[string]interface{})
	if entry["category_id"].(float64) != 1 {
		t.Fatalf("expected category_id 1, got %v", entry["category_id"])
	}
	if entry["category_name"].(string) != "工作" {
		t.Fatalf("expected category_name 工作, got %v", entry["category_name"])
	}
}

func TestGetTimeLogsByDateRange(t *testing.T) {
	svc := setupMCPTestServer(t)
	if _, err := svc.CreateTimeLogFromMCPInput(service.CreateTimeLogMCPInput{
		CategoryID: 1,
		StartTime:  service.FormatSGDateTime(time.Now().Add(-time.Hour)),
		EndTime:    service.FormatSGDateTime(time.Now()),
		Remark:     "done",
	}); err != nil {
		t.Fatalf("create timelog: %v", err)
	}

	res, _, err := GetTimeLogsByDateRange(context.Background(), nil, DateRangeParams{StartDate: service.TodaySGDateString(), EndDate: service.TodaySGDateString()})
	if err != nil {
		t.Fatalf("GetTimeLogsByDateRange: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}

	text := res.Content[0].(*mcp.TextContent).Text
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	logs := payload["time_logs"].([]interface{})
	if len(logs) != 1 {
		t.Fatalf("expected 1 time log, got %d", len(logs))
	}
	entry := logs[0].(map[string]interface{})
	if entry["category_id"].(float64) != 1 {
		t.Fatalf("expected category_id 1, got %v", entry["category_id"])
	}
	if entry["category_name"].(string) != "工作" {
		t.Fatalf("expected category_name 工作, got %v", entry["category_name"])
	}
}

func TestGetTasksByStatus(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateCategory(&domain.Category{Name: "tasks"}); err != nil {
		t.Fatalf("create category: %v", err)
	}
	if err := svc.CreateTask(&domain.Task{Title: "t1", CategoryID: 1, DueDate: time.Now()}); err != nil {
		t.Fatalf("create task: %v", err)
	}

	res, _, err := GetTasksByStatus(context.Background(), nil, TaskStatusParams{Status: "pending"})
	if err != nil {
		t.Fatalf("GetTasksByStatus: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestGetTasksByStatusInvalid(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := GetTasksByStatus(context.Background(), nil, TaskStatusParams{Status: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetActiveConstraints(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateConstraint(&domain.Constraint{Description: "focus", PunishmentQuote: "punish", StartDate: time.Now(), IsActive: true}); err != nil {
		t.Fatalf("create constraint: %v", err)
	}

	res, _, err := GetActiveConstraints(context.Background(), nil, ConstraintParams{})
	if err != nil {
		t.Fatalf("GetActiveConstraints: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestListCategoriesMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateCategory(&domain.Category{Name: "cat"}); err != nil {
		t.Fatalf("create category: %v", err)
	}

	res, _, err := ListCategories(context.Background(), nil, CategoryListParams{})
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestCreateTimeLogMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateCategory(&domain.Category{Name: "work"}); err != nil {
		t.Fatalf("create category: %v", err)
	}

	res, _, err := CreateTimeLog(context.Background(), nil, CreateTimeLogParams{CategoryID: 1, Remark: "hello"})
	if err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestUpdateTimeLogMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateCategory(&domain.Category{Name: "work"}); err != nil {
		t.Fatalf("create category: %v", err)
	}
	log, err := svc.CreateTimeLogFromMCPInput(service.CreateTimeLogMCPInput{CategoryID: 1, Remark: "orig"})
	if err != nil {
		t.Fatalf("create timelog: %v", err)
	}

	res, _, err := UpdateTimeLog(context.Background(), nil, UpdateTimeLogParams{ID: log.ID, Remark: "updated"})
	if err != nil {
		t.Fatalf("UpdateTimeLog: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestRecordMetricMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "pages", MetricType: "counter", Unit: "count", CurrentValue: float64Ptr(0)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}

	res, _, err := RecordMetric(context.Background(), nil, RecordMetricParams{MetricName: "pages", Value: 5, Source: "test"})
	if err != nil {
		t.Fatalf("RecordMetric: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestIncrementMetricMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "steps", MetricType: "counter", Unit: "count", CurrentValue: float64Ptr(0)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}

	res, _, err := IncrementMetric(context.Background(), nil, IncrementMetricParams{MetricName: "steps", Delta: 10, Source: "test"})
	if err != nil {
		t.Fatalf("IncrementMetric: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestGetMetricMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "weight", MetricType: "gauge", Unit: "kg", CurrentValue: float64Ptr(70)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}

	res, _, err := GetMetric(context.Background(), nil, GetMetricParams{Name: "weight"})
	if err != nil {
		t.Fatalf("GetMetric: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestListMetricsMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "m1", MetricType: "counter", Unit: "c", CurrentValue: float64Ptr(1)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}

	res, _, err := ListMetrics(context.Background(), nil, ListMetricsParams{})
	if err != nil {
		t.Fatalf("ListMetrics: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestEvaluateConstraintMCP(t *testing.T) {
	svc := setupMCPTestServer(t)
	if err := svc.CreateMetric(&domain.Metric{Name: "pages", MetricType: "counter", Unit: "count", CurrentValue: float64Ptr(10)}); err != nil {
		t.Fatalf("create metric: %v", err)
	}
	c := &domain.Constraint{Description: "focus", PunishmentQuote: "punish", StartDate: time.Now(), IsActive: true}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("create constraint: %v", err)
	}
	c.MetricID = int32Ptr(1)
	c.MetricOperator = strPtr("gte")
	c.MetricTargetValue = float64Ptr(5)
	if err := svc.UpdateConstraint(c); err != nil {
		t.Fatalf("update constraint: %v", err)
	}

	res, _, err := EvaluateConstraintMCP(context.Background(), nil, EvaluateConstraintParams{ConstraintID: 1})
	if err != nil {
		t.Fatalf("EvaluateConstraintMCP: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestCreateTimeLogMCPInvalidCategory(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := CreateTimeLog(context.Background(), nil, CreateTimeLogParams{CategoryID: 9999})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateTimeLogMCPNotFound(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := UpdateTimeLog(context.Background(), nil, UpdateTimeLogParams{ID: 9999})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRecordMetricMCPNotFound(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := RecordMetric(context.Background(), nil, RecordMetricParams{MetricName: "missing", Value: 1, Source: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetMetricMCPNotFound(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := GetMetric(context.Background(), nil, GetMetricParams{Name: "missing"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEvaluateConstraintMCPNotFound(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := EvaluateConstraintMCP(context.Background(), nil, EvaluateConstraintParams{ConstraintID: 9999})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateTaskMCP(t *testing.T) {
	setupMCPTestServer(t)

	res, _, err := CreateTask(context.Background(), nil, CreateTaskParams{
		Title:            "New task",
		Description:      "A test task",
		CategoryID:       1,
		DueDate:          service.TodaySGDateString(),
		EstimatedMinutes: 30,
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}

	text := res.Content[0].(*mcp.TextContent).Text
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if payload["title"].(string) != "New task" {
		t.Fatalf("expected title New task, got %v", payload["title"])
	}
	if payload["category_id"].(float64) != 1 {
		t.Fatalf("expected category_id 1, got %v", payload["category_id"])
	}
}

func TestCreateTaskMCPMissingTitle(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := CreateTask(context.Background(), nil, CreateTaskParams{CategoryID: 1})
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestCreateTaskMCPMissingCategory(t *testing.T) {
	setupMCPTestServer(t)
	_, _, err := CreateTask(context.Background(), nil, CreateTaskParams{Title: "No category"})
	if err == nil {
		t.Fatal("expected error for missing category")
	}
}
