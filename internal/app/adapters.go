package app

import (
	"fmt"

	syncadapter "github.com/justEstif/todo-open/internal/sync"
	"github.com/justEstif/todo-open/internal/sync/noop"
	"github.com/justEstif/todo-open/internal/view"
	viewjson "github.com/justEstif/todo-open/internal/view/json"
)

// NewViewRegistry returns a registry pre-loaded with built-in view adapters.
func NewViewRegistry() (*view.Registry, error) {
	r := view.NewRegistry()
	if err := r.Register(viewjson.NewAdapter()); err != nil {
		return nil, fmt.Errorf("register json view adapter: %w", err)
	}
	return r, nil
}

// NewSyncRegistry returns a registry pre-loaded with built-in sync adapters.
func NewSyncRegistry() (*syncadapter.Registry, error) {
	r := syncadapter.NewRegistry()
	if err := r.Register(noop.NewAdapter()); err != nil {
		return nil, fmt.Errorf("register noop sync adapter: %w", err)
	}
	return r, nil
}
