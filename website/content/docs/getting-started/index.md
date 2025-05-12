# Getting Started with Goofer ORM

Goofer ORM is a powerful, type-safe ORM for Go that provides an amazing developer experience with relationships, migrations, and zero drama. This guide will help you get started quickly.

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

- See the [quickstart](quickstart) page for a more detailed walkthrough
- Learn about [entity relationships](../examples/relationships)
- Explore [migrations](../examples/migrations) for schema evolution
- Dive into the [CLI](../cli) for automation

Goofer ORM is designed to make working with databases in Go a pleasant experience. Enjoy building with type-safe, relationship-aware database access!