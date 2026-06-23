package model

import "time"

type Driver string

const (
	DriverMSSQL      Driver = "mssql"
	DriverPostgres   Driver = "postgres"
	DriverMySQL      Driver = "mysql"
	DriverDatabricks Driver = "databricks"
)

func (d Driver) DefaultSchema() string {
	switch d {
	case DriverPostgres:
		return "public"
	default:
		return "dbo"
	}
}

type Connection struct {
	DSN        string
	Driver     Driver
	Database   string
	AuthMethod string
}

type Migration struct {
	InstalledRank int
	Version       string
	Description   string
	Type          string
	Script        string
	InstalledBy   string
	InstalledOn   time.Time
	ExecutionTime int
	Success       bool
}

type EnvironmentResult struct {
	Environment string
	Database    string
	Schema      string
	Migrations  []Migration
	Error       error
}
