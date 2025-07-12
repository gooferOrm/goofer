# Goofer ORM CLI Application Example

This example demonstrates how to use Goofer ORM in a command-line interface application. The CLI app manages a simple blog with users and posts.

## Features

- **User Management**: Create, list, and retrieve users
- **Post Management**: Create, list, retrieve, and search posts
- **Database Statistics**: View counts of users and posts
- **Auto-migration**: Tables are automatically created when the app starts
- **Type-safe Operations**: All database operations are type-safe using Goofer ORM

## Installation

1. Navigate to the CLI app directory:
   ```bash
   cd examples/cli_app
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build -o blog-cli
   ```

## Usage

### Initialize the Database

First, initialize the database (this creates the SQLite database file and tables):

```bash
./blog-cli init
```

### User Commands

#### Create a User
```bash
./blog-cli user create "John Doe" "john@example.com"
```

#### List All Users
```bash
./blog-cli user list
```

#### Get User by ID
```bash
./blog-cli user get 1
```

### Post Commands

#### Create a Post
```bash
./blog-cli post create 1 "My First Post" "This is the content of my first post"
```

#### List All Posts
```bash
./blog-cli post list
```

#### Get Post by ID
```bash
./blog-cli post get 1
```

#### Search Posts
```bash
./blog-cli post search "first"
```

### Statistics

View database statistics:
```bash
./blog-cli stats
```

## Complete Example Workflow

Here's a complete example of using the CLI app:

```bash
# Initialize the database
./blog-cli init

# Create some users
./blog-cli user create "Alice Johnson" "alice@example.com"
./blog-cli user create "Bob Smith" "bob@example.com"

# List users to see their IDs
./blog-cli user list

# Create posts for the users (using their IDs)
./blog-cli post create 1 "Hello World" "This is my first blog post"
./blog-cli post create 1 "Go Programming" "Learning Go is fun!"
./blog-cli post create 2 "Database Design" "Understanding ORMs and databases"

# List all posts
./blog-cli post list

# Search for posts containing "Go"
./blog-cli post search "Go"

# View statistics
./blog-cli stats
```

## Key Goofer ORM Features Demonstrated

### 1. Entity Definition
```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}
```

### 2. Repository Pattern
```go
userRepo := repository.NewRepository[User](db, sqliteDialect)
postRepo := repository.NewRepository[Post](db, sqliteDialect)
```

### 3. Auto-migration
```go
client, err = engine.NewClient(db, sqliteDialect, &User{}, &Post{})
```

### 4. Type-safe Queries
```go
// Find all users
users, err := userRepo.Find().All()

// Find by ID
user, err := userRepo.FindByID(1)

// Search with conditions
posts, err := postRepo.Find().
    Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
    All()

// Count records
userCount, err := userRepo.Find().Count()
```

### 5. Save Operations
```go
// Create new record
user := &User{Name: "John", Email: "john@example.com"}
err := userRepo.Save(user)

// Update existing record
user.Name = "Jane"
err := userRepo.Save(user)
```

## Database Schema

The CLI app creates two tables:

### Users Table
- `id` (PRIMARY KEY, AUTO_INCREMENT)
- `name` (VARCHAR(255), NOT NULL)
- `email` (VARCHAR(255), UNIQUE, NOT NULL)
- `created_at` (TIMESTAMP, DEFAULT CURRENT_TIMESTAMP)

### Posts Table
- `id` (PRIMARY KEY, AUTO_INCREMENT)
- `title` (VARCHAR(255), NOT NULL)
- `content` (TEXT)
- `user_id` (INDEX, NOT NULL, FOREIGN KEY)
- `created_at` (TIMESTAMP, DEFAULT CURRENT_TIMESTAMP)

## Extending the CLI

You can easily extend this CLI app by:

1. **Adding New Entities**: Define new structs with ORM tags
2. **Adding New Commands**: Create new cobra commands for your entities
3. **Adding Validation**: Use Goofer's validation features
4. **Adding Relationships**: Define relationships between entities
5. **Adding Transactions**: Use Goofer's transaction support

## Dependencies

- `github.com/gooferOrm/goofer`: The ORM library
- `github.com/spf13/cobra`: CLI framework
- `github.com/mattn/go-sqlite3`: SQLite driver

## File Structure

```
cli_app/
├── main.go          # Main CLI application
├── go.mod           # Go module file
├── go.sum           # Dependency checksums
└── README.md        # This file
```

The database file (`cli_app.db`) will be created automatically when you run the app. 