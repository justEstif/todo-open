package view

import (
	"context"
	"fmt"

	"github.com/justEstif/todo-open/internal/adapterregistry"
	"github.com/justEstif/todo-open/internal/core"
)

var (
	ErrAdapterNameRequired = adapterregistry.ErrAdapterNameRequired
	ErrAdapterExists       = adapterregistry.ErrAdapterExists
	ErrAdapterNotFound     = adapterregistry.ErrAdapterNotFound
)

// Adapter renders tasks for a specific view target.
type Adapter interface {
	Name() string
	RenderTasks(ctx context.Context, tasks []core.Task) ([]byte, error)
}

// Registry stores runtime view adapters by name.
type Registry struct {
	reg *adapterregistry.Registry[Adapter]
}

func NewRegistry() *Registry {
	return &Registry{reg: adapterregistry.New[Adapter]()}
}

func (r *Registry) Register(adapter Adapter) error {
	if adapter == nil {
		return fmt.Errorf("nil adapter: %w", ErrAdapterNameRequired)
	}
	return r.reg.Register(adapter)
}

func (r *Registry) Get(name string) (Adapter, error) {
	return r.reg.Get(name)
}

func (r *Registry) Names() []string {
	return r.reg.Names()
}
