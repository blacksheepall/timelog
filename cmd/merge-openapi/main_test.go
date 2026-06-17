package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunMergesSchemas(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir("/Users/gear/dev/timelog")
	}()

	if err := os.MkdirAll(filepath.Join("api", "openapi"), 0o755); err != nil {
		t.Fatalf("mkdir api: %v", err)
	}
	if err := os.MkdirAll(filepath.Join("gen", "openapi", "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	if err := os.MkdirAll(filepath.Join("router", "docs"), 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}

	rest := `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths: {}
`
	if err := os.WriteFile(filepath.Join("api", "openapi", "rest.yaml"), []byte(rest), 0o644); err != nil {
		t.Fatalf("write rest: %v", err)
	}

	schema := `{
  "definitions": {
    "timelog.v1.Timelog": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "self": { "$ref": "#/definitions/timelog.v1.Timelog" }
      }
    }
  }
}
`
	if err := os.WriteFile(filepath.Join("gen", "openapi", "schemas", "timelog.gen.schema.json"), []byte(schema), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}

	out := filepath.Join("router", "docs", "openapi.yaml")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !contains(string(data), "Timelog") {
		t.Fatalf("expected Timelog schema in output")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestRunMissingRestSpec(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir("/Users/gear/dev/timelog")
	}()

	if err := os.MkdirAll(filepath.Join("gen", "openapi", "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}

	if err := run(); err == nil {
		t.Fatal("expected error for missing rest spec")
	}
}

func TestRunInvalidRestSpec(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir("/Users/gear/dev/timelog")
	}()

	if err := os.MkdirAll(filepath.Join("api", "openapi"), 0o755); err != nil {
		t.Fatalf("mkdir api: %v", err)
	}
	if err := os.MkdirAll(filepath.Join("gen", "openapi", "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	if err := os.WriteFile(filepath.Join("api", "openapi", "rest.yaml"), []byte("bad: ["), 0o644); err != nil {
		t.Fatalf("write rest: %v", err)
	}

	if err := run(); err == nil {
		t.Fatal("expected error for invalid rest spec")
	}
}

func TestRunInvalidSchemaJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir("/Users/gear/dev/timelog")
	}()

	if err := os.MkdirAll(filepath.Join("api", "openapi"), 0o755); err != nil {
		t.Fatalf("mkdir api: %v", err)
	}
	if err := os.MkdirAll(filepath.Join("gen", "openapi", "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	if err := os.WriteFile(filepath.Join("api", "openapi", "rest.yaml"), []byte("openapi: 3.0.0\ninfo:\n  title: T\n  version: 1.0.0\npaths: {}\n"), 0o644); err != nil {
		t.Fatalf("write rest: %v", err)
	}
	if err := os.WriteFile(filepath.Join("gen", "openapi", "schemas", "bad.gen.schema.json"), []byte("not json"), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	if err := run(); err == nil {
		t.Fatal("expected error for invalid schema json")
	}
}

func TestNormalizeSchemaName(t *testing.T) {
	if got := normalizeSchemaName("foo.bar.Baz"); got != "Baz" {
		t.Fatalf("expected Baz, got %s", got)
	}
	if got := normalizeSchemaName("Baz"); got != "Baz" {
		t.Fatalf("expected Baz, got %s", got)
	}
}

func TestRewriteRef(t *testing.T) {
	refMap := map[string]string{"timelog.v1.Timelog": "Timelog"}
	if got := rewriteRef("#/definitions/timelog.v1.Timelog", refMap); got != "#/components/schemas/Timelog" {
		t.Fatalf("unexpected ref: %s", got)
	}
	if got := rewriteRef("#/other", refMap); got != "#/other" {
		t.Fatalf("expected unchanged ref: %s", got)
	}
}
