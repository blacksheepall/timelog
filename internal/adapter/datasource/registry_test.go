package datasource

import (
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
)

func TestRegistry_NamesAndGet(t *testing.T) {
	reg, err := NewRegistry([]config.DatasourceConfig{
		{Name: "maimemo", Type: "maimemo", Enabled: true, Config: map[string]interface{}{
			"token": "test-token",
		}},
	})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}

	names := reg.Names()
	if len(names) != 1 || names[0] != "maimemo" {
		t.Fatalf("unexpected names: %v", names)
	}

	_, err = reg.Get("maimemo")
	if err != nil {
		t.Fatalf("Get maimemo: %v", err)
	}

	_, err = reg.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing datasource")
	}
}

func TestRegistry_DisabledSkipped(t *testing.T) {
	reg, err := NewRegistry([]config.DatasourceConfig{
		{Name: "maimemo", Type: "maimemo", Enabled: false, Config: map[string]interface{}{}},
	})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	if len(reg.Names()) != 0 {
		t.Fatalf("expected disabled source to be skipped")
	}
}

func TestRegistry_UnknownType(t *testing.T) {
	_, err := NewRegistry([]config.DatasourceConfig{
		{Name: "bad", Type: "unknown", Enabled: true, Config: map[string]interface{}{}},
	})
	if err == nil {
		t.Fatal("expected error for unknown datasource type")
	}
}
