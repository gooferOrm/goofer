# Using Goofer ORM in CLI Applications

This guide explains how to use Goofer ORM in command-line interface (CLI) applications, from simple scripts to complex multi-command tools.

## Table of Contents

1. [Basic Setup](#basic-setup)
2. [Simple CLI Example](#simple-cli-example)
3. [Advanced CLI with Cobra](#advanced-cli-with-cobra)
4. [Using Goofer CLI Commands](#using-goofer-cli-commands)
5. [Best Practices](#best-practices)
6. [Common Patterns](#common-patterns)

## Basic Setup

### 1. Initialize Your Project

First, create a new Go module and add Goofer ORM as a dependency:

```bash
mkdir my-cli-app
cd my-cli-app
go mod init my-cli-app
go get github.com/gooferOrm/goofer
```

### 2. Choose Your Database Driver

Depending on your database choice, add the appropriate driver:

```bash
# For SQLite
go get github.com/mattn/go-sqlite3

# For PostgreSQL
go get github.com/lib/pq

# For MySQL
go get github.com/go-sql-driver/mysql
```

### 3. Basic Project Structure

```
my-cli-app/
├── main.go
├── go.mod
├── go.sum
├── models/
│   └── user.go
├── commands/
│   └── user_commands.go
└── config/
    └── database.go
```

## Simple CLI Example

Here's a minimal CLI application that demonstrates core Goofer ORM concepts:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/engine"
	"github.com/gooferOrm/goofer/repository"
)

// User entity
type User struct {
	ID    uint   `orm:"primaryKey;autoIncrement"`
	Name  string `orm:"type:varchar(255);notnull"`
	Email string `orm:"unique;type:varchar(255);notnull"`
}

func (User) TableName() string {
	return "users"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command> [args...]")
		fmt.Println("Commands: create, list, get")
		os.Exit(1)
	}

	// Initialize database
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup Goofer ORM
	dialect := dialect.NewSQLiteDialect()
	client, err := engine.NewClient(db, dialect, &User{})
	if err != nil {
		log.Fatal(err)
	}

	// Create repository
	userRepo := repository.NewRepository[User](db, dialect)

	// Handle commands
	command := os.Args[1]
	switch command {
	case "create":
		if len(os.Args) != 4 {
			fmt.Println("Usage: create <name> <email>")
			os.Exit(1)
		}
		createUser(userRepo, os.Args[2], os.Args[3])
	case "list":
		listUsers(userRepo)
	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: get <id>")
			os.Exit(1)
		}
		getUser(userRepo, os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func createUser(repo *repository.Repository[User], name, email string) {
	user := &User{Name: name, Email: email}
	if err := repo.Save(user); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)
}

func listUsers(repo *repository.Repository[User]) {
	users, err := repo.Find().All()
	if err != nil {
		log.Fatal(err)
	}
	
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}

func getUser(repo *repository.Repository[User], idStr string) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Fatal("Invalid ID")
	}
	
	user, err := repo.FindByID(uint(id))
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}
```

## Advanced CLI with Cobra

For more complex CLI applications, use the Cobra library for better command structure:

### 1. Install Cobra

```bash
go get github.com/spf13/cobra
```

### 2. Create a Structured CLI

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/engine"
	"github.com/gooferOrm/goofer/repository"
)

// Global variables
var (
	db       *sql.DB
	client   *engine.Client
	userRepo *repository.Repository[User]
)

// User entity
type User struct {
	ID    uint   `orm:"primaryKey;autoIncrement"`
	Name  string `orm:"type:varchar(255);notnull"`
	Email string `orm:"unique;type:varchar(255);notnull"`
}

func (User) TableName() string {
	return "users"
}

var rootCmd = &cobra.Command{
	Use:   "user-manager",
	Short: "A CLI tool for managing users with Goofer ORM",
}

var createCmd = &cobra.Command{
	Use:   "create [name] [email]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		createUser(args[0], args[1])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		listUsers()
	},
}

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getUser(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
}

func main() {
	// Initialize database
	initDatabase()
	defer db.Close()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	dialect := dialect.NewSQLiteDialect()
	client, err = engine.NewClient(db, dialect, &User{})
	if err != nil {
		log.Fatal(err)
	}

	userRepo = repository.NewRepository[User](db, dialect)
}

func createUser(name, email string) {
	user := &User{Name: name, Email: email}
	if err := userRepo.Save(user); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)
}

func listUsers() {
	users, err := userRepo.Find().All()
	if err != nil {
		log.Fatal(err)
	}
	
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}

func getUser(idStr string) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Fatal("Invalid ID")
	}
	
	user, err := userRepo.FindByID(uint(id))
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}
```

## Using Goofer CLI Commands

Goofer ORM provides its own CLI tool for scaffolding and managing projects:

### 1. Install Goofer CLI

```bash
go install github.com/gooferOrm/goofer/cmd/goofer@latest
```

### 2. Initialize a New Project

```bash
# Initialize a new project with SQLite
goofer init my-app --dialect=sqlite --with-examples

# Initialize with PostgreSQL
goofer init my-app --dialect=postgres --with-examples

# Initialize with MySQL
goofer init my-app --dialect=mysql --with-examples
```

### 3. Generate Code

```bash
# Generate an entity
goofer generate entity User name:string:notnull email:string:unique,notnull --with-hooks

# Generate with validation
goofer generate entity Product name:string:notnull price:float64:notnull --with-validate
```

### 4. Run Migrations

```bash
# Create a migration
goofer migrate create add_users_table

# Run migrations
goofer migrate up

# Rollback migrations
goofer migrate down
```

## Best Practices

### 1. Configuration Management

Create a configuration structure for your CLI app:

```go
type Config struct {
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func loadConfig() (*Config, error) {
	// Load from file, environment variables, or flags
}
```

### 2. Error Handling

Implement proper error handling in your CLI:

```go
func handleError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", message, err)
		os.Exit(1)
	}
}

func createUser(name, email string) {
	user := &User{Name: name, Email: email}
	if err := userRepo.Save(user); err != nil {
		handleError(err, "Failed to create user")
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)
}
```

### 3. Input Validation

Validate user input before processing:

```go
func validateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func createUser(name, email string) {
	if err := validateEmail(email); err != nil {
		handleError(err, "Invalid email")
	}
	
	user := &User{Name: name, Email: email}
	if err := userRepo.Save(user); err != nil {
		handleError(err, "Failed to create user")
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)
}
```

### 4. Database Connection Management

Implement proper database connection management:

```go
type App struct {
	DB       *sql.DB
	Client   *engine.Client
	UserRepo *repository.Repository[User]
}

func NewApp(config *Config) (*App, error) {
	db, err := sql.Open(config.Database.Driver, config.Database.DSN)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	dialect := getDialect(config.Database.Driver)
	client, err := engine.NewClient(db, dialect, &User{})
	if err != nil {
		return nil, err
	}

	return &App{
		DB:       db,
		Client:   client,
		UserRepo: repository.NewRepository[User](db, dialect),
	}, nil
}

func (app *App) Close() error {
	return app.DB.Close()
}
```

## Common Patterns

### 1. Interactive CLI

For interactive command-line applications:

```go
func interactiveMode() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case "create":
			if len(args) >= 2 {
				createUser(args[0], args[1])
			} else {
				fmt.Println("Usage: create <name> <email>")
			}
		case "list":
			listUsers()
		case "quit":
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}
```

### 2. Batch Operations

For processing multiple records:

```go
func batchCreateUsers(users []User) error {
	return userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
		for _, user := range users {
			if err := txRepo.Save(&user); err != nil {
				return err
			}
		}
		return nil
	})
}
```

### 3. Search and Filter

Implement search functionality:

```go
func searchUsers(query string) {
	users, err := userRepo.Find().
		Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").
		All()
	if err != nil {
		handleError(err, "Search failed")
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}
```

### 4. Export/Import

For data export and import:

```go
func exportUsers(filename string) {
	users, err := userRepo.Find().All()
	if err != nil {
		handleError(err, "Failed to fetch users")
	}

	file, err := os.Create(filename)
	if err != nil {
		handleError(err, "Failed to create file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, user := range users {
		if err := encoder.Encode(user); err != nil {
			handleError(err, "Failed to encode user")
		}
	}

	fmt.Printf("Exported %d users to %s\n", len(users), filename)
}
```

## Conclusion

Goofer ORM provides a powerful and type-safe way to work with databases in CLI applications. By following these patterns and best practices, you can create robust, maintainable command-line tools that leverage the full power of Goofer ORM's features.

For more examples, check out the `examples/` directory in the Goofer ORM repository, which includes both simple and advanced CLI application examples. 