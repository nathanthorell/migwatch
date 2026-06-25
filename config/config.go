package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Environments     map[string]EnvironmentConfig `toml:"environments"`
	EnvironmentOrder []string                     `toml:"-"`
}

func Load(path string) (*Config, error) {
	resolved, err := resolve(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	meta, err := toml.DecodeFile(resolved, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config %s: %w", resolved, err)
	}

	for _, k := range meta.Keys() {
		if len(k) == 2 && k[0] == "environments" {
			cfg.EnvironmentOrder = append(cfg.EnvironmentOrder, k[1])
		}
	}

	for key, env := range cfg.Environments {
		if env.Table == "" {
			env.Table = "flyway_schema_history"
			cfg.Environments[key] = env
		}
	}

	return &cfg, nil
}

func resolve(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	candidates := []string{
		"migwatch.toml",
		filepath.Join(userConfigDir(), "migwatch", "migwatch.toml"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("no config file found; checked: %v (use --config to specify)", candidates)
}

func userConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return dir
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}
