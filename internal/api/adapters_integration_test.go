package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func TestAdaptersStatusEndpoint(t *testing.T) {
	repo := memory.NewTaskRepo()
	svc := core.NewService(repo, time.Now, func() string { return "task_1" })
	runtime := adapters.Runtime{
		Config: adapters.Config{
			EnabledViews:        []string{"json"},
			EnabledSyncAdapters: []string{"noop"},
		},
		Status: []adapters.Status{
			{Kind: "sync", Name: "noop", Enabled: true, Healthy: true},
			{Kind: "view", Name: "json", Enabled: true, Healthy: true},
		},
		Ready: true,
	}
	server := httptest.NewServer(api.NewRouter(svc, runtime))
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/v1/adapters")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var out adapters.Runtime
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if !out.Ready {
		t.Fatal("ready = false, want true")
	}
	if len(out.Status) != 2 {
		t.Fatalf("status entries = %d, want 2", len(out.Status))
	}
}
