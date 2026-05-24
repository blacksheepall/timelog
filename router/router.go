package router

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"sync"
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
var log logger.Logger
var appConfig *config.Config

func Register(r *gin.Engine, cfg *config.Config, l logger.Logger, staticFiles embed.FS) *gin.Engine {
	log = l
	appConfig = cfg

	r.Use(GinLogger)
	r.Use(middleware.Cors(cfg))

	api := r.Group("/api")
	protected := api.Group("")
	if cfg.Passkey.Enabled {
		protected.Use(middleware.Auth())
	}

	// 注册 TimeLog 路由
	RegisterTimeLogRoutes(protected)

	// 注册 Task 路由
	setupTaskRoutes(protected)

	// 注册 Constraint 路由
	setupConstraintRoutes(protected)

	// 注册 Passkey 路由（仅当 passkey 功能启用时）
	if cfg.Passkey.Enabled {
		setupPasskeyRoutes(api, protected)
	}

	// 注册 Swagger 文档路由（仅非 prod 构建）
	setupSwagger(r)

	// 静态文件服务 - 嵌入的Vue前端
	distFS, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem", err)
	}

	// 创建assets子目录的文件系统
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		log.Fatal("Failed to create assets sub filesystem", err)
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

func LaunchServer(ctx context.Context, wg *sync.WaitGroup, r *gin.Engine, cfg *config.Config) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	log.Info("[Startup] Server is starting...")
	log.Info(fmt.Sprintf("[Startup] Listen address: %s", addr))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		var err error
		if cfg.Server.HTTPSEnabled {
			err = srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server startup failed: %v\n", err)
			log.Fatal("Server startup failed", err)
		}
	}()

	<-ctx.Done()
	log.Info("Server received stop signal, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server shutdown failed", err)
	}
	log.Info("Server exited gracefully.")
}
