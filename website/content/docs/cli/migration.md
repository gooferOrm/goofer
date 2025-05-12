# Migration Commands

The migration commands in the Goofer ORM CLI help you manage your database schema. They allow you to create, apply, and revert migrations, ensuring that your database schema stays in sync with your entity definitions.

## Overview

Migrations are a way to evolve your database schema over time. They are especially useful when:

- Working in a team where multiple developers need to make schema changes
- Deploying your application to different environments
- Tracking schema changes in version control
- Rolling back changes if something goes wrong

## Available Commands

| Command | Description |
|---------|-------------|
| `goofer migrate create <name>` | Create a new migration |
| `goofer migrate up` | Apply pending migrations |
| `goofer migrate down` | Revert the last migration |
| `goofer migrate status` | Show migration status |
| `goofer migrate reset` | Revert all migrations |

## Creating Migrations

### Manual Creation

To create a new migration manually:

```bash
goofer migrate create create_users_table
```

This will create two files in your migrations directory:

- `YYYYMMDDHHMMSS_create_users_table.up.sql`: Contains SQL to apply the migration
- `YYYYMMDDHHMMSS_create_users_table.down.sql`: Contains SQL to revert the migration

The timestamp (`YYYYMMDDHHMMSS`) ensures that migrations are applied in the correct order.

You can then edit these files to add your SQL statements:

```sql
-- YYYYMMDDHHMMSS_create_users_table.up.sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- YYYYMMDDHHMMSS_create_users_table.down.sql
DROP TABLE users;
```

### Automatic Generation

You can also generate migrations automatically from your entity definitions:

```bash
goofer generate migration initial_schema
```

This will analyze your registered entities and generate SQL statements to create the corresponding tables, indexes, and constraints.

## Applying Migrations

To apply pending migrations:

```bash
goofer migrate up
```

This will:

1. Check which migrations have already been applied
2. Apply any pending migrations in order
3. Record the applied migrations in the database

You can also specify the number of migrations to apply:

```bash
goofer migrate up --steps 2
```

This will apply at most 2 pending migrations.

## Reverting Migrations

To revert the last applied migration:

```bash
goofer migrate down
```

This will:

1. Find the most recently applied migration
2. Execute its down migration script
3. Remove the migration record from the database

You can also specify the number of migrations to revert:

```bash
goofer migrate down --steps 2
```

This will revert the last 2 applied migrations.

## Checking Migration Status

To see the status of your migrations:

```bash
goofer migrate status
```

This will show you:

- Which migrations have been applied and when
- Which migrations are pending
- The total number of migrations

Example output:

```
Applied migrations:
  20230101120000 - create_users_table (applied at 2023-01-01 12:00:00)
  20230101130000 - create_posts_table (applied at 2023-01-01 13:00:00)

Pending migrations:
  20230101140000 - add_user_profile_table
  20230101150000 - add_comments_table

Total: 4 migrations (2 applied, 2 pending)
```

## Resetting Migrations

To revert all applied migrations:

```bash
goofer migrate reset
```

This will revert all migrations in reverse order, effectively resetting your database to its initial state.

**Warning**: This will delete all data in the affected tables. Use with caution, especially in production environments.

## Migration Options

### Global Options

These options can be used with any migration command:

| Option | Description |
|--------|-------------|
| `--dir`, `-d` | Specify the migrations directory (default: ./migrations) |
| `--verbose`, `-v` | Enable verbose output |
| `--config`, `-c` | Specify a config file (default: ./goofer.yaml) |

### Command-Specific Options

#### `migrate create`

| Option | Description |
|--------|-------------|
| `--template`, `-t` | Specify a migration template (default: empty) |

#### `migrate up`

| Option | Description |
|--------|-------------|
| `--steps`, `-s` | Number of migrations to apply (default: all) |
| `--dry-run` | Show what would be done without actually applying migrations |

#### `migrate down`

| Option | Description |
|--------|-------------|
| `--steps`, `-s` | Number of migrations to revert (default: 1) |
| `--dry-run` | Show what would be done without actually reverting migrations |

## Migration Files

Migration files are SQL files that contain the statements to apply or revert a migration. They follow a specific naming convention:

```
<timestamp>_<name>.<direction>.sql
```

Where:
- `<timestamp>` is a timestamp in the format `YYYYMMDDHHMMSS`
- `<name>` is a descriptive name for the migration
- `<direction>` is either `up` (to apply) or `down` (to revert)

Example:

```
20230101120000_create_users_table.up.sql
20230101120000_create_users_table.down.sql
```

## Migration Table

Goofer ORM keeps track of applied migrations in a `migrations` table in your database. This table has the following schema:

```sql
CREATE TABLE migrations (
  id VARCHAR(255) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  applied_at TIMESTAMP NOT NULL,
  script TEXT NOT NULL,
  checksum VARCHAR(32) NOT NULL
);
```

The `checksum` field is used to detect if a migration file has been modified after it was applied.

## Best Practices

### Migration Naming

Use descriptive names for your migrations that clearly indicate what they do:

- `create_users_table`
- `add_email_to_users`
- `create_posts_table`
- `add_user_posts_relationship`

### Migration Content

- Each migration should be focused on a specific change
- Always provide a down migration that reverts the changes
- Test migrations before applying them to production
- Include comments in your SQL to explain complex changes

### Migration Workflow

1. Create a new migration for each schema change
2. Apply migrations to your development database
3. Test your application with the new schema
4. Commit the migration files to version control
5. Apply migrations to other environments (staging, production)

### Handling Conflicts

If multiple developers create migrations at the same time, you might end up with conflicts. To avoid this:

1. Pull the latest changes before creating a new migration
2. Use descriptive names to avoid confusion
3. Coordinate schema changes with your team
4. Consider using a feature branch workflow

## Examples

### Creating and Applying a Simple Migration

```bash
# Create a migration
goofer migrate create create_users_table

# Edit the migration files
# ...

# Apply the migration
goofer migrate up
```

### Generating a Migration from Entities

```bash
# Register your entities
# ...

# Generate a migration
goofer generate migration initial_schema

# Apply the migration
goofer migrate up
```

### Reverting a Migration

```bash
# Revert the last migration
goofer migrate down

# Check the status
goofer migrate status
```

### Applying Multiple Migrations

```bash
# Apply the next 3 migrations
goofer migrate up --steps 3
```

## Troubleshooting

### Migration Failed

If a migration fails, Goofer ORM will automatically roll back the transaction, leaving your database in a consistent state. Check the error message for details on what went wrong.

### Checksum Mismatch

If you modify a migration file after it has been applied, you'll get a checksum mismatch error. To resolve this:

1. Revert the migration: `goofer migrate down`
2. Fix the migration file
3. Apply the migration again: `goofer migrate up`

### Missing Migration Files

If you're missing migration files that are recorded in the database, you'll get an error when trying to revert them. To resolve this:

1. Restore the missing migration files
2. Try the operation again

## Next Steps

- Learn about [Generate Commands](./generate) for generating entities and repositories
- Explore [Configuration](./config) for customizing the CLI behavior
- Check out the [Migration Engine](../features/migration-engine) feature for more details on how migrations work