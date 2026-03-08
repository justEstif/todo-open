package app

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	apiclient "github.com/justEstif/todo-open/internal/client/api"
)

func TestNewServer_PersistsTasksAcrossRestart(t *testing.T) {
	workspace := t.TempDir()
	t.Setenv("TODOOPEN_WORKSPACE_ROOT", workspace)
	t.Setenv("TODOOPEN_STORE", "")

	srv1, err := NewServer(":0")
	if err != nil {
		t.Fatalf("new server #1: %v", err)
	}
	ts1 := httptest.NewServer(srv1.Handler)
	client1 := apiclient.New(ts1.URL)

	created, err := client1.CreateTask("persist me")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created task id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = srv1.Shutdown(ctx)
	ts1.Close()

	srv2, err := NewServer(":0")
	if err != nil {
		t.Fatalf("new server #2: %v", err)
	}
	ts2 := httptest.NewServer(srv2.Handler)
	t.Cleanup(ts2.Close)

	items, err := apiclient.New(ts2.URL).ListTasks()
	if err != nil {
		t.Fatalf("list after restart: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items after restart = %d, want 1", len(items))
	}
	if items[0].Title != "persist me" {
		t.Fatalf("title after restart = %q, want %q", items[0].Title, "persist me")
	}
}

func TestNewServer_FailsWhenEnabledPluginHandshakeFails(t *testing.T) {
	workspace := t.TempDir()
	t.Setenv("TODOOPEN_WORKSPACE_ROOT", workspace)

	metaDir := filepath.Join(workspace, ".todoopen")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		t.Fatalf("mkdir metadata dir: %v", err)
	}
	metaPayload := `{"workspace_version":1,"schema_version":"todo.open.task.v1"}`
	if err := os.WriteFile(filepath.Join(metaDir, "meta.json"), []byte(metaPayload), 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}
	cfgPayload := `
[views]
  enabled = ["json", "markdown"]

[sync]
  enabled = ["noop"]

[adapters.markdown]
  bin  = "sh"
  kind = "view"
  args = ["-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"wrong\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 1"]
`
	if err := os.WriteFile(filepath.Join(metaDir, "config.toml"), []byte(cfgPayload), 0o644); err != nil {
		t.Fatalf("write adapter config: %v", err)
	}

	_, err := NewServer(":0")
	if err == nil {
		t.Fatal("expected startup error")
	}
	if !strings.Contains(err.Error(), "adapter initialization failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
