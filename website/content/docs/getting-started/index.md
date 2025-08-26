# Getting Started with Goofer ORM

Goofer ORM is a powerful, type-safe ORM for Go that provides an amazing developer experience with relationships, migrations, and zero drama. This guide will help you get started quickly.

## Quick Navigation

- **[Complete Tutorial](./complete-tutorial)** - Build a full blog application from scratch
- **[Quickstart Guide](./quickstart)** - Get running in 5 minutes  
- **[Advanced Setup](./advanced)** - Production-ready configuration
- **[CLI Guide](./cli)** - Using the Goofer CLI tools

## What Makes Goofer Special

Goofer ORM transforms your Go structs into powerful database interfaces without the complexity of traditional ORMs. Built with Go's type system in mind, it provides:

- **Type Safety**: Full compile-time type checking with generics
- **Zero Magic**: Transparent SQL generation you can inspect and understand  
- **Multiple Databases**: SQLite, MySQL, PostgreSQL support
- **Migrations**: Version-controlled schema changes
- **Relationships**: One-to-one, one-to-many, many-to-many with eager loading
- **Validation**: Built-in validation with struct tags
- **Hooks**: Lifecycle events for custom logic
- **Performance**: Connection pooling, query optimization, caching strategies

## Installation

First, install Goofer ORM using Go modules:

```bash
go get github.com/gooferOrm/goofer
```

## Basic Setup

To use Goofer ORM, you'll need to define your entities, register them, and create a repository to interact with your database.

### 1. Define Your Entities

Create Go structs with ORM tags to represent your database tables:

```go
// User entity
type User struct {
    ID    uint   `orm:"primaryKey;autoIncrement"`
    Name  string `orm:"type:varchar(255);notnull"`
    Email string `orm:"unique;type:varchar(255);notnull"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
    return "users"
}
```

### 2. Connect to Database

Open a connection to your database:

```go
// Open SQLite database
db, err := sql.Open("sqlite3", "./db.db")
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}
defer db.Close()
```

### 3. Create Dialect

Goofer supports multiple database dialects. Choose the one that matches your database:

```go
// Create dialect for SQLite
sqliteDialect := dialect.NewSQLiteDialect()

// For MySQL
// mysqlDialect := dialect.NewMySQLDialect()

// For PostgreSQL
// postgresDialect := dialect.NewPostgresDialect()
```

### 4. Register Entity

Register your entity with Goofer's schema registry:

```go
if err := schema.Registry.RegisterEntity(User{}); err != nil {
    log.Fatalf("Failed to register User entity: %v", err)
}
```

### 5. Create Table

Generate and execute the SQL to create your table:

```go
userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
userSQL := sqliteDialect.CreateTableSQL(userMeta)
_, err = db.Exec(userSQL)
if err != nil {
    log.Fatalf("Failed to create users table: %v", err)
}
```

### 6. Create Repository

Create a typed repository to interact with your entity:

```go
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

## Basic CRUD Operations

### Create

```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}

if err := userRepo.Save(user); err != nil {
    log.Fatalf("Failed to save user: %v", err)
}

fmt.Printf("Created user with ID: %d\n", user.ID)
```

### Read

```go
// Find by ID
foundUser, err := userRepo.FindByID(user.ID)
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Query with conditions
users, err := userRepo.Find().
    Where("name LIKE ?", "%John%").
    OrderBy("name ASC").
    Limit(10).
    All()
if err != nil {
    log.Fatalf("Failed to query users: %v", err)
}
```

### Update

```go
user.Name = "Jane Doe"
if err := userRepo.Save(user); err != nil {
    log.Fatalf("Failed to update user: %v", err)
}
```

### Delete

```go
if err := userRepo.Delete(user); err != nil {
    log.Fatalf("Failed to delete user: %v", err)
}
```

## Next Steps

- **[Complete Tutorial](./complete-tutorial)** - Build a full blog application from scratch
- **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** - Deep dive into all ORM features
- **[Migration Guide](../../MIGRATION_GUIDE)** - Master database migrations
- Learn about [entity relationships](../examples/relationships)
- Explore [validation system](../features/validation)
- Check out the [CLI tools](../cli) for automation
- See [performance tips](../reference/performance) for optimization

## Community and Support

- **GitHub**: [github.com/gooferOrm/goofer](https://github.com/gooferOrm/goofer)
- **Issues**: Report bugs and request features
- **Discussions**: Ask questions and share tips
- **Examples**: Complete working examples in the repository

Goofer ORM is designed to make working with databases in Go a pleasant experience. Enjoy building with type-safe, relationship-aware database access!