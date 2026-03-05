package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func NewServer(addr string) (*http.Server, error) {
	configPath := os.Getenv("TODOOPEN_ADAPTER_CONFIG")
	if strings.TrimSpace(configPath) == "" {
		configPath = ".todoopen/adapters.json"
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

	repo := memory.NewTaskRepo()
	taskService := core.NewService(repo, time.Now, nil)

	return &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(taskService),
	}, nil
}
