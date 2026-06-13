package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/blacksheepaul/timelog/core/config"
	log "github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/internal/app"
	"github.com/blacksheepaul/timelog/service"
)

func resolveTTL(args []string, defaultTTL int) (int, error) {
	if len(args) > 1 {
		return 0, errors.New("too many positional arguments")
	}
	if len(args) == 0 {
		return defaultTTL, nil
	}

	ttl, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid ttl: %w", err)
	}
	if ttl < 0 {
		return 0, errors.New("ttl must be >= 0")
	}

	return ttl, nil
}

func usage() string {
	return "Usage: go run ./cmd/passkey-temp-admin [ttl]"
}

func runCreate(args []string, svc *service.Service, cfg *config.Config, stdout io.Writer) error {
	ttl, err := resolveTTL(args, cfg.Passkey.TempPassword.TTL)
	if err != nil {
		return err
	}

	record, password, err := svc.CreateTempPassword(ttl)
	if err != nil {
		return fmt.Errorf("failed to create temp password: %w", err)
	}

	fmt.Fprintf(stdout, "temp password: %s\n", password)
	fmt.Fprintf(stdout, "expires at: %s\n", record.ExpiresAt.Format("2006-01-02 15:04:05"))
	return nil
}

func main() {
	cfg := config.GetConfig("config.yml")
	logger := log.SetZapLogger(*cfg)

	application, err := app.New(cfg, logger, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	if err := runCreate(os.Args[1:], application.Service, cfg, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}
}
