package plugin

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLoaderLoadSuccess(t *testing.T) {
	t.Parallel()

	loader := NewLoader(2 * time.Second)
	def := Definition{
		Name:    "markdown",
		Kind:    AdapterKindView,
		Command: "sh",
		Args:    []string{"-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"markdown\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 5"},
	}

	p, err := loader.Load(context.Background(), def)
	if err != nil {
		t.Fatalf("load plugin: %v", err)
	}
	t.Cleanup(func() {
		_ = p.Close()
	})

	if p.Name() != "markdown" {
		t.Fatalf("name = %q", p.Name())
	}
	if p.Kind() != AdapterKindView {
		t.Fatalf("kind = %q", p.Kind())
	}
	if got := p.Health().State; got != HealthReady {
		t.Fatalf("health = %q, want ready", got)
	}
}

func TestLoaderLoadMissingCommand(t *testing.T) {
	t.Parallel()

	loader := NewLoader(time.Second)
	_, err := loader.Load(context.Background(), Definition{Name: "x", Kind: AdapterKindView})
	if !errors.Is(err, ErrPluginCommandRequired) {
		t.Fatalf("err = %v, want ErrPluginCommandRequired", err)
	}
}

func TestLoaderLoadHandshakeValidationFailure(t *testing.T) {
	t.Parallel()

	loader := NewLoader(time.Second)
	def := Definition{
		Name:    "markdown",
		Kind:    AdapterKindView,
		Command: "sh",
		Args:    []string{"-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"wrong\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; sleep 1"},
	}

	_, err := loader.Load(context.Background(), def)
	if !errors.Is(err, ErrAdapterNameMismatch) {
		t.Fatalf("err = %v, want ErrAdapterNameMismatch", err)
	}
}

func TestLoaderLoadHandshakeTimeout(t *testing.T) {
	t.Parallel()

	loader := NewLoader(50 * time.Millisecond)
	def := Definition{
		Name:    "markdown",
		Kind:    AdapterKindView,
		Command: "sh",
		Args:    []string{"-c", "sleep 1"},
	}

	_, err := loader.Load(context.Background(), def)
	if !errors.Is(err, ErrPluginStartTimeout) {
		t.Fatalf("err = %v, want ErrPluginStartTimeout", err)
	}
}

func TestLoadedPluginHealthAfterExit(t *testing.T) {
	t.Parallel()

	loader := NewLoader(time.Second)
	def := Definition{
		Name:    "markdown",
		Kind:    AdapterKindView,
		Command: "sh",
		Args:    []string{"-c", "printf '{\"protocol_version\":\"todoopen.plugin.v1\",\"name\":\"markdown\",\"kind\":\"view\",\"capabilities\":[\"render_tasks\"],\"health\":{\"state\":\"ready\"}}\\n'; exit 0"},
	}

	p, err := loader.Load(context.Background(), def)
	if err != nil {
		t.Fatalf("load plugin: %v", err)
	}
	defer func() { _ = p.Close() }()

	time.Sleep(50 * time.Millisecond)
	if got := p.Health().State; got != HealthUnhealthy {
		t.Fatalf("health = %q, want unhealthy", got)
	}
}
