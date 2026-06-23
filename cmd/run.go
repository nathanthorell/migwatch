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
		label := env.Name
		if label == "" {
			label = key
		}

		rawDSN := os.Getenv(env.DSNEnv)
		if rawDSN == "" {
			display.PrintEnvironmentHeader(model.EnvironmentResult{Environment: label})
			fmt.Printf("  error: env var %q is not set\n", env.DSNEnv)
			fmt.Println()
			continue
		}

		conn, err := config.ResolveConnection(ctx, config.BuildConnection(rawDSN))
		if err != nil {
			display.PrintEnvironmentHeader(model.EnvironmentResult{Environment: label})
			fmt.Printf("  error: %v\n\n", err)
			continue
		}

		display.PrintEnvironmentHeader(model.EnvironmentResult{
			Environment: label,
			Database:    conn.Database,
		})

		schemas := env.Schemas(conn.Driver.DefaultSchema())
		for _, schema := range schemas {
			result := fetchSchema(ctx, conn, env, schema)
			if len(schemas) > 1 {
				display.PrintSchemaLabel(schema)
			}
			display.PrintMigrationTable(result)
		}
	}

	return nil
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
