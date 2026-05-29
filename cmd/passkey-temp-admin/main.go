package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blacksheepaul/timelog/core/config"
	log "github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.GetConfig("config.yml")
	logger := log.SetZapLogger(*cfg)
	dao, err := model.NewDao(cfg, logger)
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	model.RegisterDao(dao, logger)
	service.InitService(logger, cfg, dao)

	command := strings.ToLower(os.Args[1])
	switch command {
	case "create":
		ttl := cfg.Passkey.TempPassword.TTL
		if len(os.Args) >= 3 {
			value, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("invalid ttl: %v\n", err)
				os.Exit(1)
			}
			ttl = value
		}
		record, password, err := service.CreateTempPassword(ttl)
		if err != nil {
			fmt.Printf("failed to create temp password: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("temp password: %s\n", password)
		fmt.Printf("expires at: %s\n", record.ExpiresAt.Format("2006-01-02 15:04:05"))
	case "list":
		passwords, err := service.ListTempPasswords()
		if err != nil {
			fmt.Printf("failed to list temp passwords: %v\n", err)
			os.Exit(1)
		}
		if len(passwords) == 0 {
			fmt.Println("no temp passwords found")
			return
		}
		for _, password := range passwords {
			fmt.Printf("id: %d\t expires_at: %s\n", password.ID, password.ExpiresAt.Format("2006-01-02 15:04:05"))
		}
	case "revoke":
		if len(os.Args) < 3 {
			fmt.Println("revoke requires an id")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("invalid id: %v\n", err)
			os.Exit(1)
		}
		if err := service.DeleteTempPassword(uint(id)); err != nil {
			fmt.Printf("failed to revoke temp password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("revoked")
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run ./cmd/passkey-temp-admin <create|list|revoke> [ttl|id]")
}
