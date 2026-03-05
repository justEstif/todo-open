package app

import (
	"fmt"

	syncadapter "github.com/justEstif/todo-open/internal/sync"
	"github.com/justEstif/todo-open/internal/sync/noop"
	"github.com/justEstif/todo-open/internal/view"
	viewjson "github.com/justEstif/todo-open/internal/view/json"
)

// NewViewRegistry loads built-in view adapters.
func NewViewRegistry() (*view.Registry, error) {
	registry := view.NewRegistry()
	if err := registry.Register(viewjson.NewAdapter()); err != nil {
		return nil, fmt.Errorf("register json view adapter: %w", err)
	}
	return registry, nil
}

// NewSyncRegistry loads built-in sync adapters.
func NewSyncRegistry() (*syncadapter.Registry, error) {
	registry := syncadapter.NewRegistry()
	if err := registry.Register(noop.NewAdapter()); err != nil {
		return nil, fmt.Errorf("register noop sync adapter: %w", err)
	}
	return registry, nil
}
