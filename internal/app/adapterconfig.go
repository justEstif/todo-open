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

const defaultPluginHandshakeTimeout = 2 * time.Second

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

func BuildAdapterRuntimeFromMeta(ctx context.Context, meta WorkspaceMeta, viewRegistry *view.Registry, syncRegistry *syncadapter.Registry) AdapterRuntime {
	resolver := newAdapterRuntimeResolver(viewRegistry, syncRegistry, newPluginProbe(defaultPluginHandshakeTimeout))
	return resolver.Resolve(ctx, meta)
}

type pluginProbe struct {
	loader *plugin.Loader
}

func newPluginProbe(handshakeTimeout time.Duration) pluginProbe {
	return pluginProbe{loader: plugin.NewLoader(handshakeTimeout)}
}

func (p pluginProbe) probe(ctx context.Context, cfg AdapterPluginConfig) error {
	def := plugin.Definition{Name: cfg.Name, Kind: cfg.Kind, Command: cfg.Command, Args: cfg.Args}
	loaded, err := p.loader.Load(ctx, def)
	if err != nil {
		return err
	}
	return loaded.Close()
}

type adapterRuntimeResolver struct {
	viewRegistry *view.Registry
	syncRegistry *syncadapter.Registry
	probe        pluginProbe
}

func newAdapterRuntimeResolver(viewRegistry *view.Registry, syncRegistry *syncadapter.Registry, probe pluginProbe) adapterRuntimeResolver {
	return adapterRuntimeResolver{viewRegistry: viewRegistry, syncRegistry: syncRegistry, probe: probe}
}

func (r adapterRuntimeResolver) Resolve(ctx context.Context, meta WorkspaceMeta) AdapterRuntime {
	cfg := AdapterConfig{
		EnabledViews:        meta.EnabledViews,
		EnabledSyncAdapters: meta.EnabledSyncAdapters,
	}
	runtime := BuildAdapterRuntime(cfg, r.viewRegistry, r.syncRegistry)

	pluginByKey := map[string]AdapterPluginConfig{}
	for _, p := range meta.AdapterPlugins {
		pluginByKey[string(p.Kind)+":"+p.Name] = p
	}

	r.resolvePluginBacked(ctx, &runtime, pluginByKey, plugin.AdapterKindView, meta.EnabledViews)
	r.resolvePluginBacked(ctx, &runtime, pluginByKey, plugin.AdapterKindSync, meta.EnabledSyncAdapters)

	sort.Slice(runtime.Status, func(i, j int) bool {
		if runtime.Status[i].Kind == runtime.Status[j].Kind {
			return runtime.Status[i].Name < runtime.Status[j].Name
		}
		return runtime.Status[i].Kind < runtime.Status[j].Kind
	})
	runtime.Ready = len(runtime.Errors) == 0

	return runtime
}

func (r adapterRuntimeResolver) resolvePluginBacked(ctx context.Context, runtime *AdapterRuntime, pluginByKey map[string]AdapterPluginConfig, kind plugin.AdapterKind, enabled []string) {
	kindText := string(kind)
	for _, name := range enabled {
		if r.isBuiltinRegistered(kind, name) {
			continue
		}

		p, ok := pluginByKey[kindText+":"+name]
		if !ok {
			continue
		}

		if err := r.probe.probe(ctx, p); err != nil {
			runtime.Errors = append(runtime.Errors, fmt.Sprintf("%s plugin %q failed startup handshake: %v", kindText, name, err))
			upsertStatus(runtime, AdapterStatus{Kind: kindText, Name: name, Source: "plugin", Enabled: true, Healthy: false, Message: "plugin handshake failed"})
			continue
		}

		removeMissingRegistrationError(runtime, kindText, name)
		upsertStatus(runtime, AdapterStatus{Kind: kindText, Name: name, Source: "plugin", Enabled: true, Healthy: true, Message: "plugin handshake ok"})
	}
}

func (r adapterRuntimeResolver) isBuiltinRegistered(kind plugin.AdapterKind, name string) bool {
	switch kind {
	case plugin.AdapterKindView:
		_, err := r.viewRegistry.Get(name)
		return err == nil
	case plugin.AdapterKindSync:
		_, err := r.syncRegistry.Get(name)
		return err == nil
	default:
		return false
	}
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
