package router

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

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
	svc := service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, repos, nil, testutil.FakeLogger{}, cfg, nil)
	deps := Dependencies{Service: svc}

	r := gin.New()
	api := r.Group("/api")
	setupCategoryRoutes(api, deps)
	RegisterTimeLogRoutes(api, deps)
	setupTaskRoutes(api, deps)
	setupMetricRoutes(api, deps)
	setupConstraintRoutes(api, deps)

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
