# Repository Pattern

The Repository Pattern in Goofer ORM provides a type-safe, intuitive API for database operations. It abstracts away the complexities of SQL and database interactions, allowing you to work with your entities in a more object-oriented way.

## Overview

The Repository Pattern offers the following capabilities:

- Generic Repository[T] for each entity type
- CRUD operations (Create, Read, Update, Delete)
- Fluent query building
- Filtering, sorting, and pagination
- Transaction support

## Repository Interface

Goofer ORM implements the Repository Pattern using Go's generics:

```go
// Repository provides type-safe database operations
type Repository[T schema.Entity] struct {
    db       *sql.DB
    dialect  Dialect
    metadata *schema.EntityMetadata
    ctx      context.Context
}
```

The generic type parameter `T` ensures that the repository is type-safe and can only be used with the specified entity type.

## Creating a Repository

To create a repository for an entity, use the `NewRepository` function:

```go
// Create a repository for the User entity
userRepo := repository.NewRepository[User](db, sqliteDialect)
```

This creates a repository that is specifically typed for the `User` entity, providing type-safe operations.

## Basic CRUD Operations

### Create

To create a new entity, use the `Save` method:

```go
// Create a new user
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

The `Save` method will insert a new record if the entity's primary key is zero, or update an existing record if the primary key has a value.

### Read

To read entities from the database, use the `Find`, `FindByID`, or query builder methods:

```go
// Find by ID
user, err := userRepo.FindByID(1)
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Find with conditions
users, err := userRepo.Find().
    Where("name LIKE ?", "%John%").
    OrderBy("name ASC").
    Limit(10).
    All()
if err != nil {
    log.Fatalf("Failed to find users: %v", err)
}

// Find one with conditions
user, err := userRepo.Find().
    Where("email = ?", "john@example.com").
    One()
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Count entities
count, err := userRepo.Find().
    Where("name LIKE ?", "%John%").
    Count()
if err != nil {
    log.Fatalf("Failed to count users: %v", err)
}
```

### Update

To update an entity, modify its properties and use the `Save` method:

```go
// Find the user
user, err := userRepo.FindByID(1)
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Update properties
user.Name = "Jane Doe"
user.Email = "jane@example.com"

// Save the changes
if err := userRepo.Save(user); err != nil {
    log.Fatalf("Failed to update user: %v", err)
}
```

### Delete

To delete an entity, use the `Delete` or `DeleteByID` method:

```go
// Delete by entity
if err := userRepo.Delete(user); err != nil {
    log.Fatalf("Failed to delete user: %v", err)
}

// Delete by ID
if err := userRepo.DeleteByID(1); err != nil {
    log.Fatalf("Failed to delete user: %v", err)
}
```

## Query Builder

The Repository Pattern includes a fluent query builder that allows you to construct complex queries:

```go
// Create a query builder
query := userRepo.Find().
    Where("name LIKE ?", "%John%").
    Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
    OrderBy("name ASC").
    Limit(10).
    Offset(20)

// Execute the query
users, err := query.All()
if err != nil {
    log.Fatalf("Failed to find users: %v", err)
}
```

### Where Conditions

You can add multiple `Where` conditions to filter your query:

```go
query := userRepo.Find().
    Where("name LIKE ?", "%John%").
    Where("email LIKE ?", "%@example.com").
    Where("created_at > ?", time.Now().AddDate(0, -1, 0))
```

### Ordering

You can order the results using the `OrderBy` method:

```go
query := userRepo.Find().
    OrderBy("name ASC").
    OrderBy("created_at DESC")
```

### Pagination

You can paginate the results using the `Limit` and `Offset` methods:

```go
// Get the first page (10 items per page)
page1, err := userRepo.Find().
    OrderBy("name ASC").
    Limit(10).
    Offset(0).
    All()

// Get the second page
page2, err := userRepo.Find().
    OrderBy("name ASC").
    Limit(10).
    Offset(10).
    All()
```

### Counting

You can count the number of entities that match your query:

```go
count, err := userRepo.Find().
    Where("name LIKE ?", "%John%").
    Count()
```

## Transactions

The Repository Pattern supports transactions to ensure data integrity:

```go
// Start a transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }

    // Save the user in the transaction
    if err := txRepo.Save(user); err != nil {
        return err
    }

    // Create a profile for the user
    profile := &Profile{
        UserID: user.ID,
        Bio:    "Software developer",
    }

    // Save the profile in the transaction
    profileRepo := repository.NewRepository[Profile](txRepo.DB(), txRepo.Dialect())
    if err := profileRepo.Save(profile); err != nil {
        return err
    }

    // If we return nil, the transaction will be committed
    return nil
})

if err != nil {
    log.Fatalf("Transaction failed: %v", err)
}
```

If the function returns an error, the transaction will be automatically rolled back. If it returns nil, the transaction will be committed.

## Context Support

The Repository Pattern supports context for cancellation and timeouts:

```go
// Create a context with a timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use the context with the repository
userRepoWithCtx := userRepo.WithContext(ctx)

// Execute a query with the context
users, err := userRepoWithCtx.Find().
    Where("name LIKE ?", "%John%").
    All()
```

## Hooks

The Repository Pattern integrates with the [Hooks](./hooks) system to allow you to execute code at specific points in an entity's lifecycle:

```go
// User entity with hooks
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    UpdatedAt time.Time `orm:"type:timestamp"`
}

// BeforeSave is called before the entity is saved
func (u *User) BeforeSave() error {
    u.UpdatedAt = time.Now()
    return nil
}

// AfterCreate is called after the entity is created
func (u *User) AfterCreate() error {
    fmt.Printf("User created: %s\n", u.Name)
    return nil
}
```

## Best Practices

### Repository Creation

Create repositories at application startup and reuse them:

```go
// Create repositories at startup
userRepo := repository.NewRepository[User](db, dialect)
profileRepo := repository.NewRepository[Profile](db, dialect)
postRepo := repository.NewRepository[Post](db, dialect)

// Use them throughout your application
```

### Error Handling

Always check errors returned by repository methods:

```go
user, err := userRepo.FindByID(1)
if err != nil {
    if err == sql.ErrNoRows {
        // Handle not found case
        return nil, fmt.Errorf("user not found: %w", err)
    }
    // Handle other errors
    return nil, fmt.Errorf("failed to find user: %w", err)
}
```

### Transactions

Use transactions for operations that need to be atomic:

```go
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Multiple operations that need to succeed or fail together
    return nil
})
```

### Query Optimization

Be mindful of the queries you're generating:

- Use appropriate indexes on your database tables
- Limit the number of rows returned when possible
- Use `Count()` instead of loading all entities when you only need the count
- Consider the N+1 query problem when working with relationships

## Next Steps

- Learn about [Hooks](./hooks) to understand how to add lifecycle events to your entities
- Explore [Validation](./validation) to see how to validate your entities before saving them
- Check out [Transactions](./transactions) for more details on transaction support
- See the [Examples](../examples/basic) section for more examples of using the Repository Pattern