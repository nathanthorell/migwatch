package model

import "time"

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
