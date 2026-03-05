package noop

import (
	"context"
	"testing"

	"github.com/justEstif/todo-open/internal/core"
)

func TestNoopAdapter(t *testing.T) {
	a := NewAdapter()
	if a.Name() != "noop" {
		t.Fatalf("name = %q, want noop", a.Name())
	}

	if err := a.Push(context.Background(), []core.Task{{ID: "t1"}}); err != nil {
		t.Fatalf("push: %v", err)
	}

	got, err := a.Pull(context.Background())
	if err != nil {
		t.Fatalf("pull: %v", err)
	}
	if got != nil {
		t.Fatalf("pull = %#v, want nil", got)
	}
}
