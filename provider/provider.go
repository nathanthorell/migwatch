package provider

import (
	"context"
	"fmt"

	"github.com/nathanthorell/migwatch/config"
	"github.com/nathanthorell/migwatch/model"
	"github.com/nathanthorell/migwatch/provider/flyway"
)

type MigrationProvider interface {
	Name() string
	FetchMigrations(ctx context.Context, dsn string) ([]model.Migration, error)
}

func New(cfg config.EnvironmentConfig) (MigrationProvider, error) {
	switch cfg.Provider {
	case "flyway":
		return flyway.New(cfg.Schema, cfg.Table), nil
	default:
		return nil, fmt.Errorf("unknown provider: %q", cfg.Provider)
	}
}
