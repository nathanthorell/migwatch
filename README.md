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
│   └── config.go                # Config code
├── model/
│   └── migration.go             # Shared types
├── provider/
│   ├── provider.go              # MigrationProvider interface
│   └── flyway/
│       └── flyway.go            # Flyway implementation
├── display/
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

[environments.test]
name     = "Test"
dsn_env  = "TEST_DSN"
provider = "flyway"
schema   = "dbo"
table    = "flyway_schema_history"
```

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
