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
   ./migwatch status

   # Show status for a specific environment
   ./migwatch status --env dev
   ```

## Configuration

### Project Structure

```none
├── /
│── main.go                      # App entry point
├── cmd/                         # CLI commands
│   ├── root.go                  # Root command
│   └── status.go                # Status command
├── config/                      # Configuration management
│   ├── config.go                  # Config code
├── migwatch.toml                # App config (gitignored)
├── migwatch.toml.example        # Example config
├── .env                         # Connection strings (gitignored)
└── .env.example                 # Example environment file
```
