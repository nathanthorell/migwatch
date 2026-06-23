package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nathanthorell/migwatch/model"
)

type EnvironmentConfig struct {
	Name       string   `toml:"name"`
	DSNEnv     string   `toml:"dsn_env"`
	Provider   string   `toml:"provider"`
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

// BuildConnection parses a raw DSN into a Connection, applying any needed adjustments.
func BuildConnection(rawDSN string) model.Connection {
	dsn := AdjustDSN(rawDSN)
	u, err := url.Parse(dsn)
	if err != nil {
		return model.Connection{DSN: dsn}
	}

	driver := driverFromScheme(u.Scheme)
	authMethod := u.Query().Get("fedauth")

	var database string
	switch driver {
	case model.DriverPostgres:
		database = strings.TrimPrefix(u.Path, "/")
	default:
		database = u.Query().Get("database")
	}

	return model.Connection{
		DSN:        dsn,
		Driver:     driver,
		Database:   database,
		AuthMethod: authMethod,
	}
}

func WrapAuthError(err error, conn model.Connection) error {
	msg := err.Error()
	if !strings.Contains(msg, "Login failed") && !strings.Contains(msg, "login error") {
		return err
	}

	switch conn.AuthMethod {
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

func driverFromScheme(scheme string) model.Driver {
	switch scheme {
	case "postgres", "postgresql":
		return model.DriverPostgres
	default:
		return model.DriverMSSQL
	}
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
