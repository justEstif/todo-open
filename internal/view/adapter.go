package view

import (
	"context"
	"errors"
	"fmt"
	"slices"
	stdsync "sync"

	"github.com/justEstif/todo-open/internal/core"
)

var (
	ErrAdapterNameRequired = errors.New("adapter name is required")
	ErrAdapterExists       = errors.New("adapter already registered")
	ErrAdapterNotFound     = errors.New("adapter not found")
)

// Adapter renders tasks for a specific view target.
type Adapter interface {
	Name() string
	RenderTasks(ctx context.Context, tasks []core.Task) ([]byte, error)
}

// Registry stores runtime view adapters by name.
type Registry struct {
	mu       stdsync.RWMutex
	adapters map[string]Adapter
}

func NewRegistry() *Registry {
	return &Registry{adapters: make(map[string]Adapter)}
}

func (r *Registry) Register(adapter Adapter) error {
	if adapter == nil {
		return fmt.Errorf("nil adapter: %w", ErrAdapterNameRequired)
	}

	name := adapter.Name()
	if name == "" {
		return ErrAdapterNameRequired
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adapters[name]; exists {
		return fmt.Errorf("%s: %w", name, ErrAdapterExists)
	}

	r.adapters[name] = adapter
	return nil
}

func (r *Registry) Get(name string) (Adapter, error) {
	if name == "" {
		return nil, ErrAdapterNameRequired
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[name]
	if !ok {
		return nil, fmt.Errorf("%s: %w", name, ErrAdapterNotFound)
	}
	return adapter, nil
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
