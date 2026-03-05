package store

import (
	"context"

	"github.com/ebeyene/todo-open/internal/core"
)

type TaskRepository interface {
	Create(ctx context.Context, task core.Task) (core.Task, error)
}
