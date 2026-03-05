package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	syncadapter "github.com/justEstif/todo-open/internal/sync"
	"github.com/justEstif/todo-open/internal/view"
)

// AdapterConfig captures user-facing adapter configuration.
type AdapterConfig struct {
	EnabledViews        []string                  `json:"enabled_views"`
	EnabledSyncAdapters []string                  `json:"enabled_sync_adapters"`
	ViewSettings        map[string]map[string]any `json:"view_settings,omitempty"`
	SyncSettings        map[string]map[string]any `json:"sync_settings,omitempty"`
}

// AdapterStatus represents runtime visibility for one adapter.
type AdapterStatus struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
}

// AdapterRuntime represents validated adapter configuration and startup health.
type AdapterRuntime struct {
	Config AdapterConfig   `json:"config"`
	Status []AdapterStatus `json:"status"`
	Ready  bool            `json:"ready"`
	Errors []string        `json:"errors,omitempty"`
}

func DefaultAdapterConfig() AdapterConfig {
	return AdapterConfig{
		EnabledViews:        []string{"json"},
		EnabledSyncAdapters: []string{"noop"},
	}
}

func LoadAdapterConfig(path string) (AdapterConfig, error) {
	if strings.TrimSpace(path) == "" {
		return DefaultAdapterConfig(), nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultAdapterConfig(), nil
	}
	if err != nil {
		return AdapterConfig{}, fmt.Errorf("read adapter config: %w", err)
	}

	cfg := DefaultAdapterConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AdapterConfig{}, fmt.Errorf("decode adapter config: %w", err)
	}
	if err := validateAdapterConfig(cfg); err != nil {
		return AdapterConfig{}, err
	}
	return cfg, nil
}

func BuildAdapterRuntime(cfg AdapterConfig, viewRegistry *view.Registry, syncRegistry *syncadapter.Registry) AdapterRuntime {
	runtime := AdapterRuntime{Config: cfg, Ready: true}

	enabledViews := asSet(cfg.EnabledViews)
	enabledSync := asSet(cfg.EnabledSyncAdapters)

	for _, name := range viewRegistry.Names() {
		enabled := enabledViews[name]
		runtime.Status = append(runtime.Status, AdapterStatus{Kind: "view", Name: name, Enabled: enabled, Healthy: enabled})
	}
	for _, name := range cfg.EnabledViews {
		if _, err := viewRegistry.Get(name); err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("view adapter %q is enabled but not registered", name))
			runtime.Status = append(runtime.Status, AdapterStatus{Kind: "view", Name: name, Enabled: true, Healthy: false, Message: "adapter is not registered"})
		}
	}

	for _, name := range syncRegistry.Names() {
		enabled := enabledSync[name]
		runtime.Status = append(runtime.Status, AdapterStatus{Kind: "sync", Name: name, Enabled: enabled, Healthy: enabled})
	}
	for _, name := range cfg.EnabledSyncAdapters {
		if _, err := syncRegistry.Get(name); err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("sync adapter %q is enabled but not registered", name))
			runtime.Status = append(runtime.Status, AdapterStatus{Kind: "sync", Name: name, Enabled: true, Healthy: false, Message: "adapter is not registered"})
		}
	}

	sort.Slice(runtime.Status, func(i, j int) bool {
		if runtime.Status[i].Kind == runtime.Status[j].Kind {
			return runtime.Status[i].Name < runtime.Status[j].Name
		}
		return runtime.Status[i].Kind < runtime.Status[j].Kind
	})

	return runtime
}

func validateAdapterConfig(cfg AdapterConfig) error {
	if err := validateNames("enabled_views", cfg.EnabledViews); err != nil {
		return err
	}
	if err := validateNames("enabled_sync_adapters", cfg.EnabledSyncAdapters); err != nil {
		return err
	}
	return nil
}

func validateNames(field string, names []string) error {
	seen := map[string]struct{}{}
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return fmt.Errorf("%s contains empty adapter name", field)
		}
		if _, ok := seen[trimmed]; ok {
			return fmt.Errorf("%s contains duplicate adapter name %q", field, trimmed)
		}
		seen[trimmed] = struct{}{}
	}
	return nil
}

func asSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}
