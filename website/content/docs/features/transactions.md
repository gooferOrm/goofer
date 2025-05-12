# Transactions

Transactions in Goofer ORM ensure data integrity for complex operations. They allow you to group multiple database operations into a single atomic unit, ensuring that either all operations succeed or none of them are applied.

## Overview

The Transactions system offers the following capabilities:

- First-class transaction support
- Automatic rollback on error
- Nested transaction support
- Transaction-aware hooks
- Context support for cancellation and timeouts

## Why Use Transactions?

Transactions are essential for maintaining data integrity in scenarios where multiple related operations need to succeed or fail together. For example:

- Creating a user and their profile
- Transferring money between accounts
- Processing an order with multiple items
- Updating related entities

Without transactions, if one operation fails, you might end up with inconsistent data in your database.

## Basic Transaction Usage

To use transactions in Goofer ORM, use the `Transaction` method on a repository:

```go
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

The `Transaction` method takes a function that receives a transaction-specific repository. Inside this function, you perform your database operations using the transaction repository. If the function returns an error, the transaction is automatically rolled back. If it returns nil, the transaction is committed.

## Error Handling and Rollback

If any operation inside the transaction function returns an error, the transaction is automatically rolled back:

```go
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user in the transaction
    if err := txRepo.Save(user); err != nil {
        return err // This will cause the transaction to roll back
    }
    
    // Simulate an error
    return errors.New("something went wrong") // This will cause the transaction to roll back
})

if err != nil {
    log.Printf("Transaction failed: %v", err) // This will print the error
}
```

You can also explicitly return an error to roll back the transaction if a business rule is violated:

```go
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
    
    // Check if the user already has a profile
    existingProfile, err := profileRepo.Find().Where("user_id = ?", user.ID).One()
    if err == nil {
        // User already has a profile
        return errors.New("user already has a profile")
    } else if err != sql.ErrNoRows {
        // Some other error occurred
        return err
    }
    
    // Create a profile for the user
    profile := &Profile{
        UserID: user.ID,
        Bio:    "Software developer",
    }
    
    // Save the profile in the transaction
    if err := profileRepo.Save(profile); err != nil {
        return err
    }
    
    return nil
})
```

## Panic Recovery

Transactions also handle panics, automatically rolling back if a panic occurs:

```go
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
    
    // Simulate a panic
    panic("something went wrong") // This will cause the transaction to roll back
    
    // This code will never be reached
    return nil
})

// The panic will be recovered, and err will contain the panic message
if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

## Working with Multiple Repositories

You can use multiple repositories within a single transaction:

```go
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
    
    // Create a profile repository using the same transaction
    profileRepo := repository.NewRepository[Profile](txRepo.DB(), txRepo.Dialect())
    
    // Create a profile for the user
    profile := &Profile{
        UserID: user.ID,
        Bio:    "Software developer",
    }
    
    // Save the profile in the transaction
    if err := profileRepo.Save(profile); err != nil {
        return err
    }
    
    // Create a post repository using the same transaction
    postRepo := repository.NewRepository[Post](txRepo.DB(), txRepo.Dialect())
    
    // Create a post for the user
    post := &Post{
        UserID:  user.ID,
        Title:   "My First Post",
        Content: "Hello, world!",
    }
    
    // Save the post in the transaction
    if err := postRepo.Save(post); err != nil {
        return err
    }
    
    return nil
})
```

## Hooks in Transactions

Hooks are transaction-aware, meaning they are executed within the transaction:

```go
// User entity with hooks
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp"`
    UpdatedAt time.Time `orm:"type:timestamp"`
}

// BeforeSave is called before the entity is saved
func (u *User) BeforeSave() error {
    now := time.Now()
    if u.ID == 0 {
        u.CreatedAt = now
    }
    u.UpdatedAt = now
    return nil
}

// Transaction with hooks
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user in the transaction (BeforeSave hook will be executed)
    if err := txRepo.Save(user); err != nil {
        return err
    }
    
    return nil
})
```

If a hook returns an error, the transaction is rolled back:

```go
// BeforeSave with validation
func (u *User) BeforeSave() error {
    if u.Name == "" {
        return errors.New("name cannot be empty")
    }
    
    now := time.Now()
    if u.ID == 0 {
        u.CreatedAt = now
    }
    u.UpdatedAt = now
    return nil
}

// Transaction with hooks
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user with an empty name
    user := &User{
        Name:  "", // This will cause the BeforeSave hook to return an error
        Email: "john@example.com",
    }
    
    // Save the user in the transaction (BeforeSave hook will return an error)
    if err := txRepo.Save(user); err != nil {
        return err // This will cause the transaction to roll back
    }
    
    return nil
})
```

## Context Support

Transactions support context for cancellation and timeouts:

```go
// Create a context with a timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use the context with the repository
userRepoWithCtx := userRepo.WithContext(ctx)

// Transaction with context
err := userRepoWithCtx.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user in the transaction
    if err := txRepo.Save(user); err != nil {
        return err
    }
    
    // Simulate a long-running operation
    time.Sleep(10 * time.Second) // This will exceed the context timeout
    
    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err) // This will print a context deadline exceeded error
}
```

## Best Practices

### Keep Transactions Focused

Each transaction should have a single responsibility:

```go
// Good: Focused transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a user and their profile
    return nil
})

// Bad: Transaction doing too much
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create multiple users
    // Process payments
    // Send emails
    // Update inventory
    return nil
})
```

### Minimize Transaction Duration

Keep transactions as short as possible to reduce lock contention:

```go
// Good: Short transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Only perform database operations
    return nil
})

// Bad: Long-running transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Perform database operations
    // Make HTTP requests
    // Process files
    return nil
})
```

### Handle Errors Properly

Always check errors returned by repository methods:

```go
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user in the transaction
    if err := txRepo.Save(user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### Use Transactions for Related Operations

Use transactions when operations are related and need to succeed or fail together:

```go
// Good: Related operations in a transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a user and their profile
    return nil
})

// Bad: Unrelated operations in a transaction
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a user
    // Update a product (unrelated to the user)
    return nil
})
```

## Next Steps

- Learn about [Hooks](./hooks) to understand how hooks work in transactions
- Explore the [Repository Pattern](./repository-pattern) to see how transactions integrate with repositories
- Check out the [Examples](../examples/basic) section for more examples of using transactions