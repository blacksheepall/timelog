package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/blacksheepaul/timelog/core/config"
	log "github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/router"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var staticFiles embed.FS

func main() {
	cfg := config.GetConfig("config.yml")
	logger := log.SetZapLogger(*cfg)

	service.InitService(logger, cfg)
	model.InitDao(cfg, logger)
	if cfg.Passkey.Enabled {
		if err := service.InitWebAuthn(); err != nil {
			panic("Failed to initialize WebAuthn: " + err.Error())
		}
	}

	r := router.Register(gin.New(), cfg, logger, staticFiles)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Get server address
	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	protocol := "http"
	if cfg.Server.HTTPSEnabled {
		protocol = "https"
	}

	wg.Add(1)
	go router.LaunchServer(ctx, &wg, r, cfg)

	byebye := make(chan os.Signal, 1) // Listen for system signal，such as SIGINT, SIGTERM
	signal.Notify(byebye, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Server started, press Ctrl+C to stop.")

	fmt.Println("Program is running ...")
	fmt.Printf("Server is running at %s://%s\n", protocol, addr)
	logger.Info("Program is running, waiting for termination signal...")
	someonesaidbye := <-byebye // waiting for signal
	logger.Info("Received signal: %s, shutting down...", someonesaidbye)

	cancel() // tell other goroutines to stop
	logger.Info("Program exited gracefully.")
	fmt.Println("Program exited gracefully.")
}
