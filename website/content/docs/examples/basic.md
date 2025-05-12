# Basic CRUD Operations

This example demonstrates the fundamental CRUD (Create, Read, Update, Delete) operations using Goofer ORM with an SQLite database.

## Setup

First, let's set up our project and define a simple entity:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}
```

## Database Connection

Connect to the SQLite database:

```go
// Open SQLite database
db, err := sql.Open("sqlite3", "./basic.db")
if err != nil {
	log.Fatalf("Failed to open database: %v", err)
}
defer db.Close()

// Create dialect
sqliteDialect := dialect.NewSQLiteDialect()
```

## Entity Registration

Register the entity with Goofer's schema registry:

```go
// Register entity
if err := schema.Registry.RegisterEntity(User{}); err != nil {
	log.Fatalf("Failed to register User entity: %v", err)
}

// Get entity metadata
userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
```

## Table Creation

Generate and execute the SQL to create the table:

```go
// Create table
userSQL := sqliteDialect.CreateTableSQL(userMeta)
fmt.Println("SQL for table creation:")
fmt.Println(userSQL)

_, err = db.Exec(userSQL)
if err != nil {
	log.Fatalf("Failed to create users table: %v", err)
}
```

## Repository Creation

Create a repository for the User entity:

```go
// Create repository
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

## Create Operation

Create and save a new user:

```go
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
```

## Read Operations

### Find by ID

```go
// Find the user by ID
foundUser, err := userRepo.FindByID(user.ID)
if err != nil {
	log.Fatalf("Failed to find user: %v", err)
}

fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)
```

### Query with Conditions

```go
// Find users with conditions
users, err := userRepo.Find().
	Where("name LIKE ?", "%John%").
	OrderBy("name ASC").
	Limit(10).
	All()
if err != nil {
	log.Fatalf("Failed to query users: %v", err)
}

fmt.Printf("Found %d users matching criteria\n", len(users))
for _, u := range users {
	fmt.Printf("- %s (%s)\n", u.Name, u.Email)
}
```

### Count Users

```go
// Count users
count, err := userRepo.Find().Count()
if err != nil {
	log.Fatalf("Failed to count users: %v", err)
}

fmt.Printf("Total users in database: %d\n", count)
```

## Update Operation

Update an existing user:

```go
// Update the user
foundUser.Name = "Jane Doe"
foundUser.Email = "jane@example.com"

if err := userRepo.Save(foundUser); err != nil {
	log.Fatalf("Failed to update user: %v", err)
}

// Verify update
updatedUser, err := userRepo.FindByID(foundUser.ID)
if err != nil {
	log.Fatalf("Failed to find updated user: %v", err)
}

fmt.Printf("Updated user: %s (%s)\n", updatedUser.Name, updatedUser.Email)
```

## Delete Operation

Delete a user:

```go
// Delete the user
if err := userRepo.Delete(updatedUser); err != nil {
	log.Fatalf("Failed to delete user: %v", err)
}

// Verify deletion
_, err = userRepo.FindByID(updatedUser.ID)
if err == nil {
	log.Fatalf("User still exists after deletion")
} else {
	fmt.Println("User successfully deleted")
}
```

## Complete Example

Here's the complete code:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./basic.db")
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
	fmt.Println("SQL for table creation:")
	fmt.Println(userSQL)

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

	// Find users with conditions
	users, err := userRepo.Find().
		Where("name LIKE ?", "%John%").
		OrderBy("name ASC").
		Limit(10).
		All()
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}

	fmt.Printf("Found %d users matching criteria\n", len(users))
	for _, u := range users {
		fmt.Printf("- %s (%s)\n", u.Name, u.Email)
	}

	// Count users
	count, err := userRepo.Find().Count()
	if err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}

	fmt.Printf("Total users in database: %d\n", count)

	// Update the user
	foundUser.Name = "Jane Doe"
	foundUser.Email = "jane@example.com"

	if err := userRepo.Save(foundUser); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	updatedUser, err := userRepo.FindByID(foundUser.ID)
	if err != nil {
		log.Fatalf("Failed to find updated user: %v", err)
	}

	fmt.Printf("Updated user: %s (%s)\n", updatedUser.Name, updatedUser.Email)

	// Delete the user
	if err := userRepo.Delete(updatedUser); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = userRepo.FindByID(updatedUser.ID)
	if err == nil {
		log.Fatalf("User still exists after deletion")
	} else {
		fmt.Println("User successfully deleted")
	}
}
```

## Output

When you run this example, you should see output similar to:

```
SQL for table creation:
CREATE TABLE IF NOT EXISTS "users" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "email" TEXT NOT NULL UNIQUE,
  "created_at" TEXT DEFAULT CURRENT_TIMESTAMP
);

Created user with ID: 1
Found user: John Doe (john@example.com)
Found 1 users matching criteria
- John Doe (john@example.com)
Total users in database: 1
Updated user: Jane Doe (jane@example.com)
User successfully deleted
```

This demonstrates the basic CRUD operations using Goofer ORM with an SQLite database. The same patterns apply to other dialects like MySQL and PostgreSQL with minimal changes.