package adapterregistry

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

var (
	ErrAdapterNameRequired = errors.New("adapter name is required")
	ErrAdapterExists       = errors.New("adapter already registered")
	ErrAdapterNotFound     = errors.New("adapter not found")
)

type Named interface {
	Name() string
}

type Registry[T Named] struct {
	mu       sync.RWMutex
	adapters map[string]T
}

func New[T Named]() *Registry[T] {
	return &Registry[T]{adapters: make(map[string]T)}
}

func (r *Registry[T]) Register(adapter T) error {
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

func (r *Registry[T]) Get(name string) (T, error) {
	var zero T
	if name == "" {
		return zero, ErrAdapterNameRequired
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[name]
	if !ok {
		return zero, fmt.Errorf("%s: %w", name, ErrAdapterNotFound)
	}
	return adapter, nil
}

func (r *Registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
