package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigPathEnv(t *testing.T) {
	want := "/tmp/test-config.yml"
	t.Setenv("TIMELOG_CONFIG_PATH", want)
	if got := ResolveConfigPath("config.yml"); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveConfigPathDefault(t *testing.T) {
	t.Setenv("TIMELOG_CONFIG_PATH", "")
	if got := ResolveConfigPath("config.yml"); got != "config.yml" {
		t.Fatalf("expected config.yml, got %q", got)
	}
}

func TestDatasourceConfig(t *testing.T) {
	ResetConfig()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
dev_mode: false
server:
  port: 8080
database:
  host: ":memory:"
log:
  level: info
datasources:
  - name: maimemo
    type: maimemo
    enabled: true
    config:
      token: "abc"
      endpoint: "https://open.maimemo.com"
  - name: disabled
    type: noop
    enabled: false
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("TIMELOG_CONFIG_PATH", path)

	cfg := GetConfig(path)
	if len(cfg.Datasources) != 2 {
		t.Fatalf("expected 2 datasources, got %d", len(cfg.Datasources))
	}
	ds := cfg.Datasources[0]
	if ds.Name != "maimemo" || ds.Type != "maimemo" || !ds.Enabled {
		t.Fatalf("unexpected first datasource: %+v", ds)
	}
	if ds.Config["token"] != "abc" || ds.Config["endpoint"] != "https://open.maimemo.com" {
		t.Fatalf("unexpected first datasource config: %+v", ds.Config)
	}
	if cfg.Datasources[1].Enabled {
		t.Fatalf("expected second datasource to be disabled")
	}

	ResetConfig()
}

func TestGetConfigAndReset(t *testing.T) {
	ResetConfig()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
dev_mode: true
server:
  port: 9999
database:
  host: ":memory:"
log:
  level: info
  path: "app.log"
  rotation:
    max_size: 1
    max_backups: 1
    max_age: 1
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("TIMELOG_CONFIG_PATH", path)

	cfg := GetConfig(path)
	if cfg == nil {
		t.Fatal("expected config")
	}
	if !cfg.DevMode || cfg.Server.Port != 9999 {
		t.Fatalf("unexpected config values: %+v", cfg)
	}

	SetConfig(&Config{DevMode: false})
	cfg2 := GetConfig(path)
	if cfg2.DevMode {
		t.Fatal("expected DevMode false after SetConfig")
	}

	ResetConfig()
}

func TestDotEnvLoadedIntoProcessEnv(t *testing.T) {
	ResetConfig()
	dir := t.TempDir()
	t.Chdir(dir) // GetConfig loads ".env" from the working directory
	t.Cleanup(func() {
		os.Unsetenv("TIMELOG_CONFIG_PATH")
		os.Unsetenv("MCP_TOKEN")
	})

	ymlPath := filepath.Join(dir, "custom.yml")
	yml := `
dev_mode: false
server:
  port: 8080
database:
  host: ":memory:"
log:
  level: info
`
	if err := os.WriteFile(ymlPath, []byte(yml), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	dotenv := "TIMELOG_CONFIG_PATH=" + ymlPath + "\nMCP_TOKEN=from-dotenv\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(dotenv), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	// The default path points at a missing file; only the .env-provided
	// TIMELOG_CONFIG_PATH can make this load succeed.
	cfg := GetConfig(filepath.Join(dir, "missing.yml"))
	if cfg.Server.Port != 8080 {
		t.Fatalf("expected config from .env-provided path, got port %d", cfg.Server.Port)
	}
	if got := os.Getenv("MCP_TOKEN"); got != "from-dotenv" {
		t.Fatalf("expected .env entry injected into process env, got %q", got)
	}

	ResetConfig()
}
