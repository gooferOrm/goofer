# Goofer ORM Features

Goofer ORM is a powerful, type-safe ORM for Go that provides an amazing developer experience. It's designed to make working with databases in Go a pleasant experience with zero drama.

## 🚀 Quick Navigation

- **[Complete Tutorial](../getting-started/complete-tutorial)** - Build a blog app from scratch
- **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** - Deep dive into all features  
- **[Client & Engine Guide](../../CLIENT_ENGINE_GUIDE)** - Simplified usage patterns
- **[Migration Guide](../../MIGRATION_GUIDE)** - Master database migrations

## Core Features

### Entity System

The [Entity System](./entity-system) is the foundation of Goofer ORM. It allows you to define your database schema using Go structs with tags for metadata. This approach provides:

- **Type Safety**: Full compile-time type checking with generics
- **Zero Magic**: Transparent SQL generation you can inspect
- **Intuitive Design**: Natural Go struct syntax
- **Flexible Configuration**: Extensive struct tag options

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
    Name      string    `orm:"type:varchar(255);notnull" validate:"required,min=2"`
    Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}
```

### Repository Pattern

The [Repository Pattern](./repository-pattern) implementation provides a type-safe API for database operations:

- **Type-Safe Operations**: Generic Repository[T] for each entity type
- **Fluent Query Builder**: Chainable methods for complex queries
- **CRUD Operations**: Create, Read, Update, Delete with hooks
- **Advanced Queries**: Filtering, sorting, pagination, aggregation

```go
// Type-safe repository operations
userRepo := repository.NewRepository[User](db, dialect)

// Fluent query building
users, err := userRepo.Find().
    Where("age > ?", 18).
    Where("name LIKE ?", "%John%").
    OrderBy("created_at DESC").
    Limit(10).
    WithRelation("Posts").
    All()
```

### Relationship Management

[Relation Mapping](./relation-mapping) in Goofer ORM makes it easy to work with related entities:

- **All Relationship Types**: One-to-One, One-to-Many, Many-to-Many
- **Eager & Lazy Loading**: Control when related data is loaded
- **Automatic Join Tables**: Managed many-to-many relationships
- **Nested Loading**: Load relations of relations

```go
// Define relationships
type User struct {
    Posts []Post `orm:"relation:OneToMany;foreignKey:UserID"`
    Roles []Role `orm:"relation:ManyToMany;joinTable:user_roles"`
}

// Eager loading
users, err := userRepo.Find().
    WithRelation("Posts").
    WithRelation("Roles").
    All()
```

### Migration System

The [Migration Engine](./migration-engine) helps you evolve your database schema over time:

- **Automatic Generation**: SQL generation from entity metadata
- **Version Control**: Track schema changes over time
- **Bidirectional**: Up and down migrations for rollbacks
- **Production Ready**: Zero-downtime migration strategies

```go
// Generate migrations from entities
generator := migration.NewMigrationGenerator(schema.Registry, dialect, "./migrations")
err := generator.Generate("create_initial_tables")

// Run migrations
migrator := migration.NewMigrator(db, dialect, "./migrations")
err := migrator.Up()
```

### Validation System

[Validation](./validation) ensures your data meets your requirements before it hits the database:

- **Struct Tag Validation**: Integration with go-playground/validator
- **Custom Validation**: Business logic validation in hooks
- **Rich Error Messages**: Detailed validation feedback
- **Pre-save Validation**: Automatic validation before database operations

```go
type User struct {
    Name     string `validate:"required,min=2,max=100"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"gte=0,lte=130"`
    Password string `validate:"required,min=8"`
}

// Custom validation in hooks
func (u *User) BeforeSave() error {
    if u.Role == "admin" && u.Age < 18 {
        return fmt.Errorf("admin users must be at least 18 years old")
    }
    return nil
}
```

### Lifecycle Hooks

[Hooks](./hooks) allow you to execute code at specific points in an entity's lifecycle:

- **Complete Lifecycle**: Before/After Create, Update, Delete, Save
- **Data Transformation**: Normalize data before saving
- **Audit Logging**: Track changes automatically
- **Business Logic**: Enforce complex business rules

```go
// Hooks for data transformation and auditing
func (u *User) BeforeSave() error {
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    return nil
}

func (u *User) AfterCreate() error {
    go sendWelcomeEmail(u.Email)
    return nil
}
```

### Database Dialects

[Dialects Support](./dialects) allows Goofer ORM to work with multiple database systems:

- **Multiple Databases**: SQLite, MySQL, PostgreSQL
- **Database-Specific Features**: Proper type mapping and SQL generation
- **Consistent API**: Same code works across databases
- **Custom Dialects**: Extend support for other databases

```go
// Database-specific SQL generation
sqliteDialect := dialect.NewSQLiteDialect()
mysqlDialect := dialect.NewMySQLDialect()
postgresDialect := dialect.NewPostgresDialect()

// Same entity, different SQL output
userSQL := dialect.CreateTableSQL(userMeta)
```

### Transaction Management

[Transactions](./transactions) ensure data integrity for complex operations:

- **First-class Support**: Built-in transaction handling
- **Automatic Rollback**: Rollback on errors
- **Repository Integration**: Use repositories within transactions
- **Nested Transactions**: Support for savepoints

```go
// Transaction with automatic rollback
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    user := &User{Name: "John", Email: "john@example.com"}
    if err := txRepo.Save(user); err != nil {
        return err // Automatically rolls back
    }
    
    profile := &Profile{UserID: user.ID, Bio: "Developer"}
    return profileRepo.WithTx(txRepo.DB()).Save(profile)
})
```

## Advanced Features

### Query Builder

Build complex queries with a fluent, type-safe API:

```go
// Complex query with multiple conditions
posts, err := postRepo.Find().
    Where("published = ?", true).
    Where("created_at > ?", lastWeek).
    WhereIn("category_id", []interface{}{1, 2, 3}).
    WhereBetween("view_count", 100, 1000).
    OrderBy("view_count DESC").
    OrderBy("created_at DESC").
    Limit(20).
    Offset(40).
    All()
```

### Schema Introspection

Analyze and understand your database schema:

```go
// Inspect entity metadata
userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
fmt.Printf("Table: %s", userMeta.TableName)
for _, field := range userMeta.Fields {
    fmt.Printf("Field: %s -> %s", field.Name, field.DBName)
}
```

### Custom Queries

Execute raw SQL when needed while maintaining type safety:

```go
// Custom query with struct mapping
type UserSummary struct {
    Name      string `db:"name"`
    PostCount int    `db:"post_count"`
}

var summaries []UserSummary
err := userRepo.DB().Select(&summaries, `
    SELECT u.name, COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.id = p.user_id
    GROUP BY u.id
    ORDER BY post_count DESC
`)
```

### Performance Features

- **Connection Pooling**: Configurable connection management
- **Query Optimization**: Efficient SQL generation
- **Eager Loading**: Prevent N+1 query problems
- **Caching Strategies**: Repository-level caching patterns
- **Batch Operations**: Efficient bulk operations

### Development Features

- **CLI Tools**: Code generation and migration management
- **Debug Logging**: Inspect generated SQL queries
- **Health Checks**: Monitor database connectivity
- **Testing Support**: In-memory databases and mocks

## Getting Started

Ready to dive in? Check out our guides:

1. **[Quickstart](../getting-started/quickstart)** - Get running in 5 minutes
2. **[Complete Tutorial](../getting-started/complete-tutorial)** - Build a real application
3. **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** - Master all features
4. **[Examples](../examples)** - Working code samples

## Philosophy

Goofer ORM is built on these principles:

- **Simplicity**: Easy to learn and use
- **Type Safety**: Leverage Go's type system
- **Transparency**: No hidden magic or surprises
- **Performance**: Efficient by default
- **Flexibility**: Extensible and customizable

Explore each feature in depth by clicking on the links above, or jump straight into our [complete tutorial](../getting-started/complete-tutorial) to see everything in action!

### Entity System

The [Entity System](./entity-system) is the foundation of Goofer ORM. It allows you to define your database schema using Go structs with tags for metadata. This approach provides:

- Type safety through Go's type system
- Compile-time checks for your database models
- Clear and concise schema definition

### Schema Parser

The [Schema Parser](./schema-parser) uses Go's reflection capabilities to analyze your entity structs at runtime. It:

- Extracts metadata from struct tags
- Maps Go types to database types
- Builds a complete schema registry for your application

### Relation Mapping

[Relation Mapping](./relation-mapping) in Goofer ORM makes it easy to work with related entities. It supports:

- One-to-One relationships
- One-to-Many relationships
- Many-to-One relationships
- Many-to-Many relationships with join tables
- Eager and lazy loading strategies

### Migration Engine

The [Migration Engine](./migration-engine) helps you evolve your database schema over time. It provides:

- Automatic SQL generation for schema changes
- Versioned migrations
- Up and down migration support
- Migration status tracking

### Repository Pattern

The [Repository Pattern](./repository-pattern) implementation provides a type-safe API for database operations:

- Generic Repository[T] for each entity type
- CRUD operations (Create, Read, Update, Delete)
- Fluent query building
- Filtering, sorting, and pagination

### Validation

[Validation](./validation) ensures your data meets your requirements before it hits the database:

- Integration with go-playground/validator
- Struct tag support for validation rules
- Custom validation hooks

### Hooks

[Hooks](./hooks) allow you to execute code at specific points in an entity's lifecycle:

- BeforeCreate, AfterCreate
- BeforeUpdate, AfterUpdate
- BeforeDelete, AfterDelete
- BeforeSave, AfterSave

### Dialects Support

[Dialects Support](./dialects) allows Goofer ORM to work with multiple database systems:

- SQLite
- MySQL
- PostgreSQL
- Custom dialect support

### Transactions

[Transactions](./transactions) ensure data integrity for complex operations:

- First-class transaction support
- Automatic rollback on error
- Nested transaction support

## Additional Features

- **Type Safety**: Fully leverages Go's type system with generics
- **Zero Drama**: Simple, intuitive API with minimal boilerplate
- **Query Builder**: Fluent API for building complex queries
- **Custom Queries**: Support for raw SQL when needed
- **No Code Generation**: Uses reflection and generics instead of code generation

Explore each feature in depth by clicking on the links above.