package json

import (
	"context"
	"encoding/json"

	"github.com/justEstif/todo-open/internal/core"
)

type Adapter struct{}

func NewAdapter() Adapter {
	return Adapter{}
}

func (Adapter) Name() string {
	return "json"
}

func (Adapter) RenderTasks(_ context.Context, tasks []core.Task) ([]byte, error) {
	resp := struct {
		Items []core.Task `json:"items"`
	}{Items: tasks}

	return json.Marshal(resp)
}
