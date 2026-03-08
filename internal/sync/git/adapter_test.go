package git

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitAdapter_Push(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not in PATH")
	}

	t.Run("push with changes", func(t *testing.T) {
		dir := t.TempDir()
		adapter := setupGitRepo(t, dir)

		// Modify tasks.jsonl
		tasksFile := filepath.Join(dir, "tasks.jsonl")
		if err := os.WriteFile(tasksFile, []byte(`{"id":"test","title":"test task"}`+"\n"), 0644); err != nil {
			t.Fatalf("failed to write tasks.jsonl: %v", err)
		}

		result, err := adapter.Push()
		if err != nil {
			t.Fatalf("Push() failed: %v", err)
		}

		if msg, ok := result["message"].(string); !ok || msg != "pushed" {
			t.Errorf("Push() returned wrong message: got %v, want 'pushed'", result["message"])
		}

		// Verify the commit was made
		output, err := exec.Command("git", "-C", dir, "log", "--oneline").Output()
		if err != nil {
			t.Fatalf("failed to get git log: %v", err)
		}
		if !strings.Contains(string(output), "chore: sync todo-open workspace") {
			t.Errorf("expected commit message not found in git log: %s", string(output))
		}
	})

	t.Run("push idempotent - no changes", func(t *testing.T) {
		dir := t.TempDir()
		adapter := setupGitRepo(t, dir)

		// First push should succeed
		_, err := adapter.Push()
		if err != nil {
			t.Fatalf("first Push() failed: %v", err)
		}

		// Second push should return "nothing to commit"
		result, err := adapter.Push()
		if err != nil {
			t.Fatalf("second Push() failed: %v", err)
		}

		if msg, ok := result["message"].(string); !ok || msg != "nothing to commit" {
			t.Errorf("second Push() returned wrong message: got %v, want 'nothing to commit'", result["message"])
		}
	})
}

func TestGitAdapter_Status(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not in PATH")
	}

	t.Run("status clean", func(t *testing.T) {
		dir := t.TempDir()
		adapter := setupGitRepo(t, dir)

		// Make sure all files are committed
		_, err := adapter.Push()
		if err != nil {
			t.Fatalf("Push() failed: %v", err)
		}

		result, err := adapter.Status()
		if err != nil {
			t.Fatalf("Status() failed: %v", err)
		}

		if clean, ok := result["clean"].(bool); !ok || !clean {
			t.Errorf("Status() returned clean=%v, want true", result["clean"])
		}
	})

	t.Run("status dirty", func(t *testing.T) {
		dir := t.TempDir()
		adapter := setupGitRepo(t, dir)

		// Modify tasks.jsonl without committing
		tasksFile := filepath.Join(dir, "tasks.jsonl")
		if err := os.WriteFile(tasksFile, []byte(`{"id":"modified","title":"modified task"}`+"\n"), 0644); err != nil {
			t.Fatalf("failed to write tasks.jsonl: %v", err)
		}

		result, err := adapter.Status()
		if err != nil {
			t.Fatalf("Status() failed: %v", err)
		}

		if clean, ok := result["clean"].(bool); !ok || clean {
			t.Errorf("Status() returned clean=%v, want false", result["clean"])
		}
	})
}

func TestGitAdapter_Pull(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not in PATH")
	}

	t.Run("pull success", func(t *testing.T) {
		dir := t.TempDir()
		adapter := setupGitRepo(t, dir)

		result, err := adapter.Pull()
		if err != nil {
			t.Fatalf("Pull() failed: %v", err)
		}

		if msg, ok := result["message"].(string); !ok || msg != "pulled" {
			t.Errorf("Pull() returned wrong message: got %v, want 'pulled'", result["message"])
		}
	})
}

func TestGitAdapter_ConfigJSON(t *testing.T) {
	config := Config{
		Remote: "origin",
		Branch: "main",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if decoded.Remote != config.Remote || decoded.Branch != config.Branch {
		t.Errorf("config roundtrip failed: got %+v, want %+v", decoded, config)
	}
}

func TestGitAdapter_StatusResultJSON(t *testing.T) {
	result := StatusResult{
		Clean:  true,
		Ahead:  0,
		Behind: 0,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal StatusResult: %v", err)
	}

	var decoded StatusResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal StatusResult: %v", err)
	}

	if decoded.Clean != result.Clean || decoded.Ahead != result.Ahead || decoded.Behind != result.Behind {
		t.Errorf("StatusResult roundtrip failed: got %+v, want %+v", decoded, result)
	}
}

// setupGitRepo initializes a git repository in the given directory with
// initial files and returns a GitAdapter configured for that directory.
func setupGitRepo(t *testing.T, dir string) *GitAdapter {
	t.Helper()

	// Initialize git repo
	if err := exec.Command("git", "init", dir).Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user (required for commits)
	if err := exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run(); err != nil {
		t.Fatalf("failed to set git user email: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run(); err != nil {
		t.Fatalf("failed to set git user name: %v", err)
	}

	// Create initial files
	tasksFile := filepath.Join(dir, "tasks.jsonl")
	if err := os.WriteFile(tasksFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write tasks.jsonl: %v", err)
	}

	metaFile := filepath.Join(dir, "meta.json")
	if err := os.WriteFile(metaFile, []byte(`{"version":"1.0.0"}`), 0644); err != nil {
		t.Fatalf("failed to write meta.json: %v", err)
	}

	configFile := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(configFile, []byte("[server]\naddr = \":8080\""), 0644); err != nil {
		t.Fatalf("failed to write config.toml: %v", err)
	}

	// Add and commit initial files
	if err := exec.Command("git", "-C", dir, "add", ".").Run(); err != nil {
		t.Fatalf("failed to add files: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "commit", "-m", "initial commit").Run(); err != nil {
		t.Fatalf("failed to commit initial files: %v", err)
	}

	// Create a remote (for testing purposes)
	remoteDir := t.TempDir()
	if err := exec.Command("git", "init", "--bare", remoteDir).Run(); err != nil {
		t.Fatalf("failed to create bare remote: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "remote", "add", "origin", remoteDir).Run(); err != nil {
		t.Fatalf("failed to add remote: %v", err)
	}
	if err := exec.Command("git", "-C", dir, "push", "-u", "origin", "master").Run(); err != nil {
		// Try with main branch if master fails
		if err := exec.Command("git", "-C", dir, "push", "-u", "origin", "main").Run(); err != nil {
			t.Fatalf("failed to push to remote: %v", err)
		}
	}

	return NewGitAdapter(dir, Config{
		Remote: "origin",
		Branch: "",
	})
}
