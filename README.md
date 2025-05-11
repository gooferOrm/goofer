# Goofer ORM

Stop hand-writing SQL like it's 1999. This Go ORM gives your structs a life of their own with relationships, migrations, and zero drama.

## Features

- **Type-safe**: Fully leverages Go's type system with generics
- **Zero drama**: Simple, intuitive API with minimal boilerplate
- **Relationships**: Support for one-to-one, one-to-many, and many-to-many relationships
- **Migrations**: Easy schema migrations to evolve your database over time
- **Multiple dialects**: Support for SQLite, MySQL, and PostgreSQL
- **Validation**: Built-in validation using struct tags
- **Transactions**: First-class support for database transactions
- **Hooks**: Lifecycle hooks for entities (BeforeSave, AfterCreate, etc.)

## Installation

```bash
go get github.com/gooferOrm/goofer
```

## Quick Start

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/pkg/dialect"
	"github.com/gooferOrm/goofer/pkg/repository"
	"github.com/gooferOrm/goofer/pkg/schema"
)

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

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entity
	if err := schema.Registry.RegisterEntity(User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}

	// Get entity metadata
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))

	// Create table
	userSQL := sqliteDialect.CreateTableSQL(userMeta)
	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create repository
	userRepo := repository.NewRepository[User](db, sqliteDialect)

	// Create a user
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Save the user
	if err := userRepo.Save(user); err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Find the user by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)
}
```

## Examples

The repository includes several examples demonstrating different features:

- **Basic**: Simple CRUD operations with SQLite
- **MySQL**: Using Goofer with MySQL
- **PostgreSQL**: Using Goofer with PostgreSQL
- **Relationships**: Working with one-to-one, one-to-many, and many-to-many relationships
- **Migrations**: Evolving your database schema over time

Check the `examples` directory for complete code samples.

## Documentation

### Entity Definition

Entities are defined as Go structs with ORM tags:

```go
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
```

### ORM Tags

- `primaryKey`: Marks the field as a primary key
- `autoIncrement`: Enables auto-incrementing for the primary key
- `type:<type>`: Specifies the database column type
- `notnull`: Makes the column not nullable
- `unique`: Adds a unique constraint
- `index`: Creates an index on the column
- `default:<value>`: Sets a default value
- `relation:<type>`: Defines a relationship (OneToOne, OneToMany, ManyToOne, ManyToMany)
- `foreignKey:<field>`: Specifies the foreign key field
- `joinTable:<table>`: Specifies the join table for many-to-many relationships
- `referenceKey:<field>`: Specifies the reference key for many-to-many relationships

### Repository API

```go
// Create a repository
repo := repository.NewRepository[User](db, dialect)

// Find by ID
user, err := repo.FindByID(1)

// Query builder
users, err := repo.Find().
    Where("name LIKE ?", "%John%").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(0).
    All()

// Count
count, err := repo.Find().Where("age > ?", 18).Count()

// Save (insert or update)
err := repo.Save(user)

// Delete
err := repo.Delete(user)

// Transaction
err := repo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Operations within transaction
    return nil
})
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
