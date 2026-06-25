package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanthorell/migwatch/config"
	"github.com/nathanthorell/migwatch/model"
	"github.com/nathanthorell/migwatch/provider"
)

func loadEnvFile() error {
	path := envFile
	if path == "" {
		path = ".env"
	}

	if err := godotenv.Load(path); err != nil && (envFile != "" || !os.IsNotExist(err)) {
		return fmt.Errorf("load env file %q: %w", path, err)
	}

	return nil
}

func envLabel(key string, env config.EnvironmentConfig) string {
	if env.Name != "" {
		return env.Name
	}
	return key
}

func buildBannerSummaries(envs map[string]config.EnvironmentConfig, envOrder []string) []model.EnvironmentSummary {
	var summaries []model.EnvironmentSummary
	for _, key := range envOrder {
		env := envs[key]
		label := envLabel(key, env)
		s := model.EnvironmentSummary{Label: label}
		if rawDSN := os.Getenv(env.DSNEnv); rawDSN != "" {
			conn := config.BuildConnection(rawDSN)
			s.Driver = conn.Driver
			s.Host = conn.Host
			s.Database = conn.Database
			if env.Database != "" {
				s.Database = env.Database
			}
		}
		summaries = append(summaries, s)
	}
	return summaries
}

func resolveConn(ctx context.Context, env config.EnvironmentConfig) (model.Connection, error) {
	rawDSN := os.Getenv(env.DSNEnv)
	if rawDSN == "" {
		return model.Connection{}, fmt.Errorf("env var %q is not set", env.DSNEnv)
	}

	conn, err := config.OverrideDatabase(config.BuildConnection(rawDSN), env.Database)
	if err != nil {
		return model.Connection{}, err
	}

	return config.ResolveConnection(ctx, conn)
}

func fetchSchema(ctx context.Context, conn model.Connection, env config.EnvironmentConfig, schema string) model.EnvironmentResult {
	result := model.EnvironmentResult{Schema: schema}

	p, err := provider.New(env, schema)
	if err != nil {
		result.Error = err
		return result
	}

	migrations, err := p.FetchMigrations(ctx, conn)
	if err != nil {
		result.Error = config.WrapAuthError(err, conn)
		return result
	}

	result.Migrations = migrations
	return result
}
