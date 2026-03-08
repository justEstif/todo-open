package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/justEstif/todo-open/internal/plugin"
)

func TestLoadAdapterFileConfigDefaultWhenMissing(t *testing.T) {
	t.Parallel()

	cfg, err := LoadAdapterFileConfig(t.TempDir())
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if len(cfg.Views.Enabled) != 1 || cfg.Views.Enabled[0] != "json" {
		t.Fatalf("views.enabled = %v, want [json]", cfg.Views.Enabled)
	}
	if len(cfg.Sync.Enabled) != 1 || cfg.Sync.Enabled[0] != "noop" {
		t.Fatalf("sync.enabled = %v, want [noop]", cfg.Sync.Enabled)
	}
}

func TestLoadAdapterFileConfigWithPlugin(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".todoopen"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	payload := `
[views]
  enabled = ["json", "markdown"]

[sync]
  enabled = ["noop"]

[adapters.markdown]
  bin  = "todoopen-plugin-view-markdown"
  kind = "view"

[adapters.markdown.config]
  theme = "dark"
`
	if err := os.WriteFile(filepath.Join(root, ".todoopen", "config.toml"), []byte(payload), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadAdapterFileConfig(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	entry, ok := cfg.Adapters["markdown"]
	if !ok {
		t.Fatal("expected markdown adapter entry")
	}
	if entry.Kind != plugin.AdapterKindView {
		t.Fatalf("kind = %q, want view", entry.Kind)
	}
	if entry.Config["theme"] != "dark" {
		t.Fatalf("config.theme = %q, want dark", entry.Config["theme"])
	}
}

func TestLoadAdapterFileConfigEnvExpansion(t *testing.T) {
	// t.Setenv cannot be combined with t.Parallel
	t.Setenv("MY_TOKEN", "secret123")

	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".todoopen"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	payload := `
[sync]
  enabled = ["noop", "remote"]

[adapters.remote]
  bin  = "todoopen-plugin-sync-remote"
  kind = "sync"

[adapters.remote.config]
  token = "${MY_TOKEN}"
`
	if err := os.WriteFile(filepath.Join(root, ".todoopen", "config.toml"), []byte(payload), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadAdapterFileConfig(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Adapters["remote"].Config["token"] != "secret123" {
		t.Fatalf("token = %q, want secret123", cfg.Adapters["remote"].Config["token"])
	}
}

func TestLoadAdapterFileConfigValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "missing bin",
			payload: `
[views]
  enabled = ["markdown"]
[adapters.markdown]
  kind = "view"
`,
		},
		{
			name: "invalid kind",
			payload: `
[adapters.x]
  bin  = "some-bin"
  kind = "unknown"
`,
		},
		{
			name: "enabled adapter not declared",
			payload: `
[sync]
  enabled = ["noop", "ghost"]
`,
		},
		{
			name: "duplicate enabled view",
			payload: `
[views]
  enabled = ["json", "json"]
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".todoopen"), 0o755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, ".todoopen", "config.toml"), []byte(tc.payload), 0o644); err != nil {
				t.Fatalf("write config: %v", err)
			}

			if _, err := LoadAdapterFileConfig(root); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
