package noop

import (
	"context"

	"github.com/justEstif/todo-open/internal/core"
)

type Adapter struct{}

func NewAdapter() Adapter {
	return Adapter{}
}

func (Adapter) Name() string {
	return "noop"
}

func (Adapter) Push(_ context.Context, _ []core.Task) error {
	return nil
}

func (Adapter) Pull(_ context.Context) ([]core.Task, error) {
	return nil, nil
}
