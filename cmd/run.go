package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanthorell/migwatch/config"
	"github.com/nathanthorell/migwatch/display"
	"github.com/nathanthorell/migwatch/model"
	"github.com/nathanthorell/migwatch/provider"
	"github.com/spf13/cobra"
)

func runStatus(cmd *cobra.Command, args []string) error {
	if err := loadEnvFile(); err != nil {
		return err
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	envs := cfg.Environments
	if envFilter != "" {
		env, ok := envs[envFilter]
		if !ok {
			return fmt.Errorf("environment %q not found in config", envFilter)
		}
		envs = map[string]config.EnvironmentConfig{envFilter: env}
	}

	display.PrintBanner()

	ctx := context.Background()
	for key, env := range envs {
		result := fetchEnvironment(ctx, key, env)
		display.PrintMigrationTable(result)
	}

	return nil
}

func fetchEnvironment(ctx context.Context, key string, env config.EnvironmentConfig) model.EnvironmentResult {
	label := env.Name
	if label == "" {
		label = key
	}

	result := model.EnvironmentResult{
		Environment: label,
		Schema:      env.Schema,
	}

	dsn := os.Getenv(env.DSNEnv)
	if dsn == "" {
		result.Error = fmt.Errorf("env var %q is not set", env.DSNEnv)
		return result
	}

	dsn = config.AdjustDSN(dsn)
	result.Database = config.DatabaseFromDSN(dsn)

	p, err := provider.New(env)
	if err != nil {
		result.Error = err
		return result
	}

	migrations, err := p.FetchMigrations(ctx, dsn)
	if err != nil {
		result.Error = config.WrapAuthError(err, dsn)
		return result
	}

	result.Migrations = migrations
	return result
}

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
