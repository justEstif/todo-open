package store_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/jsonl"
	"github.com/justEstif/todo-open/internal/store/memory"
)

type repoFactory struct {
	name string
	new  func(t *testing.T) core.TaskRepository
}

func TestTaskRepositoryCreateRejectsDuplicateID(t *testing.T) {
	factories := []repoFactory{
		{name: "memory", new: func(_ *testing.T) core.TaskRepository { return memory.NewTaskRepo() }},
		{name: "jsonl", new: func(t *testing.T) core.TaskRepository {
			r, err := jsonl.NewTaskRepo(t.TempDir())
			if err != nil {
				t.Fatalf("NewTaskRepo: %v", err)
			}
			return r
		}},
	}

	for _, factory := range factories {
		t.Run(factory.name, func(t *testing.T) {
			repo := factory.new(t)
			ctx := context.Background()
			now := time.Date(2026, 3, 5, 20, 0, 0, 0, time.UTC)
			task := core.Task{ID: "task_1", Title: "first", Status: core.TaskStatusOpen, CreatedAt: now, UpdatedAt: now, Version: 1}

			if _, err := repo.Create(ctx, task); err != nil {
				t.Fatalf("first create error = %v", err)
			}
			if _, err := repo.Create(ctx, task); !errors.Is(err, core.ErrInvalidInput) {
				t.Fatalf("second create error = %v, want ErrInvalidInput", err)
			}
		})
	}
}
