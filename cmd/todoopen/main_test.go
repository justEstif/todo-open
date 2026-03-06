package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func TestTaskCRUDCommands(t *testing.T) {
	t.Parallel()

	repo := memory.NewTaskRepo()
	ids := []string{"task_1", "task_2", "task_3"}
	i := 0
	svc := core.NewService(repo, func() time.Time { return time.Date(2026, 3, 5, 20, 0, 0, 0, time.UTC) }, func() string {
		id := ids[i]
		i++
		return id
	})
	ts := httptest.NewServer(api.NewRouter(svc, adapters.Runtime{}))
	t.Cleanup(ts.Close)

	var out bytes.Buffer
	var errBuf bytes.Buffer

	code := run([]string{"task", "create", "--server", ts.URL, "--title", "first"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("create failed with code %d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "\"id\":\"task_1\"") {
		t.Fatalf("unexpected create output: %s", out.String())
	}

	out.Reset()
	errBuf.Reset()
	code = run([]string{"task", "list", "--server", ts.URL}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("list failed with code %d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "task_1\topen\tfirst") {
		t.Fatalf("unexpected list output: %s", out.String())
	}

	out.Reset()
	errBuf.Reset()
	code = run([]string{"task", "update", "--server", ts.URL, "--id", "task_1", "--title", "updated"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("update failed with code %d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "\"title\":\"updated\"") {
		t.Fatalf("unexpected update output: %s", out.String())
	}

	out.Reset()
	errBuf.Reset()
	code = run([]string{"task", "get", "--server", ts.URL, "--id", "task_1"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("get failed with code %d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "\"id\":\"task_1\"") {
		t.Fatalf("unexpected get output: %s", out.String())
	}

	out.Reset()
	errBuf.Reset()
	code = run([]string{"task", "delete", "--server", ts.URL, "--id", "task_1"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("delete failed with code %d, stderr=%s", code, errBuf.String())
	}
	if strings.TrimSpace(out.String()) != "deleted" {
		t.Fatalf("unexpected delete output: %s", out.String())
	}
}

func TestTaskCommandErrorPath(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var errBuf bytes.Buffer
	code := run([]string{"task", "get", "--server", "http://127.0.0.1:1", "--id", "task_1"}, &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(errBuf.String(), "get failed:") {
		t.Fatalf("unexpected stderr: %s", errBuf.String())
	}
}

func TestHelpCommand(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	var errBuf bytes.Buffer
	code := run([]string{"--help"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr: %s", errBuf.String())
	}
	if !strings.Contains(out.String(), "todoopen web") {
		t.Fatalf("help output missing web command: %s", out.String())
	}
	if !strings.Contains(out.String(), "todoopen --version") {
		t.Fatalf("help output missing version command: %s", out.String())
	}
}

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	oldVersion := version
	version = "v1.2.3"
	t.Cleanup(func() { version = oldVersion })

	var out bytes.Buffer
	var errBuf bytes.Buffer
	code := run([]string{"--version"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr: %s", errBuf.String())
	}
	if strings.TrimSpace(out.String()) != "todoopen v1.2.3" {
		t.Fatalf("unexpected version output: %s", out.String())
	}
}

func TestAdaptersCommand_ShowsSourceForBuiltins(t *testing.T) {
	t.Parallel()

	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, ".todoopen"), 0o755); err != nil {
		t.Fatalf("mkdir metadata dir: %v", err)
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	code := run([]string{"adapters", "--workspace", workspace, "--json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), `"source":"builtin"`) {
		t.Fatalf("expected builtin source in output: %s", out.String())
	}
}

func TestAdaptersCommand_ShowsSourceForPlugin(t *testing.T) {
	t.Parallel()

	workspace := t.TempDir()
	metaDir := filepath.Join(workspace, ".todoopen")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		t.Fatalf("mkdir metadata dir: %v", err)
	}
	payload := `{
  "workspace_version": 1,
  "schema_version": "todo.open.task.v1",
  "enabled_views": ["json", "markdown"],
  "enabled_sync_adapters": ["noop"],
  "adapter_plugins": [
    {"name":"markdown","kind":"view","command":"sh","args":["-c","printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"markdown\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 1"]}
  ]
}`
	if err := os.WriteFile(filepath.Join(metaDir, "meta.json"), []byte(payload), 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	code := run([]string{"adapters", "--workspace", workspace, "--json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d stderr=%s out=%s", code, errBuf.String(), out.String())
	}
	if !strings.Contains(out.String(), `"name":"markdown"`) || !strings.Contains(out.String(), `"source":"plugin"`) {
		t.Fatalf("expected plugin source in output: %s", out.String())
	}
}
