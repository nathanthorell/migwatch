package cmd

import (
	"context"

	"github.com/nathanthorell/migwatch/config"
	"github.com/nathanthorell/migwatch/display"
	"github.com/nathanthorell/migwatch/model"
	"github.com/spf13/cobra"
)

func runCompare(cmd *cobra.Command, args []string) error {
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

	var results []model.EnvironmentResult
	for _, key := range envOrder {
		env := envs[key]
		label := envLabel(key, env)

		conn, err := resolveConn(ctx, env)
		if err != nil {
			results = append(results, model.EnvironmentResult{Environment: label, Error: err})
			continue
		}

		schema := env.Schemas(conn.Driver.DefaultSchema())[0]
		result := fetchSchema(ctx, conn, env, schema)
		result.Environment = label
		result.Database = conn.Database
		results = append(results, result)
	}

	display.PrintCompareTable(results)
	return nil
}
