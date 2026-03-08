package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWorkspaceMetaDefaultWhenMissing(t *testing.T) {
	t.Parallel()

	meta, err := LoadWorkspaceMeta(t.TempDir())
	if err != nil {
		t.Fatalf("load metadata: %v", err)
	}
	if meta.WorkspaceVersion != 1 {
		t.Fatalf("workspace_version = %d, want 1", meta.WorkspaceVersion)
	}
	if got, want := meta.SchemaVersion, "todo.open.task.v1"; got != want {
		t.Fatalf("schema_version = %q, want %q", got, want)
	}
}

func TestLoadWorkspaceMetaValidationFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "unknown field",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v1","unknown":1}`,
		},
		{
			name:    "version zero",
			payload: `{"workspace_version":0,"schema_version":"todo.open.task.v1"}`,
		},
		{
			name:    "bad schema version",
			payload: `{"workspace_version":1,"schema_version":"todo.open.task.v99"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			metaDir := filepath.Join(root, ".todoopen")
			if err := os.MkdirAll(metaDir, 0o755); err != nil {
				t.Fatalf("mkdir metadata dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(metaDir, "meta.json"), []byte(tc.payload), 0o644); err != nil {
				t.Fatalf("write metadata: %v", err)
			}

			if _, err := LoadWorkspaceMeta(root); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
