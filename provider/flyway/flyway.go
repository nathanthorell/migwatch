package flyway

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/azuread"
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

func (p *Provider) FetchMigrations(ctx context.Context, conn model.Connection) ([]model.Migration, error) {
	db, err := sql.Open(goDriverName(conn), conn.DSN)
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
			COALESCE(checksum, 0),
			success
		FROM %s
		ORDER BY installed_rank ASC`,
		tableRef(conn.Driver, p.schema, p.table),
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
			&m.Checksum,
			&m.Success,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}

func goDriverName(conn model.Connection) string {
	switch conn.Driver {
	case model.DriverPostgres:
		return "postgres"
	default:
		if conn.AuthMethod != "" {
			return "azuresql"
		}
		return "sqlserver"
	}
}

func tableRef(driver model.Driver, schema, table string) string {
	switch driver {
	case model.DriverPostgres:
		return fmt.Sprintf(`"%s"."%s"`, schema, table)
	default:
		return fmt.Sprintf("[%s].[%s]", schema, table)
	}
}
