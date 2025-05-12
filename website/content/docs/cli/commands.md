# Goofer ORM CLI Commands

This page provides an overview of all available commands in the Goofer ORM CLI. For more detailed information about specific command groups, see their dedicated pages.

## Command Structure

The Goofer ORM CLI uses a hierarchical command structure:

```
goofer [command] [subcommand] [options]
```

Where:
- `command` is the main command group (e.g., `migrate`, `generate`)
- `subcommand` is a specific action within that group (e.g., `migrate up`, `generate entity`)
- `options` are flags and arguments that modify the command's behavior

## Global Options

These options can be used with any command:

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help for a command |
| `--version`, `-v` | Show the CLI version |
| `--verbose` | Enable verbose output |
| `--config`, `-c` | Specify a config file (default: ./goofer.yaml) |

## Available Commands

### Root Commands

| Command | Description |
|---------|-------------|
| `goofer version` | Show the CLI version |
| `goofer help` | Show help for a command |
| `goofer init` | Initialize a new Goofer ORM project |
| `goofer completion` | Generate shell completion scripts |

### Migration Commands

| Command | Description |
|---------|-------------|
| `goofer migrate create <name>` | Create a new migration |
| `goofer migrate up` | Apply pending migrations |
| `goofer migrate down` | Revert the last migration |
| `goofer migrate status` | Show migration status |
| `goofer migrate reset` | Revert all migrations |

See [Migration Commands](./migration) for more details.

### Generate Commands

| Command | Description |
|---------|-------------|
| `goofer generate entity <name>` | Generate a new entity |
| `goofer generate repository <entity>` | Generate a repository for an entity |
| `goofer generate migration` | Generate a migration from entity definitions |
| `goofer generate all` | Generate all artifacts for an entity |

See [Generate Commands](./generate) for more details.

### Configuration Commands

| Command | Description |
|---------|-------------|
| `goofer config init` | Initialize a new configuration file |
| `goofer config show` | Show the current configuration |
| `goofer config set <key> <value>` | Set a configuration value |

See [Configuration](./config) for more details.

## Command Examples

### Initialize a New Project

```bash
# Initialize a new project with default settings
goofer init my-project

# Initialize with a specific database
goofer init my-project --db postgres

# Initialize with a specific directory structure
goofer init my-project --template standard
```

### Working with Migrations

```bash
# Create a new migration
goofer migrate create create_users_table

# Apply pending migrations
goofer migrate up

# Apply a specific number of migrations
goofer migrate up --steps 2

# Revert the last migration
goofer migrate down

# Show migration status
goofer migrate status
```

### Generating Code

```bash
# Generate a User entity
goofer generate entity User --fields "id:uint:primaryKey,autoIncrement name:string:notnull email:string:unique,notnull"

# Generate a repository for the User entity
goofer generate repository User

# Generate a migration from entity definitions
goofer generate migration initial_schema

# Generate all artifacts for an entity
goofer generate all User
```

### Configuration Management

```bash
# Initialize a new configuration file
goofer config init

# Show the current configuration
goofer config show

# Set a configuration value
goofer config set database.dialect postgres
```

## Command Help

You can get help for any command by using the `--help` flag:

```bash
# Get general help
goofer --help

# Get help for a specific command
goofer migrate --help

# Get help for a specific subcommand
goofer migrate create --help
```

The help output includes:
- Command description
- Usage syntax
- Available options
- Examples

## Command Aliases

Some commands have aliases for convenience:

| Command | Alias |
|---------|-------|
| `goofer migrate` | `goofer m` |
| `goofer generate` | `goofer g` |
| `goofer version` | `goofer v` |
| `goofer help` | `goofer h` |

For example, you can use `goofer m up` instead of `goofer migrate up`.

## Exit Codes

The CLI uses the following exit codes:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Command-line parsing error |
| 3 | Configuration error |
| 4 | Database error |
| 5 | Migration error |
| 6 | Generation error |

You can use these exit codes in scripts to check if a command succeeded or failed.

## Environment Variables

The CLI supports the following environment variables:

| Variable | Description |
|----------|-------------|
| `GOOFER_CONFIG` | Path to the configuration file |
| `GOOFER_DATABASE_URL` | Database connection URL |
| `GOOFER_VERBOSE` | Enable verbose output (set to "true") |

Environment variables take precedence over configuration file values.

## Next Steps

- Learn about [Migration Commands](./migration) for database schema management
- Explore [Generate Commands](./generate) for code generation
- See [Configuration](./config) for customizing the CLI behavior