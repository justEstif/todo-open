package store

import "github.com/justEstif/todo-open/internal/core"

// TaskRepository re-exports the canonical domain-owned repository contract.
// Store implementations should satisfy core.TaskRepository.
type TaskRepository = core.TaskRepository
