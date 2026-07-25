package maimemo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

func TestDataSource_Fetch(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "get_today_items.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	src, err := NewDataSource("maimemo", config.DatasourceConfig{
		Name:    "maimemo",
		Type:    "maimemo",
		Enabled: true,
		Config: map[string]interface{}{
			"token":    "test-token",
			"endpoint": server.URL,
		},
	})
	if err != nil {
		t.Fatalf("NewDataSource: %v", err)
	}

	points, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if src.Name() != "maimemo" {
		t.Fatalf("expected name maimemo, got %q", src.Name())
	}
}

func TestDataSource_New_MissingToken(t *testing.T) {
	_, err := NewDataSource("maimemo", config.DatasourceConfig{
		Config: map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}
