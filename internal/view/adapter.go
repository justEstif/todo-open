package view

import (
	"context"

	"github.com/justEstif/todo-open/internal/adapterregistry"
	"github.com/justEstif/todo-open/internal/core"
)

// Adapter renders tasks for a specific view target.
type Adapter interface {
	Name() string
	RenderTasks(ctx context.Context, tasks []core.Task) ([]byte, error)
}

// Re-export sentinel errors for package consumers and tests.
var (
	ErrAdapterNameRequired = adapterregistry.ErrAdapterNameRequired
	ErrAdapterExists       = adapterregistry.ErrAdapterExists
	ErrAdapterNotFound     = adapterregistry.ErrAdapterNotFound
)

// Registry is a named-adapter registry for view adapters.
type Registry = adapterregistry.Registry[Adapter]

// NewRegistry returns a ready-to-use view adapter registry.
func NewRegistry() *Registry {
	return adapterregistry.New[Adapter]()
}
