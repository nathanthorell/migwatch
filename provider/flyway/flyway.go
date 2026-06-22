package flyway

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/nathanthorell/migwatch/model"
)

type Provider struct {
	schema string
	table  string
}

func New(schema, table string) *Provider {
	return &Provider{schema: schema, table: table}
}

func (p *Provider) Name() string { return "flyway" }

func (p *Provider) FetchMigrations(ctx context.Context, dsn string) ([]model.Migration, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			installed_rank,
			COALESCE(version, ''),
			description,
			type,
			script,
			installed_by,
			installed_on,
			execution_time,
			success
		FROM [%s].[%s]
		ORDER BY installed_rank ASC`,
		p.schema, p.table,
	)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var migrations []model.Migration
	for rows.Next() {
		var m model.Migration
		if err := rows.Scan(
			&m.InstalledRank,
			&m.Version,
			&m.Description,
			&m.Type,
			&m.Script,
			&m.InstalledBy,
			&m.InstalledOn,
			&m.ExecutionTime,
			&m.Success,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}
