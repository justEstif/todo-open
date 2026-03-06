package app

import (
	"context"
	"testing"

	"github.com/justEstif/todo-open/internal/plugin"
)

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

func TestBuildAdapterRuntimeFromMeta_ValidPluginHandshake(t *testing.T) {
	viewRegistry, err := NewViewRegistry()
	if err != nil {
		t.Fatalf("new view registry: %v", err)
	}
	syncRegistry, err := NewSyncRegistry()
	if err != nil {
		t.Fatalf("new sync registry: %v", err)
	}

	meta := WorkspaceMeta{
		WorkspaceVersion:    1,
		SchemaVersion:       "todo.open.task.v1",
		EnabledViews:        []string{"json", "markdown"},
		EnabledSyncAdapters: []string{"noop"},
		AdapterPlugins: []AdapterPluginConfig{{
			Name:    "markdown",
			Kind:    plugin.AdapterKindView,
			Command: "sh",
			Args:    []string{"-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"markdown\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 1"},
		}},
	}

	runtime := BuildAdapterRuntimeFromMeta(context.Background(), meta, viewRegistry, syncRegistry)
	if !runtime.Ready {
		t.Fatalf("runtime should be ready, errors=%v", runtime.Errors)
	}

	found := false
	for _, st := range runtime.Status {
		if st.Kind == "view" && st.Name == "markdown" {
			found = true
			if !st.Healthy || !st.Enabled {
				t.Fatalf("status=%+v", st)
			}
		}
	}
	if !found {
		t.Fatal("expected markdown plugin status entry")
	}
}

func TestBuildAdapterRuntimeFromMeta_PluginHandshakeFailure(t *testing.T) {
	viewRegistry, err := NewViewRegistry()
	if err != nil {
		t.Fatalf("new view registry: %v", err)
	}
	syncRegistry, err := NewSyncRegistry()
	if err != nil {
		t.Fatalf("new sync registry: %v", err)
	}

	meta := WorkspaceMeta{
		WorkspaceVersion:    1,
		SchemaVersion:       "todo.open.task.v1",
		EnabledViews:        []string{"json", "markdown"},
		EnabledSyncAdapters: []string{"noop"},
		AdapterPlugins: []AdapterPluginConfig{{
			Name:    "markdown",
			Kind:    plugin.AdapterKindView,
			Command: "sh",
			Args:    []string{"-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"wrong\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 1"},
		}},
	}

	runtime := BuildAdapterRuntimeFromMeta(context.Background(), meta, viewRegistry, syncRegistry)
	if runtime.Ready {
		t.Fatal("runtime should not be ready")
	}
	if len(runtime.Errors) == 0 {
		t.Fatal("expected startup diagnostics errors")
	}
}
