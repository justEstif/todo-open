package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Config struct {
	Remote string `json:"remote"`
	Branch string `json:"branch"`
}

type StatusResult struct {
	Clean  bool `json:"clean"`
	Ahead  int  `json:"ahead"`
	Behind int  `json:"behind"`
}

type GitAdapter struct {
	WorkspaceRoot string
	Config        Config
}

func NewGitAdapter(workspaceRoot string, config Config) *GitAdapter {
	return &GitAdapter{
		WorkspaceRoot: workspaceRoot,
		Config:        config,
	}
}

func (g *GitAdapter) Push() (map[string]any, error) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, &PluginError{
			Code:    "unavailable",
			Message: "git not found",
		}
	}

	// Change to workspace directory
	if err := g.ensureInWorkspace(); err != nil {
		return nil, err
	}

	// Add the relevant files
	files := []string{"tasks.jsonl", "meta.json", "config.toml"}
	for _, file := range files {
		if err := g.runGitCmd("add", file); err != nil {
			// It's okay if the file doesn't exist
			if !strings.Contains(err.Error(), "did not match any files") {
				return nil, err
			}
		}
	}

	// Try to commit
	commitMsg := "chore: sync todo-open workspace [skip ci]"
	output, err := g.runGitCmdOutput("commit", "-m", commitMsg)
	if err != nil {
		if strings.Contains(string(output), "nothing to commit") {
			return map[string]any{"message": "nothing to commit"}, nil
		}
		return nil, err
	}

	// Determine remote and branch
	remote := g.Config.Remote
	if remote == "" {
		remote = "origin"
	}

	branch := g.Config.Branch
	if branch == "" {
		branchOutput, err := g.runGitCmdOutput("symbolic-ref", "--short", "HEAD")
		if err != nil {
			return nil, fmt.Errorf("failed to get current branch: %w", err)
		}
		branch = strings.TrimSpace(string(branchOutput))
	}

	// Push to remote
	if err := g.runGitCmd("push", remote, branch); err != nil {
		return nil, err
	}

	return map[string]any{"message": "pushed"}, nil
}

func (g *GitAdapter) Pull() (map[string]any, error) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, &PluginError{
			Code:    "unavailable",
			Message: "git not found",
		}
	}

	// Change to workspace directory
	if err := g.ensureInWorkspace(); err != nil {
		return nil, err
	}

	// Determine remote and branch
	remote := g.Config.Remote
	if remote == "" {
		remote = "origin"
	}

	branch := g.Config.Branch
	if branch == "" {
		branchOutput, err := g.runGitCmdOutput("symbolic-ref", "--short", "HEAD")
		if err != nil {
			return nil, fmt.Errorf("failed to get current branch: %w", err)
		}
		branch = strings.TrimSpace(string(branchOutput))
	}

	// Try to pull with fast-forward only
	output, err := g.runGitCmdOutput("pull", "--ff-only", remote, branch)
	if err != nil {
		if strings.Contains(string(output), "Not possible to fast-forward") {
			return nil, &PluginError{
				Code:    "invalid_input",
				Message: "cannot fast-forward, manual merge required",
			}
		}
		return nil, err
	}

	return map[string]any{"message": "pulled"}, nil
}

func (g *GitAdapter) Status() (map[string]any, error) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, &PluginError{
			Code:    "unavailable",
			Message: "git not found",
		}
	}

	// Change to workspace directory
	if err := g.ensureInWorkspace(); err != nil {
		return nil, err
	}

	result := StatusResult{}

	// Check if working directory is clean
	result.Clean = true
	files := []string{"tasks.jsonl", "meta.json", "config.toml"}
	for _, file := range files {
		output, err := g.runGitCmdOutput("status", "--porcelain", "--", file)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(string(output)) != "" {
			result.Clean = false
			break
		}
	}

	// If our specific files are clean, do a general status check
	if result.Clean {
		output, err := g.runGitCmdOutput("status", "--porcelain", "--")
		if err != nil {
			return nil, err
		}
		result.Clean = strings.TrimSpace(string(output)) == ""
	}

	// Get ahead/behind counts
	aheadCount, err := g.getRevListCount("HEAD..@{u}")
	if err == nil {
		result.Ahead = aheadCount
	}

	behindCount, err := g.getRevListCount("@{u}..HEAD")
	if err == nil {
		result.Behind = behindCount
	}

	return map[string]any{
		"clean":  result.Clean,
		"ahead":  result.Ahead,
		"behind": result.Behind,
	}, nil
}

func (g *GitAdapter) ensureInWorkspace() error {
	if g.WorkspaceRoot == "" {
		return fmt.Errorf("workspace root is required")
	}
	return nil
}

func (g *GitAdapter) runGitCmd(args ...string) error {
	cmd := exec.Command("git", args...)
	if g.WorkspaceRoot != "" {
		cmd.Dir = g.WorkspaceRoot
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git %v: %s", args, stderr.String())
	}
	return nil
}

func (g *GitAdapter) runGitCmdOutput(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	if g.WorkspaceRoot != "" {
		cmd.Dir = g.WorkspaceRoot
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return output, fmt.Errorf("git %v: %s", args, stderr.String())
	}
	return output, nil
}

func (g *GitAdapter) getRevListCount(revspec string) (int, error) {
	output, err := g.runGitCmdOutput("rev-list", revspec, "--count")
	if err != nil {
		return 0, nil // Gracefully handle errors (e.g., no upstream)
	}

	var count int
	if _, err := fmt.Sscanf(string(output), "%d", &count); err != nil {
		return 0, nil // Gracefully handle parsing errors
	}

	return count, nil
}

// PluginError implements the plugin.PluginError interface
type PluginError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *PluginError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
