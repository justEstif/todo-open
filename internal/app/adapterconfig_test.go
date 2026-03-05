package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAdapterConfig_DefaultWhenMissing(t *testing.T) {
	cfg, err := LoadAdapterConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if len(cfg.EnabledViews) != 1 || cfg.EnabledViews[0] != "json" {
		t.Fatalf("enabled views = %v, want [json]", cfg.EnabledViews)
	}
	if len(cfg.EnabledSyncAdapters) != 1 || cfg.EnabledSyncAdapters[0] != "noop" {
		t.Fatalf("enabled sync adapters = %v, want [noop]", cfg.EnabledSyncAdapters)
	}
}

func TestLoadAdapterConfig_Validation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "adapters.json")
	if err := os.WriteFile(path, []byte(`{"enabled_views":["json","json"]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadAdapterConfig(path); err == nil {
		t.Fatal("expected duplicate validation error")
	}
}

func TestBuildAdapterRuntime(t *testing.T) {
	viewRegistry, err := NewViewRegistry()
	if err != nil {
		t.Fatalf("new view registry: %v", err)
	}
	syncRegistry, err := NewSyncRegistry()
	if err != nil {
		t.Fatalf("new sync registry: %v", err)
	}

	runtime := BuildAdapterRuntime(AdapterConfig{
		EnabledViews:        []string{"json", "missing-view"},
		EnabledSyncAdapters: []string{"noop", "missing-sync"},
	}, viewRegistry, syncRegistry)

	if runtime.Ready {
		t.Fatal("runtime should not be ready")
	}
	if len(runtime.Errors) != 2 {
		t.Fatalf("errors = %v, want 2", runtime.Errors)
	}
}
