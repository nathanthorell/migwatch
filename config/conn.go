package config

import (
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

func driverFromScheme(scheme string) model.Driver {
	switch scheme {
	case "postgres", "postgresql":
		return model.DriverPostgres
	default:
		return model.DriverMSSQL
	}
}
