package plugin

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrPluginCommandRequired = errors.New("plugin command is required")
	ErrPluginStartTimeout    = errors.New("plugin start timeout")
)

// Definition describes a plugin process to launch.
type Definition struct {
	Name    string
	Kind    AdapterKind
	Command string
	Args    []string
}

type Loader struct {
	handshakeTimeout time.Duration
}

func NewLoader(handshakeTimeout time.Duration) *Loader {
	if handshakeTimeout <= 0 {
		handshakeTimeout = 2 * time.Second
	}
	return &Loader{handshakeTimeout: handshakeTimeout}
}

// LoadedPlugin represents a running plugin process.
type LoadedPlugin struct {
	def       Definition
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	handshake HandshakeResponse

	exited bool

	mu sync.RWMutex
}

func (p *LoadedPlugin) Name() string {
	return p.def.Name
}

func (p *LoadedPlugin) Kind() AdapterKind {
	return p.def.Kind
}

func (p *LoadedPlugin) Handshake() HandshakeResponse {
	return p.handshake
}

func (p *LoadedPlugin) Health() Health {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.exited {
		return Health{State: HealthUnhealthy, Message: "plugin process exited"}
	}
	return p.handshake.Health
}

func (p *LoadedPlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd == nil || p.cmd.Process == nil {
		p.exited = true
		return nil
	}

	_ = p.stdin.Close()
	_ = p.stdout.Close()
	_ = p.stderr.Close()

	if p.exited {
		return nil
	}
	if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("kill plugin process: %w", err)
	}
	p.exited = true
	return nil
}

func (l *Loader) Load(ctx context.Context, def Definition) (*LoadedPlugin, error) {
	if strings.TrimSpace(def.Command) == "" {
		return nil, ErrPluginCommandRequired
	}

	cmdPath, err := discoverCommand(def.Command)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, cmdPath, def.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("plugin stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("plugin stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("plugin stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start plugin process: %w", err)
	}

	resp, err := l.readHandshake(ctx, stdout, def)
	if err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		_ = stderr.Close()
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		return nil, err
	}

	pl := &LoadedPlugin{
		def:       def,
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		handshake: resp,
	}
	go func() {
		_ = cmd.Wait()
		pl.mu.Lock()
		pl.exited = true
		pl.mu.Unlock()
	}()

	return pl, nil
}

func (l *Loader) readHandshake(parent context.Context, stdout io.Reader, def Definition) (HandshakeResponse, error) {
	ctx, cancel := context.WithTimeout(parent, l.handshakeTimeout)
	defer cancel()

	type result struct {
		resp HandshakeResponse
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		var resp HandshakeResponse
		line, err := bufio.NewReader(stdout).ReadBytes('\n')
		if err != nil {
			ch <- result{err: fmt.Errorf("read plugin handshake: %w", err)}
			return
		}
		if err := json.Unmarshal(line, &resp); err != nil {
			ch <- result{err: fmt.Errorf("decode plugin handshake: %w", err)}
			return
		}
		ch <- result{resp: resp}
	}()

	select {
	case <-ctx.Done():
		return HandshakeResponse{}, fmt.Errorf("plugin %q: %w", def.Name, ErrPluginStartTimeout)
	case res := <-ch:
		if res.err != nil {
			return HandshakeResponse{}, res.err
		}
		req := HandshakeRequest{
			ProtocolVersion: ProtocolVersion,
			ExpectedName:    def.Name,
			ExpectedKind:    def.Kind,
		}
		if err := ValidateHandshake(req, res.resp); err != nil {
			return HandshakeResponse{}, fmt.Errorf("validate plugin handshake: %w", err)
		}
		return res.resp, nil
	}
}

func discoverCommand(command string) (string, error) {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return "", ErrPluginCommandRequired
	}
	if strings.ContainsRune(trimmed, filepath.Separator) {
		return trimmed, nil
	}
	path, err := exec.LookPath(trimmed)
	if err != nil {
		return "", fmt.Errorf("discover plugin command %q: %w", trimmed, err)
	}
	return path, nil
}
