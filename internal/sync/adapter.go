package sync

import (
	"context"

	"github.com/justEstif/todo-open/internal/adapterregistry"
	"github.com/justEstif/todo-open/internal/core"
)

// Adapter pushes and pulls task changes from an external system.
type Adapter interface {
	Name() string
	Push(ctx context.Context, tasks []core.Task) error
	Pull(ctx context.Context) ([]core.Task, error)
}

// Re-export sentinel errors for package consumers and tests.
var (
	ErrAdapterNameRequired = adapterregistry.ErrAdapterNameRequired
	ErrAdapterExists       = adapterregistry.ErrAdapterExists
	ErrAdapterNotFound     = adapterregistry.ErrAdapterNotFound
)

// Registry is a named-adapter registry for sync adapters.
type Registry = adapterregistry.Registry[Adapter]

// NewRegistry returns a ready-to-use sync adapter registry.
func NewRegistry() *Registry {
	return adapterregistry.New[Adapter]()
}
