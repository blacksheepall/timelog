package router

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/router/middleware"
	"github.com/gin-gonic/gin"
)

var GinLogger = gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
	return fmt.Sprintf(`{"time":"%s","client":"%s","method":"%s","path":"%s","latency":"%s","status":%d,"emsg":"%s"}`+"\n",
		p.TimeStamp.Format(time.DateTime),
		p.ClientIP,
		p.Method,
		p.Path,
		p.Latency,
		p.StatusCode,
		p.ErrorMessage,
	)
})

func Register(r *gin.Engine, cfg *config.Config, l logger.Logger, staticFiles embed.FS, deps Dependencies) *gin.Engine {
	r.Use(GinLogger)
	r.Use(middleware.Cors(cfg))

	api := r.Group("/api")
	protected := api.Group("")
	if cfg.Passkey.Enabled {
		protected.Use(middleware.Auth(deps.Service))
	}

	// 注册 TimeLog 路由
	RegisterTimeLogRoutes(protected, deps)

	// 注册 Category 路由
	setupCategoryRoutes(protected, deps)

	// 注册 Task 路由
	setupTaskRoutes(protected, deps)

	// 注册 Constraint 路由
	setupConstraintRoutes(protected, deps)

	// 注册通用鉴权路由
	setupAuthRoutes(api, cfg, deps)

	// 注册 Passkey 路由（仅当 passkey 功能启用时）
	if cfg.Passkey.Enabled {
		setupPasskeyRoutes(api, protected, cfg, deps)
	}

	// 注册 API 文档路由（仅非 prod 构建）
	setupDocs(r)

	// 静态文件服务 - 嵌入的Vue前端
	distFS, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		l.Fatal("Failed to create sub filesystem", err)
	}
	if _, err := fs.Stat(distFS, "index.html"); err != nil {
		l.Fatal("Embedded frontend is missing index.html; run `make web` or `make buildx` first", err)
	}

	// 创建assets子目录的文件系统
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		l.Fatal("Failed to create assets sub filesystem", err)
	}

	// 服务静态资源文件 (JS, CSS, images等)
	r.StaticFS("/assets", http.FS(assetsFS))
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(distFS))
	r.StaticFileFS("/vite.svg", "vite.svg", http.FS(distFS))

	// SPA路由 - 对于所有未匹配的路由返回index.html
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:5] == "/api/" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// 其他所有路由返回index.html给Vue Router处理
		indexData, err := distFS.Open("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load frontend"})
			return
		}
		defer indexData.Close()

		c.DataFromReader(http.StatusOK, -1, "text/html; charset=utf-8", indexData, nil)
	})

	return r
}

func RunServer(ctx context.Context, r *gin.Engine, cfg *config.Config, l logger.Logger) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	l.Info("[Startup] Server is starting...")
	l.Info(fmt.Sprintf("[Startup] Listen address: %s", addr))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		var err error
		if cfg.Server.HTTPSEnabled {
			err = srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		l.Info("Server received stop signal, shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		l.Info("Server exited gracefully.")
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
