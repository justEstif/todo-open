package memory

import (
	"context"
	"sync"

	"github.com/justEstif/todo-open/internal/core"
)

type TaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]core.Task
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{tasks: map[string]core.Task{}}
}

func (r *TaskRepo) Create(_ context.Context, task core.Task) (core.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		return core.Task{}, core.ErrInvalidInput
	}
	r.tasks[task.ID] = task
	return task, nil
}

func (r *TaskRepo) GetByID(_ context.Context, id string) (core.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok || task.DeletedAt != nil {
		return core.Task{}, core.ErrNotFound
	}
	return task, nil
}

func (r *TaskRepo) List(_ context.Context) ([]core.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]core.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		if task.DeletedAt == nil {
			out = append(out, task)
		}
	}
	return out, nil
}

func (r *TaskRepo) Update(_ context.Context, task core.Task) (core.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.tasks[task.ID]
	if !ok || existing.DeletedAt != nil {
		return core.Task{}, core.ErrNotFound
	}
	r.tasks[task.ID] = task
	return task, nil
}
