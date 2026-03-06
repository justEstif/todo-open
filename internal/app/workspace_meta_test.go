package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWorkspaceMetaDefaultWhenMissing(t *testing.T) {
	t.Parallel()

	meta, err := LoadWorkspaceMeta(t.TempDir())
	if err != nil {
		t.Fatalf("load metadata: %v", err)
	}
	if meta.WorkspaceVersion != 1 {
		t.Fatalf("workspace_version = %d, want 1", meta.WorkspaceVersion)
	}
	if got, want := meta.SchemaVersion, "todo.open.task.v1"; got != want {
		t.Fatalf("schema_version = %q, want %q", got, want)
	}
	if len(meta.EnabledViews) != 1 || meta.EnabledViews[0] != "json" {
		t.Fatalf("enabled views = %v, want [json]", meta.EnabledViews)
	}
	if len(meta.EnabledSyncAdapters) != 1 || meta.EnabledSyncAdapters[0] != "noop" {
		t.Fatalf("enabled sync adapters = %v, want [noop]", meta.EnabledSyncAdapters)
	}
}

func TestLoadWorkspaceMetaWithPlugins(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	metaDir := filepath.Join(root, ".todoopen")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		t.Fatalf("mkdir metadata dir: %v", err)
	}

	payload := `{
  "workspace_version": 1,
  "schema_version": "todo.open.task.v1",
  "enabled_views": ["json", "markdown"],
  "enabled_sync_adapters": ["noop", "git"],
  "adapter_plugins": [
    {"name":"markdown","kind":"view","command":"todoopen-plugin-view-markdown"},
    {"name":"git","kind":"sync","command":"todoopen-plugin-sync-git","args":["--mode","fast"]}
  ],
  "ext": {
    "adapter_settings": {
      "git": {"remote":"origin"}
    }
  }
}`
	if err := os.WriteFile(filepath.Join(metaDir, "meta.json"), []byte(payload), 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	meta, err := LoadWorkspaceMeta(root)
	if err != nil {
		t.Fatalf("load metadata: %v", err)
	}
	if len(meta.AdapterPlugins) != 2 {
		t.Fatalf("adapter_plugins len = %d, want 2", len(meta.AdapterPlugins))
	}
}

func TestLoadWorkspaceMetaValidationFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "unknown field",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","unknown":1}`,
		},
		{
			name:    "invalid kind",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","adapter_plugins":[{"name":"x","kind":"invalid","command":"plugin"}]}`,
		},
		{
			name:    "duplicate plugin by kind and name",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","adapter_plugins":[{"name":"x","kind":"view","command":"a"},{"name":"x","kind":"view","command":"b"}]}`,
		},
		{
			name:    "missing command",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","adapter_plugins":[{"name":"x","kind":"view","command":""}]}`,
		},
		{
			name:    "duplicate enabled view",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","enabled_views":["json","json"]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			metaDir := filepath.Join(root, ".todoopen")
			if err := os.MkdirAll(metaDir, 0o755); err != nil {
				t.Fatalf("mkdir metadata dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(metaDir, "meta.json"), []byte(tc.payload), 0o644); err != nil {
				t.Fatalf("write metadata: %v", err)
			}

			if _, err := LoadWorkspaceMeta(root); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
