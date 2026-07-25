package maimemo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestClient_GetTodayItems(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "get_today_items.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Fatalf("expected Bearer test-token, got %q", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	resp, err := client.GetTodayItems(context.Background())
	if err != nil {
		t.Fatalf("GetTodayItems: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if resp.Items[0].Value != 120 {
		t.Fatalf("expected first value 120, got %v", resp.Items[0].Value)
	}
}

func TestClient_GetTodayItems_StatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-token")
	_, err := client.GetTodayItems(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}
