package router

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

type routerFakeDataSource struct {
	name   string
	points []domain.MetricDataPoint
}

func (f *routerFakeDataSource) Name() string { return f.name }
func (f *routerFakeDataSource) Fetch(ctx context.Context) ([]domain.MetricDataPoint, error) {
	return f.points, nil
}

type routerFakeRegistry struct{}

func (r *routerFakeRegistry) Get(name string) (ports.DataSource, error) {
	if name == "maimemo" {
		return &routerFakeDataSource{
			name: "maimemo",
			points: []domain.MetricDataPoint{
				{MetricName: "今日已背单词", Value: 120, RecordedAt: time.Now().UTC(), Source: "maimemo"},
			},
		}, nil
	}
	return nil, fmt.Errorf("datasource %q not found", name)
}

func setupTestRouter(t *testing.T) (*gin.Engine, *service.Service) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	cfg := testutil.NewTestConfig()
	dao, err := model.NewDao(cfg, testutil.FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	testutil.ApplyMigrations(t, dao)

	repos := adapter.NewRepositories(dao, testutil.FakeLogger{})
	svc := service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, repos, &routerFakeRegistry{}, testutil.FakeLogger{}, cfg, nil)
	deps := Dependencies{Service: svc}

	r := gin.New()
	api := r.Group("/api")
	setupCategoryRoutes(api, deps)
	RegisterTimeLogRoutes(api, deps)
	setupTaskRoutes(api, deps)
	setupMetricRoutes(api, deps)
	setupConstraintRoutes(api, deps)
	setupDatasourceRoutes(api, deps)

	return r, svc
}

func seedCategory(t *testing.T, svc *service.Service) {
	t.Helper()
	if err := svc.CreateCategory(&domain.Category{Name: "seed"}); err != nil {
		t.Fatalf("seed category: %v", err)
	}
}

func int32Ptr(v int32) *int32       { return &v }
func strPtr(v string) *string       { return &v }
func float64Ptr(v float64) *float64 { return &v }
