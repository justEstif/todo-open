package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/justEstif/todo-open/internal/plugin"
)

const (
	defaultSchemaVersion = "todo.open.task.v1"
	metaDirName          = ".todoopen"
	metaFileName         = "meta.json"
)

type WorkspaceMeta struct {
	WorkspaceVersion    int                   `json:"workspace_version"`
	SchemaVersion       string                `json:"schema_version"`
	DefaultSort         []string              `json:"default_sort,omitempty"`
	EnabledViews        []string              `json:"enabled_views,omitempty"`
	EnabledSyncAdapters []string              `json:"enabled_sync_adapters,omitempty"`
	AdapterPlugins      []AdapterPluginConfig `json:"adapter_plugins,omitempty"`
	Ext                 map[string]any        `json:"ext,omitempty"`
}

type AdapterPluginConfig struct {
	Name    string             `json:"name"`
	Kind    plugin.AdapterKind `json:"kind"`
	Command string             `json:"command"`
	Args    []string           `json:"args,omitempty"`
}

func DefaultWorkspaceMeta() WorkspaceMeta {
	return WorkspaceMeta{
		WorkspaceVersion:    1,
		SchemaVersion:       defaultSchemaVersion,
		EnabledViews:        []string{"json"},
		EnabledSyncAdapters: []string{"noop"},
	}
}

func LoadWorkspaceMeta(workspaceRoot string) (WorkspaceMeta, error) {
	metaPath := filepath.Join(workspaceRoot, metaDirName, metaFileName)

	data, err := os.ReadFile(metaPath)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultWorkspaceMeta(), nil
	}
	if err != nil {
		return WorkspaceMeta{}, fmt.Errorf("read workspace metadata: %w", err)
	}

	meta := DefaultWorkspaceMeta()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&meta); err != nil {
		return WorkspaceMeta{}, fmt.Errorf("decode workspace metadata: %w", err)
	}

	if err := validateWorkspaceMeta(meta); err != nil {
		return WorkspaceMeta{}, err
	}
	return meta, nil
}

func validateWorkspaceMeta(meta WorkspaceMeta) error {
	if meta.WorkspaceVersion < 1 {
		return fmt.Errorf("workspace_version must be >= 1")
	}
	if strings.TrimSpace(meta.SchemaVersion) != defaultSchemaVersion {
		return fmt.Errorf("unsupported schema_version: %s", meta.SchemaVersion)
	}
	if err := validateNames("enabled_views", meta.EnabledViews); err != nil {
		return err
	}
	if err := validateNames("enabled_sync_adapters", meta.EnabledSyncAdapters); err != nil {
		return err
	}

	pluginSeen := map[string]struct{}{}
	for i, p := range meta.AdapterPlugins {
		if strings.TrimSpace(p.Name) == "" {
			return fmt.Errorf("adapter_plugins[%d].name is required", i)
		}
		if p.Kind != plugin.AdapterKindView && p.Kind != plugin.AdapterKindSync {
			return fmt.Errorf("adapter_plugins[%d].kind must be one of: view, sync", i)
		}
		if strings.TrimSpace(p.Command) == "" {
			return fmt.Errorf("adapter_plugins[%d].command is required", i)
		}

		key := string(p.Kind) + ":" + p.Name
		if _, ok := pluginSeen[key]; ok {
			return fmt.Errorf("adapter_plugins contains duplicate adapter %q for kind %q", p.Name, p.Kind)
		}
		pluginSeen[key] = struct{}{}
	}

	return nil
}
