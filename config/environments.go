package config

import "fmt"

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

func (c *Config) OrderedEnvs(filter string) (map[string]EnvironmentConfig, []string, error) {
	if filter != "" {
		env, ok := c.Environments[filter]
		if !ok {
			return nil, nil, fmt.Errorf("environment %q not found in config", filter)
		}
		return map[string]EnvironmentConfig{filter: env}, []string{filter}, nil
	}
	return c.Environments, c.EnvironmentOrder, nil
}
