package jsonl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

func mustNewTaskRepo(t *testing.T, dir string) *TaskRepo {
	t.Helper()
	r, err := NewTaskRepo(dir)
	if err != nil {
		t.Fatalf("NewTaskRepo: %v", err)
	}
	return r
}

func TestTaskRepoCRUDAndMetaBootstrap(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	repo := mustNewTaskRepo(t, dir)
	now := time.Date(2026, 3, 5, 20, 0, 0, 0, time.UTC)
	task := core.Task{ID: "task_1", Title: "one", Status: core.TaskStatusOpen, CreatedAt: now, UpdatedAt: now, Version: 1}

	created, err := repo.Create(context.Background(), task)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != "task_1" {
		t.Fatalf("unexpected id: %s", created.ID)
	}

	metaPath := filepath.Join(dir, ".todoopen", "meta.json")
	if _, err := os.Stat(metaPath); err != nil {
		t.Fatalf("expected metadata file: %v", err)
	}

	got, err := repo.GetByID(context.Background(), "task_1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "one" {
		t.Fatalf("unexpected title: %s", got.Title)
	}

	got.Title = "updated"
	got.Version = 2
	got.UpdatedAt = now.Add(time.Minute)
	if _, err := repo.Update(context.Background(), got); err != nil {
		t.Fatalf("update: %v", err)
	}

	items, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].Title != "updated" {
		t.Fatalf("unexpected list result: %+v", items)
	}
}

func TestTaskRepoRejectsCorruptJSONL(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	// Write a corrupt tasks file; workspace init must succeed since meta is absent.
	repo := mustNewTaskRepo(t, dir)
	if err := os.WriteFile(filepath.Join(dir, "tasks.jsonl"), []byte("{bad json}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.List(context.Background()); err == nil {
		t.Fatal("expected corruption error")
	}
}

func TestTaskRepoRejectsSchemaMismatch(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".todoopen"), 0o755); err != nil {
		t.Fatal(err)
	}
	badMeta := []byte(`{"workspace_version":1,"schema_version":"todo.open.task.v0"}`)
	if err := os.WriteFile(filepath.Join(dir, ".todoopen", "meta.json"), badMeta, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := NewTaskRepo(dir); err == nil {
		t.Fatal("expected schema mismatch error from NewTaskRepo")
	}
}
