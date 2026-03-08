package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultSchemaVersion = "todo.open.task.v1"
	metaDirName          = ".todoopen"
	metaFileName         = "meta.json"
)

// WorkspaceMeta holds workspace-level identity and schema metadata.
// Adapter configuration lives separately in .todoopen/config.toml.
type WorkspaceMeta struct {
	WorkspaceVersion int            `json:"workspace_version"`
	SchemaVersion    string         `json:"schema_version"`
	DefaultSort      []string       `json:"default_sort,omitempty"`
	Ext              map[string]any `json:"ext,omitempty"`
}

func DefaultWorkspaceMeta() WorkspaceMeta {
	return WorkspaceMeta{
		WorkspaceVersion: 1,
		SchemaVersion:    defaultSchemaVersion,
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
	return nil
}
