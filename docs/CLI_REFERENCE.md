# Goofer ORM CLI Reference

## Overview

The Goofer ORM Command Line Interface (CLI) provides a set of tools to streamline development with the Goofer ORM framework. This reference documents all available commands, their options, and usage examples.

## Global Options

The following options are available for all commands:

| Option | Shorthand | Description |
|--------|-----------|-------------|
| `--verbose` | `-v` | Enable verbose output |
| `--config` | `-c` | Config file (default is ./goofer.yaml) |
| `--help` | `-h` | Help for any command |

## Core Commands

### goofer

```
goofer [command]
```

The root command. When run without any subcommands, displays help information.

**Options:**
- `--version`: Print the version number and exit

### goofer version

```
goofer version
```

Displays the current version of Goofer ORM CLI.

## Project Management

### goofer init

```
goofer init [project-name]
```

Initializes a new Goofer ORM project with the recommended directory structure and files.

**Arguments:**
- `project-name`: Name of the project (optional, defaults to current directory name)

**Options:**
- `--dialect`, `-d`: Database dialect to use (sqlite, mysql, postgres) (default: "sqlite")
- `--with-examples`: Include example entity models in the project

**Examples:**
```
goofer init my-app
goofer init my-app --dialect=postgres --with-examples
```

## Code Generation

### goofer generate entity

```
goofer generate entity [name] [field:type:tag...]
```

Generates a new entity struct with the specified fields and ORM tags.

**Arguments:**
- `name`: Entity name (required)
- `field:type:tag...`: Field definitions (optional)

**Options:**
- `--out`, `-o`: Output directory for generated code (default: ".")
- `--package`, `-p`: Package name for generated code (default: "models")
- `--with-validate`: Add validation tags
- `--with-hooks`: Add lifecycle hooks

**Examples:**
```
goofer generate entity User
goofer generate entity User id:uint:primaryKey,autoIncrement name:string:notnull email:string:unique,notnull
goofer generate entity User --with-hooks --with-validate
```

## Database Management

### goofer migrate create

```
goofer migrate create [name]
```

Creates a new migration with up/down SQL files.

**Arguments:**
- `name`: Migration name (required)

**Options:**
- `--migrations-dir`, `-d`: Directory for migration files (default: "migrations")
- `--dialect`, `-t`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--db-url`, `-u`: Database connection URL
- `--provider`, `-p`: Migration provider (sql, gorm) (default: "sql")

**Example:**
```
goofer migrate create add_users_table
```

### goofer migrate up

```
goofer migrate up
```

Runs all pending migrations that have not yet been applied.

**Options:**
- `--migrations-dir`, `-d`: Directory for migration files (default: "migrations")
- `--dialect`, `-t`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--db-url`, `-u`: Database connection URL
- `--provider`, `-p`: Migration provider (sql, gorm) (default: "sql")

### goofer migrate down

```
goofer migrate down
```

Rolls back the most recently applied migration.

**Options:**
- `--migrations-dir`, `-d`: Directory for migration files (default: "migrations")
- `--dialect`, `-t`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--db-url`, `-u`: Database connection URL
- `--provider`, `-p`: Migration provider (sql, gorm) (default: "sql")

### goofer migrate status

```
goofer migrate status
```

Displays the current status of all migrations.

**Options:**
- `--migrations-dir`, `-d`: Directory for migration files (default: "migrations")
- `--dialect`, `-t`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--db-url`, `-u`: Database connection URL
- `--provider`, `-p`: Migration provider (sql, gorm) (default: "sql")

## Schema Management

### goofer schema generate

```
goofer schema generate
```

Generates SQL schema DDL from registered entities.

**Options:**
- `--dialect`, `-d`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--entities-dir`, `-e`: Directory containing entity definitions (default: ".")
- `--package`, `-p`: Package name for entity definitions (default: "models")
- `--output`, `-o`: Output file for generated schema (default: "schema.sql")

### goofer schema dump

```
goofer schema dump
```

Exports the current database schema as SQL statements.

**Options:**
- `--dialect`, `-d`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--entities-dir`, `-e`: Directory containing entity definitions (default: ".")
- `--package`, `-p`: Package name for entity definitions (default: "models")
- `--output`, `-o`: Output file for schema dump (default: "dump.sql")

### goofer schema diff

```
goofer schema diff
```

Compares entity schemas with database schemas and shows differences.

**Options:**
- `--dialect`, `-d`: Database dialect (sqlite, mysql, postgres) (default: "sqlite")
- `--entities-dir`, `-e`: Directory containing entity definitions (default: ".")
- `--package`, `-p`: Package name for entity definitions (default: "models")

## Configuration File

Goofer CLI supports configuration via a `goofer.yaml` file. This file can contain default values for common options. Example:

```yaml
# goofer.yaml
dialect: postgres
migrations_dir: ./migrations
entities_dir: ./internal/models
package: models
db_url: host=localhost port=5432 user=postgres password=postgres dbname=myapp sslmode=disable
```

## Environment Variables

The following environment variables can be used to configure the CLI:

- `GOOFER_DIALECT`: Default database dialect (sqlite, mysql, postgres)
- `GOOFER_MIGRATIONS_DIR`: Default directory for migrations
- `GOOFER_ENTITIES_DIR`: Default directory for entity definitions
- `GOOFER_DB_URL`: Default database connection URL
- `GOOFER_CONFIG`: Path to config file

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Command line parsing error
- `3`: Database error
- `4`: Migration error