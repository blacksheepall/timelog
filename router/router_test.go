package router

import (
	"bytes"
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/testutil"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

//go:embed all:web
var testStaticFS embed.FS

func newTestServiceForRouter(t *testing.T, cfg *config.Config) *service.Service {
	t.Helper()
	dao, err := model.NewDao(cfg, testutil.FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	testutil.ApplyMigrations(t, dao)
	repos := adapter.NewRepositories(dao)
	return service.NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, testutil.FakeLogger{}, cfg, nil)
}

func TestRegisterWithStaticFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.NewTestConfig()
	cfg.Server.Port = 0
	cfg.DevMode = true

	r := gin.New()
	deps := Dependencies{Service: newTestServiceForRouter(t, cfg)}

	registered := Register(r, cfg, testutil.FakeLogger{}, testStaticFS, deps)
	if registered == nil {
		t.Fatal("expected router")
	}

	// API 404
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/not-found", nil)
	registered.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	// SPA fallback
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	registered.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// static asset
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	registered.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRunServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.NewTestConfig()
	cfg.Server.Port = 0
	r := gin.New()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := RunServer(ctx, r, cfg, testutil.FakeLogger{}); err != nil {
		t.Fatalf("RunServer: %v", err)
	}
}

func TestDevLoginHandler(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDevLoginHandlerEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.NewTestConfig()
	cfg.DevMode = true
	svc := newTestServiceForRouter(t, cfg)
	deps := Dependencies{Service: svc}

	r := gin.New()
	api := r.Group("/api")
	setupAuthRoutes(api, cfg, deps)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
