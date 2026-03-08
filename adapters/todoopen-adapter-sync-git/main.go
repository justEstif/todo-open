// todoopen-adapter-sync-git syncs a todo-open workspace to/from a git repository.
//
// Required config (in .todoopen/config.toml):
//
//	[adapters.git]
//	  bin = "todoopen-adapter-sync-git"
//
//	[adapters.git.config]
//	  remote = "origin"   # optional, default: origin
//	  branch = ""         # optional, default: current HEAD branch
//
// Implements the todoopen.plugin.v1 protocol over stdin/stdout.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// --- plugin protocol types (inlined; no shared internal dependency) ---

const protocolVersion = "todoopen.plugin.v1"

type handshakeResponse struct {
	ProtocolVersion string   `json:"protocol_version"`
	Name            string   `json:"name"`
	Kind            string   `json:"kind"`
	Capabilities    []string `json:"capabilities"`
	Health          health   `json:"health"`
}

type health struct {
	State   string `json:"state"`
	Message string `json:"message,omitempty"`
}

type requestEnvelope struct {
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Payload map[string]any `json:"payload,omitempty"`
}

type responseEnvelope struct {
	ID      string         `json:"id"`
	Payload map[string]any `json:"payload,omitempty"`
	Error   *pluginError   `json:"error,omitempty"`
}

type pluginError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// --- git adapter ---

type config struct {
	Remote string
	Branch string
}

func configFromPayload(payload map[string]any) config {
	cfg := config{Remote: "origin"}
	if raw, ok := payload["config"]; ok {
		if m, ok := raw.(map[string]any); ok {
			if v, ok := m["remote"].(string); ok && v != "" {
				cfg.Remote = v
			}
			if v, ok := m["branch"].(string); ok {
				cfg.Branch = v
			}
		}
	}
	return cfg
}

func workspaceRoot(payload map[string]any) (string, error) {
	v, ok := payload["workspace_root"].(string)
	if !ok || v == "" {
		return "", fmt.Errorf("workspace_root is required")
	}
	return v, nil
}

func gitCmd(dir string, args ...string) ([]byte, []byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func currentBranch(dir string) (string, error) {
	out, _, err := gitCmd(dir, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func handlePush(payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}
	cfg := configFromPayload(payload)

	// Stage existing workspace files.
	for _, f := range []string{"tasks.jsonl", "meta.json", "config.toml"} {
		_, stderr, err := gitCmd(root, "add", f)
		if err != nil && !strings.Contains(string(stderr), "did not match any files") {
			return nil, &pluginError{Code: "internal", Message: string(stderr)}
		}
	}

	// Commit.
	stdout, stderr, err := gitCmd(root, "commit", "-m", "chore: sync todo-open workspace [skip ci]")
	if err != nil {
		combined := string(stdout) + string(stderr)
		if strings.Contains(combined, "nothing to commit") {
			return map[string]any{"message": "nothing to commit"}, nil
		}
		return nil, &pluginError{Code: "internal", Message: strings.TrimSpace(string(stderr))}
	}

	// Resolve branch.
	branch := cfg.Branch
	if branch == "" {
		branch, err = currentBranch(root)
		if err != nil {
			return nil, &pluginError{Code: "internal", Message: err.Error()}
		}
	}

	// Push.
	_, stderr, err = gitCmd(root, "push", cfg.Remote, branch)
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: strings.TrimSpace(string(stderr))}
	}

	return map[string]any{"message": "pushed"}, nil
}

func handlePull(payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}
	cfg := configFromPayload(payload)

	branch := cfg.Branch
	if branch == "" {
		branch, err = currentBranch(root)
		if err != nil {
			return nil, &pluginError{Code: "internal", Message: err.Error()}
		}
	}

	_, stderr, err := gitCmd(root, "pull", "--ff-only", cfg.Remote, branch)
	if err != nil {
		msg := string(stderr)
		if strings.Contains(msg, "Not possible to fast-forward") {
			return nil, &pluginError{Code: "invalid_input", Message: "cannot fast-forward, manual merge required"}
		}
		return nil, &pluginError{Code: "internal", Message: strings.TrimSpace(msg)}
	}

	return map[string]any{"message": "pulled"}, nil
}

func revListCount(dir, revspec string) int {
	out, _, err := gitCmd(dir, "rev-list", revspec, "--count")
	if err != nil {
		return 0
	}
	var n int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &n)
	return n
}

func handleStatus(payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}

	out, stderr, err := gitCmd(root, "status", "--porcelain", "--",
		"tasks.jsonl", "meta.json", "config.toml")
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: strings.TrimSpace(string(stderr))}
	}

	clean := strings.TrimSpace(string(out)) == ""
	ahead := revListCount(root, "HEAD..@{u}")
	behind := revListCount(root, "@{u}..HEAD")

	return map[string]any{"clean": clean, "ahead": ahead, "behind": behind}, nil
}

// --- main ---

func writeJSON(v any) {
	b, _ := json.Marshal(v)
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}

func main() {
	if _, err := exec.LookPath("git"); err != nil {
		writeJSON(handshakeResponse{
			ProtocolVersion: protocolVersion,
			Name:            "git",
			Kind:            "sync",
			Capabilities:    []string{"pull", "push", "status"},
			Health:          health{State: "unhealthy", Message: "git not found in PATH"},
		})
		os.Exit(1)
	}

	writeJSON(handshakeResponse{
		ProtocolVersion: protocolVersion,
		Name:            "git",
		Kind:            "sync",
		Capabilities:    []string{"pull", "push", "status"},
		Health:          health{State: "ready"},
	})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req requestEnvelope
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}

		var payload map[string]any
		var perr *pluginError

		switch req.Method {
		case "push":
			payload, perr = handlePush(req.Payload)
		case "pull":
			payload, perr = handlePull(req.Payload)
		case "status":
			payload, perr = handleStatus(req.Payload)
		default:
			perr = &pluginError{Code: "not_supported", Message: "unknown method: " + req.Method}
		}

		writeJSON(responseEnvelope{ID: req.ID, Payload: payload, Error: perr})
	}
}
