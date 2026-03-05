package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/api/handlers"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/jsonl"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func NewServer(addr string) (*http.Server, error) {
	workspaceRoot, err := resolveWorkspaceRoot()
	if err != nil {
		return nil, err
	}

	configPath := os.Getenv("TODOOPEN_ADAPTER_CONFIG")
	if strings.TrimSpace(configPath) == "" {
		configPath = filepath.Join(workspaceRoot, ".todoopen", "adapters.json")
	}

	cfg, err := LoadAdapterConfig(configPath)
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
	runtime := BuildAdapterRuntime(cfg, viewRegistry, syncRegistry)
	if !runtime.Ready {
		return nil, fmt.Errorf("adapter initialization failed: %s", strings.Join(runtime.Errors, "; "))
	}

	repo := defaultTaskRepo(workspaceRoot)
	taskService := core.NewService(repo, time.Now, nil)

	return &http.Server{
		Addr: addr,
		Handler: api.NewRouter(taskService, handlers.AdapterRuntimeResponse{
			Config: handlers.AdapterConfigResponse{
				EnabledViews:        runtime.Config.EnabledViews,
				EnabledSyncAdapters: runtime.Config.EnabledSyncAdapters,
				ViewSettings:        runtime.Config.ViewSettings,
				SyncSettings:        runtime.Config.SyncSettings,
			},
			Status: toAdapterStatusResponse(runtime.Status),
			Ready:  runtime.Ready,
			Errors: runtime.Errors,
		}),
	}, nil
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

func toAdapterStatusResponse(statuses []AdapterStatus) []handlers.AdapterStatusResponse {
	out := make([]handlers.AdapterStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		out = append(out, handlers.AdapterStatusResponse{
			Kind:    status.Kind,
			Name:    status.Name,
			Enabled: status.Enabled,
			Healthy: status.Healthy,
			Message: status.Message,
		})
	}
	return out
}
