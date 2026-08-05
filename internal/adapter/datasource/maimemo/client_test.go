package maimemo

import (
	"context"
	"io"
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
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != getTodayItemsPath {
			t.Fatalf("expected path %s, got %s", getTodayItemsPath, r.URL.Path)
		}
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
	if len(resp.TodayItems) != 3 {
		t.Fatalf("expected 3 items, got %d", len(resp.TodayItems))
	}
	if resp.TodayItems[0].VocSpelling != "apple" || !resp.TodayItems[0].IsNew {
		t.Fatalf("unexpected first item: %+v", resp.TodayItems[0])
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

func TestClient_GetTodayItems_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"errors":[{"code":"common_not_found","msg":"Resource or Api not found","info":""}],"success":false}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	_, err := client.GetTodayItems(context.Background())
	if err == nil {
		t.Fatal("expected error when success is false")
	}
}
