package store

import (
	"context"
	"time"

	"github.com/ebeyene/todo-open/internal/core"
)

type TaskRepository interface {
	Create(ctx context.Context, task core.Task) (core.Task, error)
	GetByID(ctx context.Context, id string) (core.Task, error)
	List(ctx context.Context) ([]core.Task, error)
	Update(ctx context.Context, task core.Task) (core.Task, error)
	Delete(ctx context.Context, id string, deletedAt time.Time) error
}
