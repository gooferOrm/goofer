# Getting Started with Goofer ORM - Complete Tutorial

This tutorial will guide you through building a complete blog application using Goofer ORM, covering every aspect from initial setup to advanced features.

## What You'll Build

By the end of this tutorial, you'll have built a fully functional blog API with:
- User authentication and management
- Blog posts with categories and tags
- Comments system
- File uploads for images
- Search functionality
- Admin dashboard
- API endpoints
- Database migrations

## Prerequisites

- Go 1.21 or later
- Basic knowledge of Go programming
- Understanding of SQL databases
- A code editor (VS Code recommended)

## Step 1: Project Setup

### Initialize the Project

```bash
mkdir goofer-blog-tutorial
cd goofer-blog-tutorial
go mod init goofer-blog-tutorial
```

### Install Dependencies

```bash
# Core Goofer ORM
go get github.com/gooferOrm/goofer

# Database drivers
go get github.com/mattn/go-sqlite3     # SQLite
go get github.com/go-sql-driver/mysql  # MySQL (optional)
go get github.com/lib/pq              # PostgreSQL (optional)

# Additional packages for the tutorial
go get github.com/gin-gonic/gin        # Web framework
go get github.com/golang-jwt/jwt/v5    # JWT authentication
go get golang.org/x/crypto/bcrypt      # Password hashing
go get github.com/joho/godotenv        # Environment variables
```

### Project Structure

Create the following directory structure:

```
goofer-blog-tutorial/
├── main.go
├── .env
├── .gitignore
├── cmd/
│   ├── migrate/
│   │   └── main.go
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── category.go
│   │   ├── tag.go
│   │   └── comment.go
│   ├── repositories/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── comment.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── blog.go
│   │   └── search.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── posts.go
│   │   └── admin.go
│   └── database/
│       └── database.go
├── migrations/
├── uploads/
└── README.md
```

## Step 2: Configuration Setup

### .env file

```bash
# Database Configuration
DATABASE_URL=sqlite://./blog.db
DB_DIALECT=sqlite

# Server Configuration
PORT=8080
HOST=localhost

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# File Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760  # 10MB
```

### .gitignore

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
goofer-blog-tutorial

# Test files
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Environment files
.env
.env.local

# Database files
*.db
*.sqlite
*.sqlite3

# Uploads
uploads/

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
```

### Configuration Package

Create `internal/config/config.go`:

```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
    
    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    JWT      JWTConfig
    Upload   UploadConfig
}

type DatabaseConfig struct {
    URL     string
    Dialect string
}

type ServerConfig struct {
    Host string
    Port int
}

type JWTConfig struct {
    Secret    string
    ExpiresIn time.Duration
}

type UploadConfig struct {
    Path        string
    MaxFileSize int64
}

func Load() (*Config, error) {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        // Don't fail if .env doesn't exist (for production)
        fmt.Println("Warning: .env file not found")
    }
    
    config := &Config{
        Database: DatabaseConfig{
            URL:     getEnv("DATABASE_URL", "sqlite://./blog.db"),
            Dialect: getEnv("DB_DIALECT", "sqlite"),
        },
        Server: ServerConfig{
            Host: getEnv("HOST", "localhost"),
            Port: getEnvAsInt("PORT", 8080),
        },
        JWT: JWTConfig{
            Secret:    getEnv("JWT_SECRET", "change-this-secret"),
            ExpiresIn: getEnvAsDuration("JWT_EXPIRES_IN", "24h"),
        },
        Upload: UploadConfig{
            Path:        getEnv("UPLOAD_PATH", "./uploads"),
            MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 10485760), // 10MB
        },
    }
    
    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
    valueStr := getEnv(key, defaultValue)
    if duration, err := time.ParseDuration(valueStr); err == nil {
        return duration
    }
    // Parse default value
    duration, _ := time.ParseDuration(defaultValue)
    return duration
}
```

## Step 3: Database Models

### User Model

Create `internal/models/user.go`:

```go
package models

import (
    "time"
    "golang.org/x/crypto/bcrypt"
    "strings"
)

type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" json:"id"`
    Username  string    `orm:"unique;type:varchar(50);notnull" json:"username" validate:"required,min=3,max=50"`
    Email     string    `orm:"unique;type:varchar(255);notnull" json:"email" validate:"required,email"`
    Password  string    `orm:"type:varchar(255);notnull" json:"-" validate:"required,min=8"`
    FirstName string    `orm:"type:varchar(100);notnull" json:"first_name" validate:"required,max=100"`
    LastName  string    `orm:"type:varchar(100);notnull" json:"last_name" validate:"required,max=100"`
    Bio       string    `orm:"type:text" json:"bio"`
    Avatar    string    `orm:"type:varchar(255)" json:"avatar"`
    Role      string    `orm:"type:varchar(20);default:user" json:"role" validate:"oneof=admin editor user"`
    IsActive  bool      `orm:"type:boolean;default:true" json:"is_active"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `orm:"type:timestamp" json:"updated_at"`
    
    // Relationships
    Posts    []Post    `orm:"relation:OneToMany;foreignKey:AuthorID" json:"posts,omitempty"`
    Comments []Comment `orm:"relation:OneToMany;foreignKey:AuthorID" json:"comments,omitempty"`
}

func (User) TableName() string {
    return "users"
}

// BeforeSave hook for password hashing and data normalization
func (u *User) BeforeSave() error {
    // Normalize email
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    u.Username = strings.ToLower(strings.TrimSpace(u.Username))
    
    // Hash password if it's new or changed
    if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    
    return nil
}

// BeforeUpdate hook
func (u *User) BeforeUpdate() error {
    u.UpdatedAt = time.Now()
    return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}

// IsEditor checks if user has editor or admin role
func (u *User) IsEditor() bool {
    return u.Role == "editor" || u.Role == "admin"
}

// FullName returns the user's full name
func (u *User) FullName() string {
    return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
```

### Category Model

Create `internal/models/category.go`:

```go
package models

import "time"

type Category struct {
    ID          uint      `orm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `orm:"unique;type:varchar(100);notnull" json:"name" validate:"required,max=100"`
    Slug        string    `orm:"unique;type:varchar(100);notnull" json:"slug"`
    Description string    `orm:"type:text" json:"description"`
    Color       string    `orm:"type:varchar(7);default:#3B82F6" json:"color"` // Hex color
    IsActive    bool      `orm:"type:boolean;default:true" json:"is_active"`
    CreatedAt   time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt   time.Time `orm:"type:timestamp" json:"updated_at"`
    
    // Relationships
    Posts []Post `orm:"relation:OneToMany;foreignKey:CategoryID" json:"posts,omitempty"`
}

func (Category) TableName() string {
    return "categories"
}

// BeforeSave hook to generate slug
func (c *Category) BeforeSave() error {
    if c.Slug == "" {
        c.Slug = generateSlug(c.Name)
    }
    return nil
}

// BeforeUpdate hook
func (c *Category) BeforeUpdate() error {
    c.UpdatedAt = time.Now()
    return nil
}
```

### Tag Model

Create `internal/models/tag.go`:

```go
package models

import "time"

type Tag struct {
    ID        uint      `orm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `orm:"unique;type:varchar(50);notnull" json:"name" validate:"required,max=50"`
    Slug      string    `orm:"unique;type:varchar(50);notnull" json:"slug"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
    
    // Relationships
    Posts []Post `orm:"relation:ManyToMany;joinTable:post_tags;foreignKey:TagID;referenceKey:PostID" json:"posts,omitempty"`
}

func (Tag) TableName() string {
    return "tags"
}

// BeforeSave hook to generate slug
func (t *Tag) BeforeSave() error {
    if t.Slug == "" {
        t.Slug = generateSlug(t.Name)
    }
    return nil
}

// PostTag represents the join table for many-to-many relationship
type PostTag struct {
    PostID uint `orm:"primaryKey" json:"post_id"`
    TagID  uint `orm:"primaryKey" json:"tag_id"`
    Post   *Post `orm:"relation:ManyToOne;foreignKey:PostID" json:"post,omitempty"`
    Tag    *Tag  `orm:"relation:ManyToOne;foreignKey:TagID" json:"tag,omitempty"`
}

func (PostTag) TableName() string {
    return "post_tags"
}
```

### Post Model

Create `internal/models/post.go`:

```go
package models

import (
    "time"
    "strings"
)

type PostStatus string

const (
    PostStatusDraft     PostStatus = "draft"
    PostStatusPublished PostStatus = "published"
    PostStatusArchived  PostStatus = "archived"
)

type Post struct {
    ID          uint       `orm:"primaryKey;autoIncrement" json:"id"`
    Title       string     `orm:"type:varchar(255);notnull" json:"title" validate:"required,max=255"`
    Slug        string     `orm:"unique;type:varchar(255);notnull" json:"slug"`
    Excerpt     string     `orm:"type:text" json:"excerpt"`
    Content     string     `orm:"type:text;notnull" json:"content" validate:"required"`
    FeaturedImage string   `orm:"type:varchar(255)" json:"featured_image"`
    Status      PostStatus `orm:"type:varchar(20);default:draft" json:"status" validate:"oneof=draft published archived"`
    ViewCount   uint       `orm:"type:int;default:0" json:"view_count"`
    PublishedAt *time.Time `orm:"type:timestamp" json:"published_at"`
    CreatedAt   time.Time  `orm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt   time.Time  `orm:"type:timestamp" json:"updated_at"`
    
    // Foreign Keys
    AuthorID   uint `orm:"index;notnull" json:"author_id"`
    CategoryID uint `orm:"index;notnull" json:"category_id"`
    
    // Relationships
    Author   *User     `orm:"relation:ManyToOne;foreignKey:AuthorID" json:"author,omitempty"`
    Category *Category `orm:"relation:ManyToOne;foreignKey:CategoryID" json:"category,omitempty"`
    Tags     []Tag     `orm:"relation:ManyToMany;joinTable:post_tags;foreignKey:PostID;referenceKey:TagID" json:"tags,omitempty"`
    Comments []Comment `orm:"relation:OneToMany;foreignKey:PostID" json:"comments,omitempty"`
}

func (Post) TableName() string {
    return "posts"
}

// BeforeSave hook
func (p *Post) BeforeSave() error {
    // Generate slug if not provided
    if p.Slug == "" {
        p.Slug = generateSlug(p.Title)
    }
    
    // Generate excerpt if not provided
    if p.Excerpt == "" {
        p.Excerpt = generateExcerpt(p.Content, 160)
    }
    
    // Set published_at when status changes to published
    if p.Status == PostStatusPublished && p.PublishedAt == nil {
        now := time.Now()
        p.PublishedAt = &now
    }
    
    return nil
}

// BeforeUpdate hook
func (p *Post) BeforeUpdate() error {
    p.UpdatedAt = time.Now()
    return nil
}

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
    return p.Status == PostStatusPublished
}

// CanBeEditedBy checks if a user can edit this post
func (p *Post) CanBeEditedBy(user *User) bool {
    if user.IsAdmin() {
        return true
    }
    if user.IsEditor() && p.AuthorID == user.ID {
        return true
    }
    return false
}

// IncrementViews increments the view count
func (p *Post) IncrementViews() {
    p.ViewCount++
}
```

### Comment Model

Create `internal/models/comment.go`:

```go
package models

import "time"

type CommentStatus string

const (
    CommentStatusPending  CommentStatus = "pending"
    CommentStatusApproved CommentStatus = "approved"
    CommentStatusSpam     CommentStatus = "spam"
)

type Comment struct {
    ID        uint          `orm:"primaryKey;autoIncrement" json:"id"`
    Content   string        `orm:"type:text;notnull" json:"content" validate:"required,max=1000"`
    Status    CommentStatus `orm:"type:varchar(20);default:pending" json:"status"`
    CreatedAt time.Time     `orm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time     `orm:"type:timestamp" json:"updated_at"`
    
    // Foreign Keys
    PostID   uint `orm:"index;notnull" json:"post_id"`
    AuthorID uint `orm:"index;notnull" json:"author_id"`
    ParentID *uint `orm:"index" json:"parent_id"` // For replies
    
    // Relationships
    Post     *Post      `orm:"relation:ManyToOne;foreignKey:PostID" json:"post,omitempty"`
    Author   *User      `orm:"relation:ManyToOne;foreignKey:AuthorID" json:"author,omitempty"`
    Parent   *Comment   `orm:"relation:ManyToOne;foreignKey:ParentID" json:"parent,omitempty"`
    Replies  []Comment  `orm:"relation:OneToMany;foreignKey:ParentID" json:"replies,omitempty"`
}

func (Comment) TableName() string {
    return "comments"
}

// BeforeUpdate hook
func (c *Comment) BeforeUpdate() error {
    c.UpdatedAt = time.Now()
    return nil
}

// IsApproved checks if the comment is approved
func (c *Comment) IsApproved() bool {
    return c.Status == CommentStatusApproved
}

// IsReply checks if this is a reply to another comment
func (c *Comment) IsReply() bool {
    return c.ParentID != nil
}
```

### Utility Functions

Add these utility functions to `internal/models/utils.go`:

```go
package models

import (
    "regexp"
    "strings"
    "unicode"
)

// generateSlug creates a URL-friendly slug from a string
func generateSlug(text string) string {
    // Convert to lowercase
    slug := strings.ToLower(text)
    
    // Replace spaces and non-alphanumeric characters with hyphens
    reg := regexp.MustCompile(`[^a-z0-9]+`)
    slug = reg.ReplaceAllString(slug, "-")
    
    // Remove leading and trailing hyphens
    slug = strings.Trim(slug, "-")
    
    return slug
}

// generateExcerpt creates an excerpt from content
func generateExcerpt(content string, maxLength int) string {
    // Strip HTML tags (basic implementation)
    reg := regexp.MustCompile(`<[^>]*>`)
    plainText := reg.ReplaceAllString(content, "")
    
    // Remove extra whitespace
    plainText = strings.TrimSpace(plainText)
    
    // Truncate to max length
    if len(plainText) <= maxLength {
        return plainText
    }
    
    // Find the last space before maxLength to avoid cutting words
    truncated := plainText[:maxLength]
    lastSpace := strings.LastIndex(truncated, " ")
    
    if lastSpace > 0 {
        truncated = truncated[:lastSpace]
    }
    
    return truncated + "..."
}

// stripTags removes HTML tags from a string
func stripTags(content string) string {
    reg := regexp.MustCompile(`<[^>]*>`)
    return reg.ReplaceAllString(content, "")
}

// truncateText truncates text to a specific length
func truncateText(text string, maxLength int) string {
    if len(text) <= maxLength {
        return text
    }
    
    // Find the last space to avoid cutting words
    truncated := text[:maxLength]
    lastSpace := strings.LastIndex(truncated, " ")
    
    if lastSpace > 0 {
        truncated = truncated[:lastSpace]
    }
    
    return truncated + "..."
}

// isValidSlug checks if a string is a valid slug
func isValidSlug(slug string) bool {
    // Slug should only contain lowercase letters, numbers, and hyphens
    // Should not start or end with hyphens
    reg := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
    return reg.MatchString(slug)
}
```

## Step 4: Database Connection

Create `internal/database/database.go`:

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    
    "github.com/gooferOrm/goofer/dialect"
    "github.com/gooferOrm/goofer/repository"
    "github.com/gooferOrm/goofer/schema"
    
    "goofer-blog-tutorial/internal/config"
    "goofer-blog-tutorial/internal/models"
)

type Database struct {
    DB      *sql.DB
    Dialect dialect.Dialect
}

// New creates a new database connection and registers entities
func New(cfg *config.Config) (*Database, error) {
    // Parse database URL
    driver, dataSource := parseDatabaseURL(cfg.Database.URL)
    
    // Open database connection
    db, err := sql.Open(driver, dataSource)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Create dialect
    var dialectInstance dialect.Dialect
    switch cfg.Database.Dialect {
    case "sqlite":
        dialectInstance = dialect.NewSQLiteDialect()
    case "mysql":
        dialectInstance = dialect.NewMySQLDialect()
    case "postgres":
        dialectInstance = dialect.NewPostgresDialect()
    default:
        return nil, fmt.Errorf("unsupported database dialect: %s", cfg.Database.Dialect)
    }
    
    database := &Database{
        DB:      db,
        Dialect: dialectInstance,
    }
    
    // Register entities
    if err := database.registerEntities(); err != nil {
        return nil, fmt.Errorf("failed to register entities: %w", err)
    }
    
    return database, nil
}

// registerEntities registers all models with the schema registry
func (d *Database) registerEntities() error {
    entities := []interface{}{
        models.User{},
        models.Category{},
        models.Tag{},
        models.Post{},
        models.PostTag{},
        models.Comment{},
    }
    
    for _, entity := range entities {
        if err := schema.Registry.RegisterEntity(entity); err != nil {
            return fmt.Errorf("failed to register entity %T: %w", entity, err)
        }
    }
    
    log.Printf("Registered %d entities", len(entities))
    return nil
}

// CreateTables creates all tables based on registered entities
func (d *Database) CreateTables() error {
    for _, entityType := range []interface{}{
        models.User{},
        models.Category{},
        models.Tag{},
        models.Post{},
        models.PostTag{},
        models.Comment{},
    } {
        if err := d.createTableForEntity(entityType); err != nil {
            return err
        }
    }
    
    return nil
}

// createTableForEntity creates a table for a specific entity
func (d *Database) createTableForEntity(entity interface{}) error {
    metadata, ok := schema.Registry.GetEntityMetadata(schema.GetEntityType(entity))
    if !ok {
        return fmt.Errorf("entity metadata not found for %T", entity)
    }
    
    createSQL := d.Dialect.CreateTableSQL(metadata)
    
    log.Printf("Creating table for %T", entity)
    log.Printf("SQL: %s", createSQL)
    
    _, err := d.DB.Exec(createSQL)
    if err != nil {
        return fmt.Errorf("failed to create table for %T: %w", entity, err)
    }
    
    return nil
}

// NewUserRepository creates a new user repository
func (d *Database) NewUserRepository() *repository.Repository[models.User] {
    return repository.NewRepository[models.User](d.DB, d.Dialect)
}

// NewPostRepository creates a new post repository
func (d *Database) NewPostRepository() *repository.Repository[models.Post] {
    return repository.NewRepository[models.Post](d.DB, d.Dialect)
}

// NewCategoryRepository creates a new category repository
func (d *Database) NewCategoryRepository() *repository.Repository[models.Category] {
    return repository.NewRepository[models.Category](d.DB, d.Dialect)
}

// NewTagRepository creates a new tag repository
func (d *Database) NewTagRepository() *repository.Repository[models.Tag] {
    return repository.NewRepository[models.Tag](d.DB, d.Dialect)
}

// NewCommentRepository creates a new comment repository
func (d *Database) NewCommentRepository() *repository.Repository[models.Comment] {
    return repository.NewRepository[models.Comment](d.DB, d.Dialect)
}

// Close closes the database connection
func (d *Database) Close() error {
    return d.DB.Close()
}

// parseDatabaseURL parses a database URL and returns driver and data source
func parseDatabaseURL(url string) (string, string) {
    if strings.HasPrefix(url, "sqlite://") {
        return "sqlite3", strings.TrimPrefix(url, "sqlite://")
    } else if strings.HasPrefix(url, "mysql://") {
        return "mysql", strings.TrimPrefix(url, "mysql://")
    } else if strings.HasPrefix(url, "postgres://") {
        return "postgres", url
    }
    
    // Default to SQLite for simple paths
    return "sqlite3", url
}

// SeedData seeds the database with initial data
func (d *Database) SeedData() error {
    log.Println("Seeding database with initial data...")
    
    // Create default categories
    categoryRepo := d.NewCategoryRepository()
    categories := []models.Category{
        {Name: "Technology", Description: "Posts about technology and programming", Color: "#3B82F6"},
        {Name: "Lifestyle", Description: "Lifestyle and personal posts", Color: "#10B981"},
        {Name: "Business", Description: "Business and entrepreneurship", Color: "#F59E0B"},
        {Name: "Travel", Description: "Travel experiences and guides", Color: "#EF4444"},
    }
    
    for _, category := range categories {
        if err := categoryRepo.Save(&category); err != nil {
            log.Printf("Error seeding category %s: %v", category.Name, err)
        }
    }
    
    // Create default tags
    tagRepo := d.NewTagRepository()
    tags := []models.Tag{
        {Name: "Go"},
        {Name: "Programming"},
        {Name: "Tutorial"},
        {Name: "Web Development"},
        {Name: "Database"},
        {Name: "API"},
    }
    
    for _, tag := range tags {
        if err := tagRepo.Save(&tag); err != nil {
            log.Printf("Error seeding tag %s: %v", tag.Name, err)
        }
    }
    
    // Create admin user
    userRepo := d.NewUserRepository()
    adminUser := models.User{
        Username:  "admin",
        Email:     "admin@example.com",
        Password:  "admin123", // Will be hashed by BeforeSave hook
        FirstName: "Admin",
        LastName:  "User",
        Role:      "admin",
        Bio:       "System administrator",
    }
    
    if err := userRepo.Save(&adminUser); err != nil {
        log.Printf("Error creating admin user: %v", err)
    } else {
        log.Printf("Created admin user with ID: %d", adminUser.ID)
    }
    
    log.Println("Database seeding completed")
    return nil
}
```

This is the foundation of our blog application. In the next parts, we'll continue with:

1. Authentication service
2. Blog service with CRUD operations
3. HTTP handlers and API endpoints
4. Migration system
5. Search functionality
6. File upload handling
7. Admin dashboard

Would you like me to continue with the next parts of the tutorial?