# Gofer ORM - Code Walkthrough

## Overview
Gofer ORM is a lightweight, type-safe ORM (Object-Relational Mapping) library for Go. It provides a simple way to interact with SQL databases while maintaining type safety and following Go idioms.

## Core Components

### 1. Schema Package (`pkg/schema/`)
The schema package is responsible for:
- Managing entity metadata (tables, columns, relationships)
- Parsing struct tags to extract ORM configuration
- Validating entity schemas
- Maintaining a registry of all registered entities

Key types:
- `EntityMetadata`: Contains complete schema information for an entity
- `FieldMetadata`: Contains metadata about a single field/column
- `RelationMetadata`: Describes relationships between entities
- `SchemaRegistry`: Global registry of all entity metadata

### 2. Dialect Package (`pkg/dialect/`)
Handles database-specific SQL generation and differences between database engines.

Key components:
- `Dialect` interface: Defines methods for database-specific SQL generation
- Implementations for different databases (SQLite, PostgreSQL, MySQL)
- Handles:
  - Data type mapping
  - Table creation SQL
  - Identifier quoting
  - Placeholder generation

### 3. Repository Package (`pkg/repository/`)
Provides a generic repository pattern implementation for database operations.

Key features:
- Generic `Repository[T]` type for type-safe operations
- CRUD operations (Create, Read, Update, Delete)
- Transaction support
- Query building
- Relationship handling

### 4. Validation Package (`pkg/validation/`)
Integrates with the `go-playground/validator` package to provide validation for entities.

## How It Works

### 1. Defining Entities
Entities are defined as Go structs with special struct tags:

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255)"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
```

### 2. Initialization
```go
// Open database connection
db, _ := sql.Open("sqlite3", ":memory:")

// Create dialect
sqliteDialect := dialect.NewSQLiteDialect()

// Register entities
schema.Registry.RegisterEntity(User{})

// Create repository
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

### 3. Basic Operations

#### Create
```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}
err := userRepo.Save(user)
```

#### Read
```go
// Find by ID
user, err := userRepo.FindByID(1)

// Find with conditions
users, err := userRepo.Find("name = ?", "John Doe")
```

#### Update
```go
user.Name = "Jane Doe"
err := userRepo.Save(user)
```

#### Delete
```go
err := userRepo.Delete(user)
```

## Key Design Decisions

1. **Type Safety**: Uses Go generics to provide type-safe database operations.

2. **No Code Generation**: Avoids code generation in favor of runtime reflection.

3. **Database Agnostic**: Supports multiple database backends through the Dialect interface.

4. **Explicit Over Implicit**: Encourages explicit configuration through struct tags.

5. **No Magic**: Tries to be as transparent as possible about the SQL being generated.

## Extending the ORM

### Adding a New Database Dialect
1. Implement the `dialect.Dialect` interface
2. Add type mappings for your database
3. Implement any database-specific SQL generation

### Adding Custom Validators
1. Register custom validators with the validator instance
2. Use the validation tags in your entity structs

## Best Practices

1. Always register all entities before using them with repositories
2. Use transactions for operations that modify multiple entities
3. Leverage validation tags to ensure data integrity
4. Be mindful of N+1 query problems when loading relationships
5. Use the query builder for complex queries

## Common Patterns

### Transactions
```gorepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Use txRepo for all operations in the transaction
    if err := txRepo.Save(&user1); err != nil {
        return err
    }
    if err := txRepo.Save(&user2); err != nil {
        return err
    }
    return nil
})
```

### Custom Queries
```go
var users []User
err := repo.DB().Select(&users, "SELECT * FROM users WHERE age > ?", 18)
```

## Performance Considerations

1. **Reflection Overhead**: The ORM uses reflection which has some overhead. For performance-critical paths, consider using raw SQL.

2. **Query Optimization**: Always check the generated SQL and add appropriate indexes.

3. **Connection Pooling**: Configure the underlying `sql.DB` connection pool appropriately for your workload.

## Testing

1. Unit tests for each package
2. Integration tests with an in-memory SQLite database
3. Test coverage for all major features

## Future Improvements

1. Add support for more database features (e.g., JSON operations)
2. Improve relationship loading strategies
3. Add more query builder features
4. Support for database migrations
5. Better error messages and debugging tools
