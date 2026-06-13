package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/blacksheepaul/timelog/core/config"
	log "github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/internal/app"
	"github.com/blacksheepaul/timelog/router"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var staticFiles embed.FS

func main() {
	cfg := config.GetConfig("config.yml")
	logger := log.SetZapLogger(*cfg)

	application, err := app.New(cfg, logger, nil)
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}

	if cfg.Passkey.Enabled {
		webAuthn, err := service.InitWebAuthnWithConfig(cfg)
		if err != nil {
			panic("Failed to initialize WebAuthn: " + err.Error())
		}
		application.Service.SetWebAuthn(webAuthn)
	}

	r := router.Register(gin.New(), cfg, logger, staticFiles, router.Dependencies{
		Service: application.Service,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get server address
	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	protocol := "http"
	if cfg.Server.HTTPSEnabled {
		protocol = "https"
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- router.RunServer(ctx, r, cfg, logger)
	}()

	byebye := make(chan os.Signal, 1) // Listen for system signal，such as SIGINT, SIGTERM
	signal.Notify(byebye, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Server started, press Ctrl+C to stop.")

	fmt.Println("Program is running ...")
	fmt.Printf("Server is running at %s://%s\n", protocol, addr)
	logger.Info("Program is running, waiting for termination signal...")

	select {
	case someonesaidbye := <-byebye:
		logger.Info("Received signal: %s, shutting down...", someonesaidbye)
		cancel()
		if err := <-errCh; err != nil {
			slog.Error("server exited with error", "err", err)
			os.Exit(1)
		}
	case err := <-errCh:
		if err != nil {
			slog.Error("server startup failed", "err", err)
			os.Exit(1)
		}
	}

	logger.Info("Program exited gracefully.")
	fmt.Println("Program exited gracefully.")
}
