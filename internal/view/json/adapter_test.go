package json

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/justEstif/todo-open/internal/core"
)

func TestAdapterRenderTasks(t *testing.T) {
	a := NewAdapter()
	if a.Name() != "json" {
		t.Fatalf("name = %q, want json", a.Name())
	}

	out, err := a.RenderTasks(context.Background(), []core.Task{{ID: "t1", Title: "one"}})
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	var got struct {
		Items []core.Task `json:"items"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].ID != "t1" {
		t.Fatalf("items = %#v", got.Items)
	}
}
