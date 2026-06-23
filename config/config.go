package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type EnvironmentConfig struct {
	Name       string   `toml:"name"`
	DSNEnv     string   `toml:"dsn_env"`
	Provider   string   `toml:"provider"`
	Database   string   `toml:"database"`
	Schema     string   `toml:"schema"`
	SchemaList []string `toml:"schemas"`
	Table      string   `toml:"table"`
}

func (e EnvironmentConfig) Schemas(defaultSchema string) []string {
	if len(e.SchemaList) > 0 {
		return e.SchemaList
	}
	if e.Schema != "" {
		return []string{e.Schema}
	}
	return []string{defaultSchema}
}

type Config struct {
	Environments map[string]EnvironmentConfig `toml:"environments"`
}

func Load(path string) (*Config, error) {
	resolved, err := resolve(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(resolved, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", resolved, err)
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
