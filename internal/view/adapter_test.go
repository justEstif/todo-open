package view

import (
	"context"
	"errors"
	"testing"

	"github.com/justEstif/todo-open/internal/core"
)

type testAdapter struct{ name string }

func (a testAdapter) Name() string { return a.name }

func (a testAdapter) RenderTasks(_ context.Context, tasks []core.Task) ([]byte, error) {
	return nil, nil
}

func TestRegistryRegisterAndGet(t *testing.T) {
	r := NewRegistry()
	a := testAdapter{name: "table"}
	if err := r.Register(a); err != nil {
		t.Fatalf("register: %v", err)
	}

	got, err := r.Get("table")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name() != "table" {
		t.Fatalf("name = %q, want table", got.Name())
	}
}

func TestRegistryRejectsDuplicate(t *testing.T) {
	r := NewRegistry()
	a := testAdapter{name: "table"}
	if err := r.Register(a); err != nil {
		t.Fatalf("register: %v", err)
	}

	err := r.Register(a)
	if !errors.Is(err, ErrAdapterExists) {
		t.Fatalf("duplicate err = %v, want ErrAdapterExists", err)
	}
}

func TestRegistryNamesSorted(t *testing.T) {
	r := NewRegistry()
	_ = r.Register(testAdapter{name: "z"})
	_ = r.Register(testAdapter{name: "a"})

	names := r.Names()
	if len(names) != 2 || names[0] != "a" || names[1] != "z" {
		t.Fatalf("names = %v, want [a z]", names)
	}
}
