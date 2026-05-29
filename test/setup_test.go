package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/service"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func TestMain(m *testing.M) {
	cfg := loadTestConfig()
	gin.SetMode(gin.DebugMode)
	serviceForTest(cfg, FakeLogger{})
	daoForTest(cfg, FakeLogger{})

	if cfg.Test.Flush {
		flushDb()
	}

	os.Exit(m.Run())
}

func loadTestConfig() *config.Config {
	configPath := config.ResolveConfigPath("../config-test.yml")
	if _, err := os.Stat(configPath); err == nil {
		return config.GetConfig(configPath)
	}

	cfg := config.GetConfig("../config-example.yml")
	cfg.Database.Host = "./test.db"
	cfg.Server.HTTPSEnabled = false
	cfg.Passkey.Enabled = false
	cfg.Log.Path = "./test.log"
	cfg.Test.Flush = true
	return cfg
}

func flushDb() {
	migrationFiles, err := filepath.Glob("../model/migrations/*.up.sql")
	if err != nil {
		panic(err)
	}

	rawDB := model.GetDao().RawDB
	if rawDB == nil {
		panic("raw database is nil")
	}

	if err := dropAllTables(rawDB); err != nil {
		panic(fmt.Errorf("failed to drop tables: %v", err))
	}

	for _, migrationFile := range migrationFiles {
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			panic(fmt.Errorf("failed to read migration %s: %v", migrationFile, err))
		}
		if _, err := rawDB.Exec(string(content)); err != nil {
			panic(fmt.Errorf("failed to apply migration %s: %v", migrationFile, err))
		}
	}
}

func dropAllTables(db *sql.DB) error {
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tableNames = append(tableNames, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}
	defer db.Exec("PRAGMA foreign_keys = ON")

	for _, name := range tableNames {
		if _, err := db.Exec("DROP TABLE IF EXISTS " + name); err != nil {
			return err
		}
	}

	return nil
}

func daoForTest(cfg *config.Config, logi logger.Logger) {
	dao, err := model.NewDao(cfg, logi)
	if err != nil {
		panic(err)
	}
	model.RegisterDao(dao, logi)
}

func serviceForTest(cfg *config.Config, logi logger.Logger) {
	service.InitService(logi, cfg, model.GetDao())
}

type FakeLogger struct{}

func (l FakeLogger) Debug(fields ...interface{}) {
	fmt.Println(fields...)
}

func (l FakeLogger) Debugw(msg string, keysAndValues ...interface{}) {
	fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l FakeLogger) Info(fields ...interface{}) {
	fmt.Println(fields...)
}

func (l FakeLogger) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l FakeLogger) Warn(fields ...interface{}) {
	fmt.Println(fields...)
}

func (l FakeLogger) Warnw(msg string, keysAndValues ...interface{}) {
	fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l FakeLogger) Error(fields ...interface{}) {
	fmt.Println(fields...)
}

func (l FakeLogger) Errorw(msg string, keysAndValues ...interface{}) {
	fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
}

func (l FakeLogger) Fatal(fields ...interface{}) {
	fmt.Println(fields...)
}

func (l FakeLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	fmt.Println(append([]interface{}{msg}, keysAndValues...)...)
}
