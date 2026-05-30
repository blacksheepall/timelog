package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "merge-openapi: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	restPath := filepath.Join("api", "openapi", "rest.yaml")
	schemaDir := filepath.Join("gen", "openapi", "schemas")
	outPath := filepath.Join("router", "docs", "openapi.yaml")

	restBytes, err := os.ReadFile(restPath)
	if err != nil {
		return fmt.Errorf("read rest spec: %w", err)
	}

	var spec map[string]any
	if err := yaml.Unmarshal(restBytes, &spec); err != nil {
		return fmt.Errorf("parse rest spec: %w", err)
	}

	components, _ := spec["components"].(map[string]any)
	if components == nil {
		components = map[string]any{}
		spec["components"] = components
	}
	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		schemas = map[string]any{}
		components["schemas"] = schemas
	}

	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("read schema dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "gen.schema.json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(schemaDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read schema %s: %w", entry.Name(), err)
		}

		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			return fmt.Errorf("parse schema %s: %w", entry.Name(), err)
		}

		definitions, _ := doc["definitions"].(map[string]any)
		if definitions == nil {
			continue
		}

		refMap := map[string]string{}
		for key := range definitions {
			refMap[key] = normalizeSchemaName(key)
		}

		for key, value := range definitions {
			name := refMap[key]
			schema, ok := value.(map[string]any)
			if !ok {
				continue
			}
			rewriteRefs(schema, refMap)
			if _, exists := schemas[name]; !exists {
				schemas[name] = schema
			}
		}
	}

	out, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("marshal openapi: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		return fmt.Errorf("write openapi: %w", err)
	}

	return nil
}

func normalizeSchemaName(key string) string {
	if idx := strings.LastIndex(key, "."); idx >= 0 {
		return key[idx+1:]
	}
	return key
}

func rewriteRefs(node any, refMap map[string]string) {
	switch value := node.(type) {
	case map[string]any:
		for k, v := range value {
			if k == "$ref" {
				if ref, ok := v.(string); ok {
					value[k] = rewriteRef(ref, refMap)
				}
				continue
			}
			rewriteRefs(v, refMap)
		}
	case []any:
		for _, item := range value {
			rewriteRefs(item, refMap)
		}
	}
}

func rewriteRef(ref string, refMap map[string]string) string {
	const prefix = "#/definitions/"
	if !strings.HasPrefix(ref, prefix) {
		return ref
	}
	key := strings.TrimPrefix(ref, prefix)
	name, ok := refMap[key]
	if !ok {
		name = normalizeSchemaName(key)
	}
	return "#/components/schemas/" + name
}
