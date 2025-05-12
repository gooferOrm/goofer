# Goofer ORM CLI

The Goofer ORM Command Line Interface (CLI) provides a set of tools to help you manage your database schema, generate code, and perform other common tasks. It's designed to make working with Goofer ORM easier and more efficient.

## Overview

The CLI offers the following capabilities:

- Database schema management (migrations)
- Code generation
- Entity scaffolding
- Database seeding
- Configuration management

## Installation

To use the Goofer ORM CLI, you need to install it first. See the [Installation](./installation) page for detailed instructions.

## Basic Usage

The basic syntax for the CLI is:

```bash
goofer [command] [subcommand] [options]
```

For example:

```bash
# Generate a migration
goofer migrate create initial_schema

# Apply pending migrations
goofer migrate up

# Generate an entity
goofer generate entity User
```

## Getting Help

You can get help for any command by using the `--help` flag:

```bash
# Get general help
goofer --help

# Get help for a specific command
goofer migrate --help

# Get help for a specific subcommand
goofer migrate create --help
```

## Available Commands

The CLI provides several commands for different tasks:

### Migration Commands

- `migrate create`: Create a new migration
- `migrate up`: Apply pending migrations
- `migrate down`: Revert the last migration
- `migrate status`: Show migration status

See [Migration Commands](./migration) for more details.

### Generate Commands

- `generate entity`: Generate a new entity
- `generate repository`: Generate a repository for an entity
- `generate migration`: Generate a migration from entity definitions

See [Generate Commands](./generate) for more details.

### Configuration Commands

- `config init`: Initialize a new configuration file
- `config show`: Show the current configuration
- `config set`: Set a configuration value

See [Configuration](./config) for more details.

## Configuration

The CLI can be configured using a `goofer.yaml` file in your project directory. This file allows you to customize the behavior of the CLI.

See [Configuration](./config) for more details on how to configure the CLI.

## Examples

Here are some common examples of using the CLI:

### Creating and Applying Migrations

```bash
# Create a new migration
goofer migrate create create_users_table

# Apply pending migrations
goofer migrate up

# Revert the last migration
goofer migrate down

# Check migration status
goofer migrate status
```

### Generating Entities

```bash
# Generate a User entity
goofer generate entity User --fields "id:uint:primaryKey,autoIncrement name:string:notnull email:string:unique,notnull"

# Generate a Post entity with a relation to User
goofer generate entity Post --fields "id:uint:primaryKey,autoIncrement title:string:notnull content:string:notnull user_id:uint:notnull" --relations "user:belongsTo:User:user_id"
```

### Initializing a New Project

```bash
# Initialize a new project
goofer init my-project

# Initialize with a specific database
goofer init my-project --db postgres
```

## Next Steps

- Learn about [Installation](./installation) to get started with the CLI
- Explore [Commands](./commands) for a complete list of available commands
- Check out [Migration Commands](./migration) for database schema management
- See [Generate Commands](./generate) for code generation
- Read about [Configuration](./config) to customize the CLI behavior