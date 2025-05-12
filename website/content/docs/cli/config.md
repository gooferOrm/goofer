# Configuration

The Goofer ORM CLI can be configured using a configuration file, environment variables, or command-line flags. This page explains how to configure the CLI to match your project's requirements.

## Configuration File

The default configuration file is `goofer.yaml` in your project's root directory. You can specify a different file using the `--config` flag or the `GOOFER_CONFIG` environment variable.

### Basic Configuration

A basic configuration file looks like this:

```yaml
# Database configuration
database:
  dialect: sqlite
  dsn: ./database.db

# Directory configuration
directories:
  entities: ./entity
  repositories: ./repository
  migrations: ./migrations

# Generation configuration
generation:
  timestamps: true
  package:
    entity: entity
    repository: repository
```

### Full Configuration

Here's a complete configuration file with all available options:

```yaml
# Database configuration
database:
  # Database dialect (sqlite, mysql, postgres)
  dialect: sqlite
  # Database connection string
  dsn: ./database.db
  # Maximum number of open connections
  max_open_conns: 10
  # Maximum number of idle connections
  max_idle_conns: 5
  # Connection lifetime in seconds
  conn_max_lifetime: 3600

# Directory configuration
directories:
  # Directory for entity files
  entities: ./entity
  # Directory for repository files
  repositories: ./repository
  # Directory for migration files
  migrations: ./migrations
  # Directory for template files
  templates: ./.goofer/templates

# Generation configuration
generation:
  # Add timestamps to generated entities
  timestamps: true
  # Package names
  package:
    # Package name for entities
    entity: entity
    # Package name for repositories
    repository: repository
  # Template configuration
  templates:
    # Use custom templates
    custom: false
    # Template extension
    extension: .tmpl

# Migration configuration
migration:
  # Table name for migration records
  table: migrations
  # Automatically run migrations on startup
  auto: false

# Logging configuration
logging:
  # Log level (debug, info, warn, error)
  level: info
  # Log format (text, json)
  format: text
  # Log file (stdout, stderr, or file path)
  output: stdout
```

## Configuration Commands

The CLI provides commands to manage your configuration:

### Initialize Configuration

To create a new configuration file:

```bash
goofer config init
```

This will create a `goofer.yaml` file in your project's root directory with default settings.

You can specify a different file:

```bash
goofer config init --file ./config/goofer.yaml
```

### Show Configuration

To show the current configuration:

```bash
goofer config show
```

This will display the merged configuration from all sources (file, environment variables, and defaults).

### Set Configuration Values

To set a configuration value:

```bash
goofer config set database.dialect postgres
```

This will update the configuration file with the new value.

You can set nested values using dot notation:

```bash
goofer config set database.max_open_conns 20
```

## Environment Variables

You can override configuration values using environment variables. The environment variables are prefixed with `GOOFER_` and use underscores instead of dots for nested values.

For example:

| Configuration | Environment Variable |
|---------------|----------------------|
| `database.dialect` | `GOOFER_DATABASE_DIALECT` |
| `database.dsn` | `GOOFER_DATABASE_DSN` |
| `directories.migrations` | `GOOFER_DIRECTORIES_MIGRATIONS` |

### Common Environment Variables

| Environment Variable | Description |
|----------------------|-------------|
| `GOOFER_CONFIG` | Path to the configuration file |
| `GOOFER_DATABASE_DIALECT` | Database dialect (sqlite, mysql, postgres) |
| `GOOFER_DATABASE_DSN` | Database connection string |
| `GOOFER_DIRECTORIES_MIGRATIONS` | Directory for migration files |
| `GOOFER_LOGGING_LEVEL` | Log level (debug, info, warn, error) |

## Command-Line Flags

You can also override configuration values using command-line flags. The flags take precedence over environment variables and the configuration file.

### Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--config`, `-c` | Path to the configuration file |
| `--verbose`, `-v` | Enable verbose output (overrides logging.level) |

### Command-Specific Flags

Many commands have flags that override specific configuration values:

| Command | Flag | Configuration |
|---------|------|---------------|
| `migrate` | `--dir`, `-d` | `directories.migrations` |
| `migrate` | `--dialect` | `database.dialect` |
| `generate entity` | `--package`, `-p` | `generation.package.entity` |
| `generate repository` | `--package`, `-p` | `generation.package.repository` |

## Configuration Precedence

Configuration values are resolved in the following order (highest precedence first):

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

This means that a value specified via a command-line flag will override the same value specified in an environment variable, which will override the value in the configuration file.

## Database Configuration

### SQLite

```yaml
database:
  dialect: sqlite
  dsn: ./database.db
```

### MySQL

```yaml
database:
  dialect: mysql
  dsn: user:password@tcp(localhost:3306)/dbname?parseTime=true
```

### PostgreSQL

```yaml
database:
  dialect: postgres
  dsn: postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## Directory Configuration

You can customize the directories used by Goofer ORM:

```yaml
directories:
  entities: ./src/domain/entity
  repositories: ./src/infrastructure/repository
  migrations: ./src/infrastructure/migrations
  templates: ./src/infrastructure/templates
```

## Generation Configuration

You can customize the code generation behavior:

```yaml
generation:
  timestamps: true
  package:
    entity: domain.entity
    repository: infrastructure.repository
  templates:
    custom: true
    extension: .gotmpl
```

## Migration Configuration

You can customize the migration behavior:

```yaml
migration:
  table: goofer_migrations
  auto: true
```

## Logging Configuration

You can customize the logging behavior:

```yaml
logging:
  level: debug
  format: json
  output: ./logs/goofer.log
```

## Multiple Environments

You can use different configuration files for different environments:

```bash
# Development
goofer --config ./config/goofer.dev.yaml migrate up

# Production
goofer --config ./config/goofer.prod.yaml migrate up
```

Or use environment variables:

```bash
# Development
GOOFER_DATABASE_DSN=./dev.db goofer migrate up

# Production
GOOFER_DATABASE_DSN=postgres://user:password@prod-db:5432/app goofer migrate up
```

## Best Practices

### Version Control

Include your configuration file in version control, but use environment variables or a separate configuration file for sensitive information like database credentials.

### Environment-Specific Configuration

Use environment variables or different configuration files for environment-specific values:

- `goofer.yaml`: Common configuration
- `goofer.dev.yaml`: Development-specific configuration
- `goofer.prod.yaml`: Production-specific configuration

### Sensitive Information

Don't store sensitive information like passwords in your configuration file. Use environment variables instead:

```bash
GOOFER_DATABASE_DSN=postgres://user:password@localhost:5432/dbname goofer migrate up
```

### Documentation

Document your configuration choices, especially if you're using custom directories or templates.

## Examples

### Basic SQLite Configuration

```yaml
database:
  dialect: sqlite
  dsn: ./database.db
```

### Production MySQL Configuration

```yaml
database:
  dialect: mysql
  dsn: ${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/${MYSQL_DATABASE}?parseTime=true
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: 3600
```

### Custom Directory Structure

```yaml
directories:
  entities: ./internal/domain/entity
  repositories: ./internal/infrastructure/repository
  migrations: ./internal/infrastructure/migrations
```

### Custom Package Names

```yaml
generation:
  package:
    entity: internal.domain.entity
    repository: internal.infrastructure.repository
```

## Troubleshooting

### Configuration Not Found

If you get a "configuration file not found" error, make sure:

1. The file exists at the specified path
2. The file has the correct permissions
3. You're running the command from the correct directory

### Invalid Configuration

If you get an "invalid configuration" error, check:

1. The YAML syntax is correct
2. The values have the correct types
3. Required fields are present

### Environment Variables Not Applied

If your environment variables aren't being applied, check:

1. The variable names are correct (prefixed with `GOOFER_`)
2. The variable names use underscores for nested values
3. The variables are set in the current environment

## Next Steps

- Learn about [Commands](./commands) for a complete list of available commands
- Explore [Migration Commands](./migration) for database schema management
- Check out [Generate Commands](./generate) for code generation