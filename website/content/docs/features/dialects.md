# Dialects Support

Dialects Support in Goofer ORM allows it to work with multiple database systems. It provides a consistent API while handling the differences between database systems under the hood.

## Overview

The Dialects Support offers the following capabilities:

- Support for multiple database systems (SQLite, MySQL, PostgreSQL)
- Abstracted interfaces for different SQL dialects
- Handling of data type mapping
- Database-specific SQL generation
- Custom dialect support

## Supported Dialects

Goofer ORM currently supports the following database systems:

- **SQLite**: A lightweight, file-based database
- **MySQL**: A popular open-source relational database
- **PostgreSQL**: A powerful, open-source object-relational database

## Dialect Interface

All dialects implement the `Dialect` interface:

```go
// Dialect interface for database-specific implementations
type Dialect interface {
    // Placeholder returns the placeholder for a parameter at the given index
    Placeholder(int) string
    
    // QuoteIdentifier quotes an identifier (table name, column name)
    QuoteIdentifier(string) string
    
    // DataType maps a field metadata to a database-specific type
    DataType(field schema.FieldMetadata) string
    
    // CreateTableSQL generates SQL to create a table for the entity
    CreateTableSQL(*schema.EntityMetadata) string
    
    // Name returns the name of the dialect
    Name() string
}
```

This interface ensures that all dialects provide the necessary functionality for the ORM to work with different database systems.

## Using Dialects

To use a specific dialect, create an instance of the dialect and pass it to the repository:

```go
// SQLite dialect
sqliteDialect := &dialect.SQLiteDialect{}

// MySQL dialect
mysqlDialect := &dialect.MySQLDialect{}

// PostgreSQL dialect
postgresDialect := &dialect.PostgresDialect{}

// Create a repository with the dialect
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

## Dialect-Specific Features

### SQLite

The SQLite dialect is designed for lightweight, file-based databases:

```go
// Open SQLite database
db, err := sql.Open("sqlite3", "./database.db")
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}

// Create SQLite dialect
sqliteDialect := &dialect.SQLiteDialect{}

// Create repository with SQLite dialect
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

SQLite has some specific characteristics:

- Uses `?` as parameter placeholders
- Uses double quotes for identifiers
- Has a simpler type system than other databases
- Supports `AUTOINCREMENT` for primary keys

### MySQL

The MySQL dialect is designed for the MySQL database system:

```go
// Open MySQL database
db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}

// Create MySQL dialect
mysqlDialect := &dialect.MySQLDialect{}

// Create repository with MySQL dialect
userRepo := repository.NewRepository[User](db, mysqlDialect)
```

MySQL has some specific characteristics:

- Uses `?` as parameter placeholders
- Uses backticks for identifiers
- Has specific data types like `TINYINT(1)` for booleans
- Uses `AUTO_INCREMENT` for auto-incrementing fields
- Requires an engine specification (default is InnoDB)

### PostgreSQL

The PostgreSQL dialect is designed for the PostgreSQL database system:

```go
// Open PostgreSQL database
db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/dbname?sslmode=disable")
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}

// Create PostgreSQL dialect
postgresDialect := &dialect.PostgresDialect{}

// Create repository with PostgreSQL dialect
userRepo := repository.NewRepository[User](db, postgresDialect)
```

PostgreSQL has some specific characteristics:

- Uses `$1`, `$2`, etc. as parameter placeholders
- Uses double quotes for identifiers
- Has specific data types like `SERIAL` for auto-incrementing fields
- Supports advanced features like JSON, arrays, and custom types

## Data Type Mapping

Each dialect maps Go types to database-specific types:

### SQLite Type Mapping

| Go Type | SQLite Type |
|---------|-------------|
| string | TEXT |
| int, int8, int16, int32, int64 | INTEGER |
| uint, uint8, uint16, uint32, uint64 | INTEGER |
| float32, float64 | REAL |
| bool | INTEGER |
| time.Time | TEXT |
| []byte | BLOB |

### MySQL Type Mapping

| Go Type | MySQL Type |
|---------|------------|
| string | VARCHAR(255) |
| int, int8, int16, int32, int64 | INT |
| uint, uint8, uint16, uint32, uint64 | INT UNSIGNED |
| float32 | FLOAT |
| float64 | DOUBLE |
| bool | TINYINT(1) |
| time.Time | DATETIME |
| []byte | BLOB |

### PostgreSQL Type Mapping

| Go Type | PostgreSQL Type |
|---------|-----------------|
| string | VARCHAR(255) |
| int, int8, int16, int32 | INTEGER |
| int64 | BIGINT |
| uint, uint8, uint16, uint32 | INTEGER |
| uint64 | BIGINT |
| float32 | REAL |
| float64 | DOUBLE PRECISION |
| bool | BOOLEAN |
| time.Time | TIMESTAMP |
| []byte | BYTEA |

## SQL Generation

Each dialect generates SQL statements that are compatible with the specific database system:

### Table Creation

```go
// Get entity metadata
userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))

// Generate SQL for table creation
createTableSQL := sqliteDialect.CreateTableSQL(userMeta)
fmt.Println(createTableSQL)
```

This will generate SQL that is specific to the dialect:

#### SQLite

```sql
CREATE TABLE IF NOT EXISTS "users" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "email" TEXT NOT NULL UNIQUE,
  "created_at" TEXT DEFAULT CURRENT_TIMESTAMP
);
```

#### MySQL

```sql
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### PostgreSQL

```sql
CREATE TABLE IF NOT EXISTS "users" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) NOT NULL UNIQUE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Query Parameters

Each dialect handles query parameters differently:

```go
// SQLite and MySQL use ? placeholders
users, err := userRepo.Find().
    Where("name = ?", "John").
    All()

// PostgreSQL uses $1, $2, etc. placeholders
users, err := userRepo.Find().
    Where("name = $1", "John").
    All()
```

Goofer ORM handles this automatically based on the dialect, so you can use `?` placeholders in your queries regardless of the dialect.

## Creating a Custom Dialect

You can create a custom dialect by implementing the `Dialect` interface:

```go
// CustomDialect implements the Dialect interface
type CustomDialect struct {
    dialect.BaseDialect
}

// Name returns the name of the dialect
func (d *CustomDialect) Name() string {
    return "custom"
}

// Placeholder returns the placeholder for a parameter at the given index
func (d *CustomDialect) Placeholder(index int) string {
    return "?"
}

// QuoteIdentifier quotes an identifier with custom quotes
func (d *CustomDialect) QuoteIdentifier(name string) string {
    return fmt.Sprintf("`%s`", name)
}

// DataType maps a field metadata to a custom database-specific type
func (d *CustomDialect) DataType(field schema.FieldMetadata) string {
    // Custom type mapping logic
    return "VARCHAR(255)"
}

// CreateTableSQL generates SQL to create a table for the entity
func (d *CustomDialect) CreateTableSQL(meta *schema.EntityMetadata) string {
    // Custom table creation logic
    return ""
}
```

## Best Practices

### Use the Right Dialect for Your Database

Choose the dialect that matches your database system:

```go
// For SQLite
sqliteDialect := &dialect.SQLiteDialect{}

// For MySQL
mysqlDialect := &dialect.MySQLDialect{}

// For PostgreSQL
postgresDialect := &dialect.PostgresDialect{}
```

### Be Aware of Database-Specific Features

Each database system has its own features and limitations:

- SQLite is lightweight but has limited concurrency
- MySQL has good performance but less advanced features
- PostgreSQL has advanced features but can be more complex

### Test with Your Target Database

Always test your application with the database system you'll use in production:

```go
// Test with SQLite during development
db, _ := sql.Open("sqlite3", ":memory:")
sqliteDialect := &dialect.SQLiteDialect{}

// Test with PostgreSQL before production
db, _ := sql.Open("postgres", "postgres://user:password@localhost:5432/testdb?sslmode=disable")
postgresDialect := &dialect.PostgresDialect{}
```

### Use Explicit Types in ORM Tags

To ensure consistent behavior across dialects, use explicit types in your ORM tags:

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement;type:int"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
```

## Next Steps

- Learn about the [Migration Engine](./migration-engine) to see how dialects are used in migrations
- Explore the [Repository Pattern](./repository-pattern) to understand how dialects are used in queries
- Check out the [Examples](../examples/basic) section for examples of using different dialects