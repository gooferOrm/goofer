# Migration Engine

The Migration Engine in Goofer ORM helps you evolve your database schema over time. It provides tools for creating, applying, and managing database migrations, ensuring that your database schema stays in sync with your entity definitions.

## Overview

The Migration Engine offers the following capabilities:

- Automatic SQL generation for schema changes
- Versioned migrations with timestamps
- Up and down migration support
- Migration status tracking
- Support for multiple database dialects

## Migration Concepts

### What is a Migration?

A migration is a set of changes to your database schema that moves it from one state to another. Migrations typically include operations like:

- Creating tables
- Adding columns
- Modifying columns
- Creating indexes
- Adding constraints
- And more

### Migration Files

Goofer ORM uses a pair of SQL files for each migration:

- `<timestamp>_<name>.up.sql`: Contains SQL statements to apply the migration
- `<timestamp>_<name>.down.sql`: Contains SQL statements to revert the migration

The timestamp ensures that migrations are applied in the correct order, and the name provides a description of what the migration does.

## Using the Migration Engine

### Creating a Migration

To create a new migration, use the `MigrationGenerator`:

```go
// Create a migration generator
generator := &migration.MigrationGenerator{
    Registry: schema.Registry,
    Dialect:  sqliteDialect,
    OutPath:  "./migrations",
}

// Generate a migration
if err := generator.Generate("create_users_table"); err != nil {
    log.Fatalf("Failed to generate migration: %v", err)
}
```

This will create two files in the `./migrations` directory:

- `20230101120000_create_users_table.up.sql`
- `20230101120000_create_users_table.down.sql`

The up migration will contain SQL to create tables based on your entity definitions, and the down migration will contain SQL to drop those tables.

### Applying Migrations

To apply pending migrations, use the `Migrator`:

```go
// Create a migrator
migrator := migration.NewMigrator(db, sqliteDialect, "./migrations")

// Apply pending migrations
if err := migrator.Up(); err != nil {
    log.Fatalf("Failed to apply migrations: %v", err)
}
```

This will:

1. Create a `migrations` table in your database if it doesn't exist
2. Check which migrations have already been applied
3. Apply any pending migrations in order
4. Record the applied migrations in the `migrations` table

### Reverting Migrations

To revert the most recent migration, use the `Down` method:

```go
// Revert the most recent migration
if err := migrator.Down(); err != nil {
    log.Fatalf("Failed to revert migration: %v", err)
}
```

This will:

1. Find the most recently applied migration
2. Execute its down migration script
3. Remove the migration record from the `migrations` table

### Checking Migration Status

To check the status of your migrations, use the `Status` method:

```go
// Get migration status
migrations, err := migrator.Status()
if err != nil {
    log.Fatalf("Failed to get migration status: %v", err)
}

fmt.Println("Applied migrations:")
for _, migration := range migrations {
    fmt.Printf("  %s - %s (applied at %s)\n", migration.ID, migration.Name, migration.AppliedAt)
}
```

This will show you which migrations have been applied and when.

## Migration Table

The Migration Engine creates and maintains a `migrations` table in your database with the following schema:

```sql
CREATE TABLE migrations (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP NOT NULL,
    script TEXT NOT NULL,
    checksum VARCHAR(32) NOT NULL
);
```

This table tracks:

- The ID of each migration (timestamp)
- The name of the migration
- When the migration was applied
- The script that was executed
- A checksum to detect if the migration file has been modified

## Generating Migrations from Entities

One of the most powerful features of the Migration Engine is its ability to generate migrations automatically from your entity definitions:

```go
// Register entities
if err := schema.Registry.RegisterEntity(User{}); err != nil {
    log.Fatalf("Failed to register User entity: %v", err)
}
if err := schema.Registry.RegisterEntity(Post{}); err != nil {
    log.Fatalf("Failed to register Post entity: %v", err)
}

// Generate migration
if err := generator.Generate("initial_schema"); err != nil {
    log.Fatalf("Failed to generate migration: %v", err)
}
```

This will generate a migration that creates tables for all registered entities, with the appropriate columns, types, constraints, and relationships.

## Best Practices

### Migration Naming

Use descriptive names for your migrations that clearly indicate what they do:

- `create_users_table`
- `add_email_to_users`
- `create_posts_table`
- `add_user_posts_relationship`

### Migration Organization

Keep your migrations organized:

- Store migrations in a dedicated directory
- Use a consistent naming convention
- Include a README or documentation explaining the migration process

### Testing Migrations

Always test migrations before applying them to production:

- Apply migrations to a test database
- Verify that the schema matches your expectations
- Test that your application works with the new schema
- Test that down migrations correctly revert changes

### Version Control

Include migrations in your version control system:

- Commit migrations along with the code changes that require them
- Never modify a migration that has been applied to any environment
- Create new migrations for additional changes

## Next Steps

- Learn about [Dialects Support](./dialects) to understand how Goofer ORM works with different database systems
- Explore the [Repository Pattern](./repository-pattern) to see how to perform CRUD operations on your entities
- Check out the [CLI](../cli) for automating migration tasks