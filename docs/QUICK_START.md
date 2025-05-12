# Goofer ORM - Quick Start Guide

This guide will help you get started with Goofer ORM quickly by walking through a basic example.

## Prerequisites

- Go 1.21 or higher
- Basic understanding of Go and SQL databases

## Installation

```bash
go get github.com/gooferOrm/goofer
```

## Creating a New Project

Let's create a simple project to demonstrate Goofer ORM's capabilities.

### 1. Project Setup

Create a new directory for your project and initialize a Go module:

```bash
mkdir goofer-demo
cd goofer-demo
go mod init example.com/goofer-demo
```

### 2. Define Your Entities

Create a file called `models.go`:

```go
package main

import (
	"time"

	"github.com/gooferOrm/goofer/schema"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// Post entity
type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Title     string    `orm:"type:varchar(255);notnull" validate:"required"`
	Content   string    `orm:"type:text" validate:"required"`
	UserID    uint      `orm:"index;notnull" validate:"required"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}

// TableName returns the table name for the Post entity
func (Post) TableName() string {
	return "posts"
}
```

### 3. Create the Main Application

Create a file called `main.go`:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

func main() {
	// Open SQLite database (in memory for this example)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create SQLite dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Post{}); err != nil {
		log.Fatalf("Failed to register Post entity: %v", err)
	}

	// Get entity metadata
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
	postMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Post{}))

	// Create tables
	_, err = db.Exec(sqliteDialect.CreateTableSQL(userMeta))
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(postMeta))
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	// Create repositories
	userRepo := repository.NewRepository[User](db, sqliteDialect)
	postRepo := repository.NewRepository[Post](db, sqliteDialect)

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

	// Create a post
	post := &Post{
		Title:   "My First Post",
		Content: "Hello, Goofer ORM!",
		UserID:  user.ID,
	}

	// Save the post
	if err := postRepo.Save(post); err != nil {
		log.Fatalf("Failed to save post: %v", err)
	}

	fmt.Printf("Created post with ID: %d\n", post.ID)

	// Find the user by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)

	// Find posts by user ID
	posts, err := postRepo.Find().Where("user_id = ?", user.ID).All()
	if err != nil {
		log.Fatalf("Failed to find posts: %v", err)
	}

	fmt.Printf("Found %d posts by user %s:\n", len(posts), foundUser.Name)
	for _, p := range posts {
		fmt.Printf("- %s: %s\n", p.Title, p.Content)
	}
}
```

### 4. Install Dependencies

Update your Go module dependencies:

```bash
go get github.com/gooferOrm/goofer
go get github.com/mattn/go-sqlite3
```

### 5. Run the Application

```bash
go run *.go
```

You should see output similar to:

```
Created user with ID: 1
Created post with ID: 1
Found user: John Doe (john@example.com)
Found 1 posts by user John Doe:
- My First Post: Hello, Goofer ORM!
```

## Next Steps

Now that you have a basic application working, you can:

1. **Explore relationships**: Add more complex entity relationships (one-to-many, many-to-many)
2. **Add validation**: Implement custom validation for your entities
3. **Use hooks**: Add lifecycle hooks to automate tasks during CRUD operations
4. **Try migrations**: Evolve your database schema over time with migrations
5. **Use different databases**: Switch to MySQL or PostgreSQL for production

Check out the [examples directory](https://github.com/gooferOrm/goofer/tree/main/examples) for more advanced usage patterns and the [full documentation](https://github.com/gooferOrm/goofer/blob/main/README.md) for detailed information.

## Troubleshooting

If you encounter issues:

1. Ensure all entities are registered with `schema.Registry.RegisterEntity()`
2. Check your connection string for the database
3. Make sure entity relationships are properly defined
4. For SQLite issues, check that CGO is enabled (required by mattn/go-sqlite3)

## Community Resources

- GitHub: [github.com/gooferOrm/goofer](https://github.com/gooferOrm/goofer)
- Issues: [github.com/gooferOrm/goofer/issues](https://github.com/gooferOrm/goofer/issues)

Happy coding with Goofer ORM!