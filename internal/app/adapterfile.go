package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/justEstif/todo-open/internal/plugin"
)

const (
	configFileName = "config.toml"
)

// AdapterEntry holds registration and settings for a single adapter binary.
type AdapterEntry struct {
	Bin    string             `toml:"bin"`
	Kind   plugin.AdapterKind `toml:"kind"`
	Args   []string           `toml:"args"`
	Config map[string]string  `toml:"config"`
}

// AdapterFileConfig is the fully resolved adapter configuration loaded from
// .todoopen/config.toml. Adapter config is kept separate from workspace
// metadata (meta.json) so each file has a single concern.
type AdapterFileConfig struct {
	Views    AdapterGroupConfig      `toml:"views"`
	Sync     AdapterGroupConfig      `toml:"sync"`
	Adapters map[string]AdapterEntry `toml:"adapters"`
}

// AdapterGroupConfig controls which adapters are active for a given kind.
type AdapterGroupConfig struct {
	Enabled []string `toml:"enabled"`
}

// ExpandEnv replaces ${VAR} references in each adapter's config map with their
// environment variable values. Only adapter config maps are expanded — binary
// paths, args, and kind are not touched.
func (c *AdapterFileConfig) ExpandEnv() {
	for name, entry := range c.Adapters {
		if entry.Config == nil {
			continue
		}
		expanded := make(map[string]string, len(entry.Config))
		for k, v := range entry.Config {
			expanded[k] = os.ExpandEnv(v)
		}
		entry.Config = expanded
		c.Adapters[name] = entry
	}
}

// ApplyDefaults fills in zero-value fields with sensible defaults.
func (c *AdapterFileConfig) ApplyDefaults() {
	if len(c.Views.Enabled) == 0 {
		c.Views.Enabled = []string{"json"}
	}
	if len(c.Sync.Enabled) == 0 {
		c.Sync.Enabled = []string{"noop"}
	}
}

// Validate checks the config for consistency.
func (c *AdapterFileConfig) Validate() error {
	if err := validateNames("views.enabled", c.Views.Enabled); err != nil {
		return err
	}
	if err := validateNames("sync.enabled", c.Sync.Enabled); err != nil {
		return err
	}

	for name, entry := range c.Adapters {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("adapter name must not be empty")
		}
		if entry.Kind != plugin.AdapterKindView && entry.Kind != plugin.AdapterKindSync {
			return fmt.Errorf("adapters.%s.kind must be one of: view, sync", name)
		}
		if strings.TrimSpace(entry.Bin) == "" {
			return fmt.Errorf("adapters.%s.bin is required", name)
		}
	}

	// Every enabled adapter that is not a builtin must have a declared entry.
	for _, name := range c.Views.Enabled {
		if name == "json" {
			continue // builtin
		}
		if _, ok := c.Adapters[name]; !ok {
			return fmt.Errorf("views.enabled: adapter %q is not declared in [adapters]", name)
		}
	}
	for _, name := range c.Sync.Enabled {
		if name == "noop" {
			continue // builtin
		}
		if _, ok := c.Adapters[name]; !ok {
			return fmt.Errorf("sync.enabled: adapter %q is not declared in [adapters]", name)
		}
	}

	return nil
}

// DefaultAdapterFileConfig returns defaults used when config.toml is absent.
func DefaultAdapterFileConfig() AdapterFileConfig {
	cfg := AdapterFileConfig{}
	cfg.ApplyDefaults()
	return cfg
}

// LoadAdapterFileConfig reads .todoopen/config.toml from the workspace root.
// If the file does not exist the default config is returned.
func LoadAdapterFileConfig(workspaceRoot string) (AdapterFileConfig, error) {
	cfgPath := filepath.Join(workspaceRoot, metaDirName, configFileName)

	data, err := os.ReadFile(cfgPath)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultAdapterFileConfig(), nil
	}
	if err != nil {
		return AdapterFileConfig{}, fmt.Errorf("read adapter config: %w", err)
	}

	var cfg AdapterFileConfig
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return AdapterFileConfig{}, fmt.Errorf("decode adapter config: %w", err)
	}

	cfg.ApplyDefaults()
	cfg.ExpandEnv()

	if err := cfg.Validate(); err != nil {
		return AdapterFileConfig{}, fmt.Errorf("invalid adapter config: %w", err)
	}

	return cfg, nil
}

// pluginsFromAdapterFileConfig converts AdapterFileConfig entries into the
// AdapterPluginConfig slice expected by the runtime resolver.
func pluginsFromAdapterFileConfig(cfg AdapterFileConfig) []AdapterPluginConfig {
	var plugins []AdapterPluginConfig
	for name, entry := range cfg.Adapters {
		plugins = append(plugins, AdapterPluginConfig{
			Name:    name,
			Kind:    entry.Kind,
			Command: entry.Bin,
			Args:    entry.Args,
		})
	}
	return plugins
}
