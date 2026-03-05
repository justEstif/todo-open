package jsonl

import (
	"context"
	"fmt"
	"time"

	"github.com/ebeyene/todo-open/internal/core"
)

type TaskRepo struct{}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (r *TaskRepo) Create(_ context.Context, _ core.Task) (core.Task, error) {
	return core.Task{}, fmt.Errorf("not implemented")
}

func (r *TaskRepo) GetByID(_ context.Context, _ string) (core.Task, error) {
	return core.Task{}, fmt.Errorf("not implemented")
}

func (r *TaskRepo) List(_ context.Context) ([]core.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *TaskRepo) Update(_ context.Context, _ core.Task) (core.Task, error) {
	return core.Task{}, fmt.Errorf("not implemented")
}

func (r *TaskRepo) Delete(_ context.Context, _ string, _ time.Time) error {
	return fmt.Errorf("not implemented")
}
