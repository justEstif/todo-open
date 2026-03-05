package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"net/http/httptest"

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
	ts := httptest.NewServer(api.NewRouter(svc))
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
