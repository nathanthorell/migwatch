# migwatch

migwatch is a CLI tool for visualizing database migration state across environments. It connects directly to your databases and displays migration history from your migration tool's schema history table.

## Quick Start

1. **Copy example files**:

   ```bash
   cp .env.example .env
   cp migwatch.toml.example migwatch.toml
   ```

2. **Edit `.env`** with your database connection strings

3. **Edit `migwatch.toml`** with your environment definitions

4. **Run**:

   ```bash
   # Show migration status for all environments
   ./migwatch

   # Show status for a specific environment
   ./migwatch --env dev
   ```

## Configuration

### Project Structure

```none
├── /
├── main.go                      # App entry point
├── cmd/                         # CLI commands
│   ├── root.go                  # Command definition and flags
│   └── run.go                   # Execution logic
├── config/                      # Configuration management
│   ├── config.go                # TOML types and loading
│   ├── conn.go                  # DSN parsing and connection building
│   └── auth.go                  # Token resolution and auth error wrapping
├── model/
│   └── models.go                # Shared types
├── provider/
│   ├── provider.go              # MigrationProvider interface
│   └── flyway/
│       └── flyway.go            # Flyway implementation
├── display/
│   ├── display.go               # Banner and environment header output
│   ├── styles.go                # Lipgloss style definitions
│   └── table.go                 # Styled table output
├── migwatch.toml                # App config (gitignored)
├── migwatch.toml.example        # Example config
├── .env                         # Connection strings (gitignored)
└── .env.example                 # Example environment file
```

### Application Configuration

Configuration is managed through `migwatch.toml`. Each environment block references a `.env` variable by name and sets provider-specific options:

```toml
[environments.dev]
name     = "Dev"
dsn_env  = "DEV_DSN"
provider = "flyway"
schema   = "dbo"
table    = "flyway_schema_history"

# multiple schemas on the same database
[environments.dev_multi]
name     = "Dev"
dsn_env  = "DEV_DSN"
provider = "flyway"
schemas  = ["dbo", "dba"]
table    = "flyway_schema_history"

# database override: one shared DSN, different database per environment
[environments.staging]
name     = "Staging"
dsn_env  = "DEV_DSN"
provider = "flyway"
database = "stagingdb"
schema   = "dbo"
table    = "flyway_schema_history"
```

### Authentication (SQL Server)

Connection strings go in `.env`. Use the `fedauth` parameter to select the auth method:

```env
# SQL auth
DEV_DSN=sqlserver://user:pass@myserver?database=mydb

# Entra - az login token
DEV_DSN=sqlserver://myserver.database.windows.net?database=mydb&fedauth=ActiveDirectoryAzCli

# Entra - interactive browser
DEV_DSN=sqlserver://myserver.database.windows.net?database=mydb&fedauth=ActiveDirectoryInteractive

# Entra - default credential chain
DEV_DSN=sqlserver://myserver.database.windows.net?database=mydb&fedauth=ActiveDirectoryDefault
```

Interactive auth also requires `AZURE_CLIENT_ID` in `.env` - see `.env.example` for the default value.

### Authentication (PostgreSQL)

```env
# Password auth
DEV_DSN=postgres://user:pass@myserver:5432/mydb?sslmode=disable

# Entra - az login token (Azure Flexible Server)
DEV_DSN=postgres://user%40example.com@myserver.postgres.database.azure.com/mydb?sslmode=require&fedauth=ActiveDirectoryAzCli

# Entra - default credential chain (Azure Flexible Server)
DEV_DSN=postgres://user%40example.com@myserver.postgres.database.azure.com/mydb?sslmode=require&fedauth=ActiveDirectoryDefault
```

Note: the `@` in the username should be percent-encoded as `%40` - both forms work, but `%40` is safer and standards-compliant.

## Building and Running

```bash
# Build the application
go build -o ./build/migwatch .

# Show migration status for all environments
./build/migwatch

# Filter to one environment
./build/migwatch --env dev

# Use a custom config or env file
./build/migwatch --config /path/to/migwatch.toml --env-file /path/to/.env
```

## Flags

- `-e, --env <name>` - Filter output to a single environment
- `--config <path>` - Path to config file
- `--env-file <path>` - Path to `.env` file
- `-h, --help` - help for migwatch

## Contributing

This is a work in progress. Contributions and suggestions are welcome.
