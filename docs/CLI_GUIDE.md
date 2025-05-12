# Goofer ORM CLI Guide

## Overview

The Goofer CLI provides a set of commands to help you manage your database schema, generate code, and interact with your ORM. This guide walks you through the available commands and how to use them effectively.

## Installation

The CLI is included with the Goofer ORM library. After installing the library, you can use the CLI with the `goofer` command:

```bash
# Install Goofer ORM
go get github.com/gooferOrm/goofer

# Verify installation
goofer --version
```

## Available Commands

### help

Displays help information about available commands.

```bash
goofer help [command]
```

### init

Initializes a new Goofer ORM project in the current directory.

```bash
goofer init [project-name]
```

Options:
- `--dialect=<dialect>`: Specify the database dialect (sqlite, mysql, postgres)
- `--with-examples`: Include example entities in the project

### entity

Generates a new entity struct with appropriate ORM tags.

```bash
goofer entity <entity-name> [field:type...]
```

Example:
```bash
goofer entity User id:uint:primary name:string:notnull email:string:unique:notnull
```

This generates a User entity with id, name, and email fields, with appropriate tags.

### migrate

Manages database migrations.

```bash
# Create a new migration
goofer migrate create <migration-name>

# Run pending migrations
goofer migrate up

# Rollback the most recent migration
goofer migrate down

# Show migration status
goofer migrate status
```

### schema

Manages database schema based on your entities.

```bash
# Generate schema for all registered entities
goofer schema generate

# Update existing schema based on entity changes
goofer schema update

# Dump current schema to SQL
goofer schema dump > schema.sql
```

### query

Interactive query builder for testing SQL queries.

```bash
goofer query [file-with-entities]
```

This opens an interactive prompt where you can write queries and see the generated SQL.

## Working with Dialects

Goofer supports multiple database dialects. You can specify the dialect in commands using the `--dialect` flag:

```bash
goofer schema generate --dialect=postgres
```

Supported dialects:
- `sqlite`: SQLite dialect
- `mysql`: MySQL dialect
- `postgres`: PostgreSQL dialect

## Configuration File

You can create a `goofer.yaml` configuration file in your project root to avoid repeating common options:

```yaml
dialect: postgres
connection: "host=localhost port=5432 user=postgres password=postgres dbname=goofer_db"
entities_dir: "./models"
migrations_dir: "./migrations"
```

## Environment Variables

Goofer also supports configuration through environment variables:

- `GOOFER_DIALECT`: Database dialect
- `GOOFER_CONNECTION`: Database connection string
- `GOOFER_ENTITIES_DIR`: Directory containing entity models
- `GOOFER_MIGRATIONS_DIR`: Directory for migration files

## Examples

### Creating a New Project

```bash
# Initialize a new project
goofer init my-app --dialect=postgres

# Generate entities
cd my-app
goofer entity User id:uint:primary name:string email:string:unique
goofer entity Post id:uint:primary title:string content:text user_id:uint:index

# Generate schema
goofer schema generate

# Create a migration
goofer migrate create initial_schema

# Run migrations
goofer migrate up
```

### Working with an Existing Project

```bash
# View migration status
goofer migrate status

# Add a new field to an entity
# (Edit the entity file first)
goofer schema update

# Create a migration for the schema change
goofer migrate create add_user_age

# Run the migration
goofer migrate up
```

## Troubleshooting

### Common Issues

1. **Command not found**: Ensure that the Goofer CLI is properly installed and your Go bin directory is in your PATH.

2. **Database connection failure**: Check your connection string and make sure the database is running.

3. **Entity not found**: Make sure your entity is in the correct directory and follows the required interface.

4. **Migration conflicts**: If you have conflicts between migrations, you may need to manually edit the migration files.

### Getting Help

For more detailed help on a specific command:

```bash
goofer help <command>
```

For additional help, check the documentation or file an issue on the GitHub repository.