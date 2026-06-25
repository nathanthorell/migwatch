package cmd

import (
	"context"
	"fmt"

	"github.com/nathanthorell/migwatch/config"
	"github.com/nathanthorell/migwatch/display"
	"github.com/nathanthorell/migwatch/model"
	"github.com/spf13/cobra"
)

func runSummary(cmd *cobra.Command, args []string) error {
	if err := loadEnvFile(); err != nil {
		return err
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	envs, envOrder, err := cfg.OrderedEnvs(envFilter)
	if err != nil {
		return err
	}

	display.PrintBanner(buildBannerSummaries(envs, envOrder))

	ctx := context.Background()
	for _, key := range envOrder {
		env := envs[key]
		label := envLabel(key, env)

		conn, err := resolveConn(ctx, env)
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
			display.PrintSummary(result)
		}
	}

	return nil
}
