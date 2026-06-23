package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/nathanthorell/migwatch/model"
)

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
		Host:       u.Hostname(),
		Database:   database,
		AuthMethod: authMethod,
	}
}

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

// OverrideDatabase rewrites conn to use the given database name, updating both
// the Connection fields and the DSN so the driver connects to the right database.
func OverrideDatabase(conn model.Connection, database string) (model.Connection, error) {
	if database == "" {
		return conn, nil
	}
	u, err := url.Parse(conn.DSN)
	if err != nil {
		return conn, fmt.Errorf("parse DSN: %w", err)
	}

	switch conn.Driver {
	case model.DriverPostgres:
		u.Path = "/" + database
	default:
		q := u.Query()
		q.Set("database", database)
		u.RawQuery = q.Encode()
	}

	conn.Database = database
	conn.DSN = u.String()
	return conn, nil
}

func driverFromScheme(scheme string) model.Driver {
	switch scheme {
	case "postgres", "postgresql":
		return model.DriverPostgres
	default:
		return model.DriverMSSQL
	}
}
