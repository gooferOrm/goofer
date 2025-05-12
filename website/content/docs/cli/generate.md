# Generate Commands

The generate commands in the Goofer ORM CLI help you scaffold code for your application. They can generate entities, repositories, migrations, and more, saving you time and ensuring consistency in your codebase.

## Overview

Code generation is a powerful feature that can:

- Reduce boilerplate code
- Ensure consistency across your codebase
- Speed up development
- Enforce best practices

## Available Commands

| Command | Description |
|---------|-------------|
| `goofer generate entity <name>` | Generate a new entity |
| `goofer generate repository <entity>` | Generate a repository for an entity |
| `goofer generate migration` | Generate a migration from entity definitions |
| `goofer generate all <entity>` | Generate all artifacts for an entity |

## Generating Entities

To generate a new entity:

```bash
goofer generate entity User
```

This will create a new file `entity/user.go` with a basic entity definition:

```go
package entity

import (
	"time"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}
```

### Specifying Fields

You can specify fields for your entity using the `--fields` flag:

```bash
goofer generate entity User --fields "name:string:notnull email:string:unique,notnull age:int"
```

This will generate an entity with the specified fields:

```go
package entity

import (
	"time"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	Age       int       `orm:"type:int"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}
```

The field specification format is:

```
name:type[:options]
```

Where:
- `name` is the field name
- `type` is the Go type (string, int, bool, etc.)
- `options` are comma-separated ORM tag options (notnull, unique, etc.)

### Specifying Relationships

You can specify relationships using the `--relations` flag:

```bash
goofer generate entity Post --fields "title:string:notnull content:string:notnull user_id:uint:notnull" --relations "user:belongsTo:User:user_id"
```

This will generate an entity with the specified relationship:

```go
package entity

import (
	"time"
)

// Post entity
type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Title     string    `orm:"type:varchar(255);notnull" validate:"required"`
	Content   string    `orm:"type:text;notnull" validate:"required"`
	UserID    uint      `orm:"index;notnull" validate:"required"`
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the Post entity
func (Post) TableName() string {
	return "posts"
}
```

The relationship specification format is:

```
name:type:entity[:foreignKey]
```

Where:
- `name` is the relationship field name
- `type` is the relationship type (hasOne, hasMany, belongsTo, manyToMany)
- `entity` is the related entity name
- `foreignKey` is the foreign key field name (optional)

## Generating Repositories

To generate a repository for an entity:

```bash
goofer generate repository User
```

This will create a new file `repository/user_repository.go` with a repository implementation:

```go
package repository

import (
	"database/sql"

	"github.com/gooferOrm/goofer/pkg/dialect"
	"github.com/gooferOrm/goofer/pkg/repository"
	"your-module/entity"
)

// UserRepository provides access to User entities
type UserRepository struct {
	*repository.Repository[entity.User]
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB, dialect dialect.Dialect) *UserRepository {
	return &UserRepository{
		Repository: repository.NewRepository[entity.User](db, dialect),
	}
}

// Custom methods for UserRepository can be added here
```

### Custom Repository Methods

You can add custom methods to the generated repository:

```go
// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	return r.Find().
		Where("email = ?", email).
		One()
}

// FindActive finds all active users
func (r *UserRepository) FindActive() ([]entity.User, error) {
	return r.Find().
		Where("active = ?", true).
		All()
}
```

## Generating Migrations

To generate a migration from your entity definitions:

```bash
goofer generate migration initial_schema
```

This will analyze your registered entities and generate SQL statements to create the corresponding tables, indexes, and constraints.

The generated migration files will be placed in your migrations directory:

- `YYYYMMDDHHMMSS_initial_schema.up.sql`: Contains SQL to create tables
- `YYYYMMDDHHMMSS_initial_schema.down.sql`: Contains SQL to drop tables

## Generating All Artifacts

To generate all artifacts for an entity:

```bash
goofer generate all User
```

This will generate:
- The User entity
- A UserRepository
- A migration for the User table

It's a convenient way to quickly scaffold all the necessary code for a new entity.

## Generate Options

### Global Options

These options can be used with any generate command:

| Option | Description |
|--------|-------------|
| `--output`, `-o` | Specify the output directory (default: ./src) |
| `--verbose`, `-v` | Enable verbose output |
| `--config`, `-c` | Specify a config file (default: ./goofer.yaml) |

### Command-Specific Options

#### `generate entity`

| Option | Description |
|--------|-------------|
| `--fields`, `-f` | Specify entity fields |
| `--relations`, `-r` | Specify entity relationships |
| `--timestamps` | Add CreatedAt and UpdatedAt fields (default: true) |
| `--package`, `-p` | Specify the package name (default: entity) |

#### `generate repository`

| Option | Description |
|--------|-------------|
| `--package`, `-p` | Specify the package name (default: repository) |
| `--methods`, `-m` | Specify custom methods to generate |

#### `generate migration`

| Option | Description |
|--------|-------------|
| `--dir`, `-d` | Specify the migrations directory (default: ./migrations) |
| `--dialect` | Specify the database dialect (default: from config) |

## Templates

Goofer ORM uses templates for code generation. You can customize these templates to match your project's coding style and requirements.

### Default Templates

The default templates are embedded in the CLI, but you can override them by creating your own templates in the `.goofer/templates` directory.

### Custom Templates

To create a custom template:

1. Create a `.goofer/templates` directory in your project
2. Create a template file with the same name as the default template
3. Customize the template to your needs

For example, to customize the entity template:

```
.goofer/templates/entity.tmpl
```

### Template Variables

Templates have access to various variables depending on the template type:

#### Entity Template

- `{{.Name}}`: The entity name
- `{{.Fields}}`: The entity fields
- `{{.Package}}`: The package name
- `{{.TableName}}`: The table name

#### Repository Template

- `{{.Entity}}`: The entity name
- `{{.Package}}`: The package name
- `{{.EntityPackage}}`: The entity package name
- `{{.Methods}}`: The custom methods

## Examples

### Generating a User Entity with Fields

```bash
goofer generate entity User --fields "name:string:notnull email:string:unique,notnull age:int:notnull active:bool:notnull"
```

### Generating a Post Entity with a Relationship to User

```bash
goofer generate entity Post --fields "title:string:notnull content:string:notnull user_id:uint:notnull" --relations "user:belongsTo:User:user_id"
```

### Generating a Repository with Custom Methods

```bash
goofer generate repository User --methods "FindByEmail,FindActive"
```

### Generating a Migration for All Entities

```bash
goofer generate migration initial_schema
```

### Generating All Artifacts for a Comment Entity

```bash
goofer generate all Comment --fields "content:string:notnull post_id:uint:notnull user_id:uint:notnull" --relations "post:belongsTo:Post:post_id user:belongsTo:User:user_id"
```

## Best Practices

### Entity Naming

- Use singular, PascalCase names for entities (e.g., `User`, `Post`, `Comment`)
- Use plural, snake_case names for table names (e.g., `users`, `posts`, `comments`)

### Field Naming

- Use PascalCase for field names (e.g., `FirstName`, `EmailAddress`)
- Use snake_case for column names (automatically converted by Goofer ORM)

### Relationship Naming

- Use descriptive names for relationships (e.g., `Author` instead of `User` for a post's author)
- Use plural names for one-to-many and many-to-many relationships (e.g., `Posts`, `Comments`)

### Repository Naming

- Use the entity name followed by "Repository" (e.g., `UserRepository`, `PostRepository`)
- Group related methods together

## Next Steps

- Learn about [Migration Commands](./migration) for database schema management
- Explore [Configuration](./config) for customizing the CLI behavior
- Check out the [Entity System](../features/entity-system) feature for more details on entity definitions