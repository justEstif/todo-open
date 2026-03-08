package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
	"github.com/justEstif/todo-open/internal/store/jsonl"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func NewServer(addr string) (*http.Server, error) {
	workspaceRoot, err := resolveWorkspaceRoot()
	if err != nil {
		return nil, err
	}

	_, err = LoadWorkspaceMeta(workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load workspace metadata: %w", err)
	}

	adapterCfg, err := LoadAdapterFileConfig(workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load adapter config: %w", err)
	}

	viewRegistry, err := NewViewRegistry()
	if err != nil {
		return nil, fmt.Errorf("load view adapters: %w", err)
	}
	syncRegistry, err := NewSyncRegistry()
	if err != nil {
		return nil, fmt.Errorf("load sync adapters: %w", err)
	}
	runtime := BuildAdapterRuntimeFromConfig(context.Background(), adapterCfg, viewRegistry, syncRegistry)
	if !runtime.Ready {
		return nil, fmt.Errorf("adapter initialization failed: %s", strings.Join(runtime.Errors, "; "))
	}

	repo := defaultTaskRepo(workspaceRoot)
	broker := events.NewBroker()
	taskService := core.NewService(repo, time.Now, nil)
	taskService.OnMutation(func(e core.MutationEvent) {
		var event events.Event
		event.Type = e.Type
		event.Task = e.Task
		event.OldStatus = e.OldStatus
		event.NewStatus = e.NewStatus
		event.At = e.At
		broker.Publish(event)
	})

	// Start the lease sweeper background goroutine. It stops when the server context is cancelled.
	// We use context.Background() here; the sweeper will be stopped via server shutdown in main.
	sweeperCtx, sweeperCancel := context.WithCancel(context.Background())

	srv := &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(taskService, runtime, broker),
	}

	StartLeaseSweeper(sweeperCtx, taskService, 30*time.Second)

	// Wrap srv.Shutdown to also cancel the sweeper.
	origShutdown := srv.Shutdown
	_ = origShutdown // keep linter happy; we set a custom RegisterOnShutdown callback instead
	srv.RegisterOnShutdown(sweeperCancel)

	return srv, nil
}

func defaultTaskRepo(workspaceRoot string) core.TaskRepository {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("TODOOPEN_STORE")), "memory") {
		return memory.NewTaskRepo()
	}
	return jsonl.NewTaskRepo(workspaceRoot)
}

func resolveWorkspaceRoot() (string, error) {
	if root := strings.TrimSpace(os.Getenv("TODOOPEN_WORKSPACE_ROOT")); root != "" {
		return root, nil
	}
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	return root, nil
}
