package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type EnvironmentConfig struct {
	Name     string `toml:"name"`
	DSNEnv   string `toml:"dsn_env"`
	Provider string `toml:"provider"`
	Schema   string `toml:"schema"`
	Table    string `toml:"table"`
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
		if env.Schema == "" {
			env.Schema = "dbo"
			cfg.Environments[key] = env
		}
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

func DatabaseFromDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return ""
	}
	return u.Query().Get("database")
}

// AdjustDSN injects applicationclientid from AZURE_CLIENT_ID env var when using ActiveDirectoryInteractive.
func AdjustDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	q := u.Query()
	if q.Get("fedauth") != "ActiveDirectoryInteractive" {
		return dsn
	}
	if q.Get("applicationclientid") != "" {
		return dsn
	}
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		return dsn
	}
	q.Set("applicationclientid", clientID)
	u.RawQuery = q.Encode()
	return u.String()
}

func AuthMethodFromDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "sql"
	}
	if fedauth := u.Query().Get("fedauth"); fedauth != "" {
		return fedauth
	}
	return "sql"
}

func WrapAuthError(err error, dsn string) error {
	msg := err.Error()
	if !strings.Contains(msg, "Login failed") && !strings.Contains(msg, "login error") {
		return err
	}

	switch AuthMethodFromDSN(dsn) {
	case "ActiveDirectoryAzCli":
		return fmt.Errorf("authentication failed: az login token missing or expired — run `az login`")
	case "ActiveDirectoryInteractive":
		return fmt.Errorf("authentication failed: interactive login did not complete")
	case "ActiveDirectoryDefault":
		return fmt.Errorf("authentication failed: no valid credential found in default chain (az login, env vars, managed identity)")
	default:
		return fmt.Errorf("authentication failed: invalid username or password")
	}
}

func userConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return dir
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}
