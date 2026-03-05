package jsonl

import (
	"context"
	"fmt"

	"github.com/ebeyene/todo-open/internal/core"
)

type TaskRepo struct{}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (r *TaskRepo) Create(_ context.Context, _ core.Task) (core.Task, error) {
	return core.Task{}, fmt.Errorf("not implemented")
}
