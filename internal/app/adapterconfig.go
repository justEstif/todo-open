package app

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/plugin"
	syncadapter "github.com/justEstif/todo-open/internal/sync"
	"github.com/justEstif/todo-open/internal/view"
)

// Re-export canonical adapter runtime/config types for app consumers.
type AdapterConfig = adapters.Config
type AdapterStatus = adapters.Status
type AdapterRuntime = adapters.Runtime

func BuildAdapterRuntime(cfg AdapterConfig, viewRegistry *view.Registry, syncRegistry *syncadapter.Registry) AdapterRuntime {
	runtime := AdapterRuntime{Config: cfg, Ready: true}

	enabledViews := asSet(cfg.EnabledViews)
	enabledSync := asSet(cfg.EnabledSyncAdapters)

	for _, name := range viewRegistry.Names() {
		enabled := enabledViews[name]
		runtime.Status = append(runtime.Status, AdapterStatus{Kind: "view", Name: name, Source: "builtin", Enabled: enabled, Healthy: enabled})
	}
	for _, name := range cfg.EnabledViews {
		if _, err := viewRegistry.Get(name); err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("view adapter %q is enabled but not registered", name))
			runtime.Status = append(runtime.Status, AdapterStatus{Kind: "view", Name: name, Source: "unknown", Enabled: true, Healthy: false, Message: "adapter is not registered"})
		}
	}

	for _, name := range syncRegistry.Names() {
		enabled := enabledSync[name]
		runtime.Status = append(runtime.Status, AdapterStatus{Kind: "sync", Name: name, Source: "builtin", Enabled: enabled, Healthy: enabled})
	}
	for _, name := range cfg.EnabledSyncAdapters {
		if _, err := syncRegistry.Get(name); err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("sync adapter %q is enabled but not registered", name))
			runtime.Status = append(runtime.Status, AdapterStatus{Kind: "sync", Name: name, Source: "unknown", Enabled: true, Healthy: false, Message: "adapter is not registered"})
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

func BuildAdapterRuntimeFromMeta(meta WorkspaceMeta, viewRegistry *view.Registry, syncRegistry *syncadapter.Registry) AdapterRuntime {
	cfg := AdapterConfig{
		EnabledViews:        meta.EnabledViews,
		EnabledSyncAdapters: meta.EnabledSyncAdapters,
	}
	runtime := BuildAdapterRuntime(cfg, viewRegistry, syncRegistry)

	pluginByKey := map[string]AdapterPluginConfig{}
	for _, p := range meta.AdapterPlugins {
		pluginByKey[string(p.Kind)+":"+p.Name] = p
	}

	loader := plugin.NewLoader(2 * time.Second)
	for _, name := range meta.EnabledViews {
		if _, err := viewRegistry.Get(name); err == nil {
			continue
		}
		p, ok := pluginByKey[string(plugin.AdapterKindView)+":"+name]
		if !ok {
			continue
		}
		def := plugin.Definition{Name: p.Name, Kind: p.Kind, Command: p.Command, Args: p.Args}
		loaded, err := loader.Load(context.Background(), def)
		if err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("view plugin %q failed startup handshake: %v", name, err))
			upsertStatus(&runtime, AdapterStatus{Kind: "view", Name: name, Source: "plugin", Enabled: true, Healthy: false, Message: "plugin handshake failed"})
			continue
		}
		_ = loaded.Close()
		removeMissingRegistrationError(&runtime, "view", name)
		upsertStatus(&runtime, AdapterStatus{Kind: "view", Name: name, Source: "plugin", Enabled: true, Healthy: true, Message: "plugin handshake ok"})
	}

	for _, name := range meta.EnabledSyncAdapters {
		if _, err := syncRegistry.Get(name); err == nil {
			continue
		}
		p, ok := pluginByKey[string(plugin.AdapterKindSync)+":"+name]
		if !ok {
			continue
		}
		def := plugin.Definition{Name: p.Name, Kind: p.Kind, Command: p.Command, Args: p.Args}
		loaded, err := loader.Load(context.Background(), def)
		if err != nil {
			runtime.Ready = false
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("sync plugin %q failed startup handshake: %v", name, err))
			upsertStatus(&runtime, AdapterStatus{Kind: "sync", Name: name, Source: "plugin", Enabled: true, Healthy: false, Message: "plugin handshake failed"})
			continue
		}
		_ = loaded.Close()
		removeMissingRegistrationError(&runtime, "sync", name)
		upsertStatus(&runtime, AdapterStatus{Kind: "sync", Name: name, Source: "plugin", Enabled: true, Healthy: true, Message: "plugin handshake ok"})
	}

	sort.Slice(runtime.Status, func(i, j int) bool {
		if runtime.Status[i].Kind == runtime.Status[j].Kind {
			return runtime.Status[i].Name < runtime.Status[j].Name
		}
		return runtime.Status[i].Kind < runtime.Status[j].Kind
	})
	runtime.Ready = len(runtime.Errors) == 0

	return runtime
}

func upsertStatus(runtime *AdapterRuntime, status AdapterStatus) {
	for i := range runtime.Status {
		if runtime.Status[i].Kind == status.Kind && runtime.Status[i].Name == status.Name {
			runtime.Status[i] = status
			return
		}
	}
	runtime.Status = append(runtime.Status, status)
}

func removeMissingRegistrationError(runtime *AdapterRuntime, kind, name string) {
	prefix := fmt.Sprintf("%s adapter %q is enabled but not registered", kind, name)
	filtered := make([]string, 0, len(runtime.Errors))
	for _, msg := range runtime.Errors {
		if msg == prefix {
			continue
		}
		filtered = append(filtered, msg)
	}
	runtime.Errors = filtered
}

func asSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}
