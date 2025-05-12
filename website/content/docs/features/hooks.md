# Hooks

Hooks in Goofer ORM allow you to execute code at specific points in an entity's lifecycle. They provide a way to implement cross-cutting concerns like validation, logging, and automatic field updates without cluttering your business logic.

## Overview

The Hooks system offers the following capabilities:

- Lifecycle hooks for entities (BeforeSave, AfterCreate, etc.)
- Interface-based hook registration
- Global hooks that apply to all entities
- Entity-specific hooks
- Transaction-aware hooks

## Lifecycle Hooks

Goofer ORM supports the following lifecycle hooks:

| Hook | Triggered | Use Case |
|------|-----------|----------|
| `BeforeSave` | Before an entity is saved (insert or update) | Validation, setting timestamps |
| `AfterSave` | After an entity is saved | Logging, cache invalidation |
| `BeforeCreate` | Before a new entity is created | Setting default values |
| `AfterCreate` | After a new entity is created | Logging, sending notifications |
| `BeforeUpdate` | Before an existing entity is updated | Validation, setting update timestamps |
| `AfterUpdate` | After an existing entity is updated | Logging, cache invalidation |
| `BeforeDelete` | Before an entity is deleted | Validation, archiving |
| `AfterDelete` | After an entity is deleted | Logging, cleanup related data |

## Implementing Hooks

There are two ways to implement hooks in Goofer ORM:

### 1. Entity Methods

You can add hook methods directly to your entity structs:

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp"`
    UpdatedAt time.Time `orm:"type:timestamp"`
}

// BeforeSave is called before the entity is saved
func (u *User) BeforeSave() error {
    // Set timestamps
    if u.ID == 0 {
        u.CreatedAt = time.Now()
    }
    u.UpdatedAt = time.Now()
    return nil
}

// AfterCreate is called after the entity is created
func (u *User) AfterCreate() error {
    fmt.Printf("User created: %s (%s)\n", u.Name, u.Email)
    return nil
}

// BeforeDelete is called before the entity is deleted
func (u *User) BeforeDelete() error {
    fmt.Printf("About to delete user: %s\n", u.Name)
    return nil
}
```

### 2. Hook Interfaces

You can also implement hook interfaces for more flexibility:

```go
// Hook interfaces for entity lifecycle events
type (
    BeforeCreateHook interface {
        BeforeCreate() error
    }

    AfterCreateHook interface {
        AfterCreate() error
    }

    BeforeUpdateHook interface {
        BeforeUpdate() error
    }

    AfterUpdateHook interface {
        AfterUpdate() error
    }

    BeforeDeleteHook interface {
        BeforeDelete() error
    }

    AfterDeleteHook interface {
        AfterDelete() error
    }

    BeforeSaveHook interface {
        BeforeSave() error
    }

    AfterSaveHook interface {
        AfterSave() error
    }
)
```

Your entity can implement any of these interfaces:

```go
// Implement specific hook interfaces
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp"`
    UpdatedAt time.Time `orm:"type:timestamp"`
}

// BeforeSave implements the BeforeSaveHook interface
func (u *User) BeforeSave() error {
    u.UpdatedAt = time.Now()
    return nil
}

// AfterCreate implements the AfterCreateHook interface
func (u *User) AfterCreate() error {
    fmt.Printf("User created: %s\n", u.Name)
    return nil
}
```

## Hook Execution Order

When multiple hooks are defined, they are executed in the following order:

1. Global hooks (registered with the repository)
2. Entity-specific hooks (methods on the entity struct)

For a save operation, the hooks are executed in this order:

1. `BeforeSave`
2. `BeforeCreate` (for new entities) or `BeforeUpdate` (for existing entities)
3. The actual database operation (INSERT or UPDATE)
4. `AfterCreate` (for new entities) or `AfterUpdate` (for existing entities)
5. `AfterSave`

## Error Handling

If any hook returns an error, the operation is aborted and the error is returned to the caller:

```go
func (u *User) BeforeSave() error {
    if u.Name == "" {
        return errors.New("name cannot be empty")
    }
    u.UpdatedAt = time.Now()
    return nil
}
```

In this example, if the user's name is empty, the save operation will be aborted and the error will be returned.

## Global Hooks

You can register global hooks with a repository to apply them to all entities of a specific type:

```go
// Create a timestamp hook
type TimestampHook struct{}

// BeforeSave sets timestamps
func (h *TimestampHook) BeforeSave(entity interface{}) error {
    if ts, ok := entity.(TimestampEntity); ok {
        now := time.Now()
        if ts.GetID() == 0 {
            ts.SetCreatedAt(now)
        }
        ts.SetUpdatedAt(now)
    }
    return nil
}

// TimestampEntity interface for entities with timestamps
type TimestampEntity interface {
    GetID() uint
    SetCreatedAt(time.Time)
    SetUpdatedAt(time.Time)
}

// Register the hook with the repository
userRepo.RegisterHook(&TimestampHook{})
```

## Common Hook Use Cases

### Automatic Timestamps

```go
func (u *User) BeforeSave() error {
    now := time.Now()
    if u.ID == 0 {
        u.CreatedAt = now
    }
    u.UpdatedAt = now
    return nil
}
```

### Validation

```go
func (u *User) BeforeSave() error {
    if u.Name == "" {
        return errors.New("name cannot be empty")
    }
    if !strings.Contains(u.Email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}
```

### Logging

```go
func (u *User) AfterCreate() error {
    log.Printf("User created: ID=%d, Name=%s, Email=%s", u.ID, u.Name, u.Email)
    return nil
}

func (u *User) AfterUpdate() error {
    log.Printf("User updated: ID=%d, Name=%s, Email=%s", u.ID, u.Name, u.Email)
    return nil
}

func (u *User) AfterDelete() error {
    log.Printf("User deleted: ID=%d, Name=%s, Email=%s", u.ID, u.Name, u.Email)
    return nil
}
```

### Password Hashing

```go
func (u *User) BeforeSave() error {
    // Only hash the password if it has changed
    if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}
```

### Generating Slugs

```go
func (p *Post) BeforeSave() error {
    if p.Slug == "" {
        p.Slug = generateSlug(p.Title)
    }
    return nil
}

func generateSlug(title string) string {
    // Convert to lowercase
    slug := strings.ToLower(title)
    // Replace spaces with hyphens
    slug = strings.ReplaceAll(slug, " ", "-")
    // Remove special characters
    slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")
    // Remove multiple hyphens
    slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
    // Trim hyphens from start and end
    slug = strings.Trim(slug, "-")
    return slug
}
```

## Hooks in Transactions

Hooks are transaction-aware, meaning that if a hook returns an error during a transaction, the transaction will be rolled back:

```go
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user (hooks will be executed)
    if err := txRepo.Save(user); err != nil {
        return err
    }
    
    // Create a profile for the user
    profile := &Profile{
        UserID: user.ID,
        Bio:    "Software developer",
    }
    
    // Save the profile (hooks will be executed)
    profileRepo := repository.NewRepository[Profile](txRepo.DB(), txRepo.Dialect())
    if err := profileRepo.Save(profile); err != nil {
        return err
    }
    
    return nil
})

if err != nil {
    log.Fatalf("Transaction failed: %v", err)
}
```

## Best Practices

### Keep Hooks Focused

Each hook should have a single responsibility:

```go
// Good: Focused hook
func (u *User) BeforeSave() error {
    u.UpdatedAt = time.Now()
    return nil
}

// Bad: Hook doing too much
func (u *User) BeforeSave() error {
    u.UpdatedAt = time.Now()
    u.Email = strings.ToLower(u.Email)
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    log.Printf("Saving user: %s", u.Name)
    return nil
}
```

### Use Appropriate Hook Types

Choose the right hook for the job:

- Use `BeforeSave` for validation and setting timestamps
- Use `AfterCreate` for logging and sending notifications
- Use `BeforeDelete` for validation and archiving
- Use `AfterDelete` for cleanup

### Handle Errors Properly

Always return errors from hooks to abort the operation if necessary:

```go
func (u *User) BeforeSave() error {
    if u.Email == "" {
        return errors.New("email cannot be empty")
    }
    return nil
}
```

### Test Your Hooks

Write tests for your hooks to ensure they work as expected:

```go
func TestUserBeforeSave(t *testing.T) {
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Call the hook directly
    err := user.BeforeSave()
    if err != nil {
        t.Errorf("BeforeSave returned an error: %v", err)
    }
    
    // Check that the timestamp was set
    if user.UpdatedAt.IsZero() {
        t.Error("UpdatedAt was not set")
    }
}
```

## Next Steps

- Learn about [Validation](./validation) to see how hooks can be used for validation
- Explore the [Repository Pattern](./repository-pattern) to understand how hooks integrate with CRUD operations
- Check out [Transactions](./transactions) to see how hooks work in transactions