# Validation

Validation in Goofer ORM ensures that your data meets your requirements before it hits the database. It integrates with the popular `go-playground/validator` package to provide a robust validation system.

## Overview

The Validation system offers the following capabilities:

- Integration with go-playground/validator
- Struct tag support for validation rules
- Custom validation hooks
- Validation before save operations
- Detailed validation error messages

## Basic Validation

Goofer ORM uses the `validate` tag to define validation rules for your entity fields:

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
    Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
    Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
    Age       int       `orm:"type:int" validate:"gte=0,lte=130"`
    Password  string    `orm:"type:varchar(255);notnull" validate:"required,min=8"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
```

## Validation Tags

Goofer ORM supports all validation tags from the `go-playground/validator` package. Here are some commonly used tags:

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must not be empty | `validate:"required"` |
| `email` | Field must be a valid email address | `validate:"email"` |
| `min=X` | Minimum length for strings, minimum value for numbers | `validate:"min=8"` |
| `max=X` | Maximum length for strings, maximum value for numbers | `validate:"max=100"` |
| `gte=X` | Greater than or equal to X | `validate:"gte=0"` |
| `lte=X` | Less than or equal to X | `validate:"lte=130"` |
| `oneof=X Y Z` | Field must be one of the specified values | `validate:"oneof=admin user guest"` |
| `unique` | Field must be unique (requires database check) | `validate:"unique"` |
| `url` | Field must be a valid URL | `validate:"url"` |
| `uuid` | Field must be a valid UUID | `validate:"uuid"` |
| `datetime` | Field must be a valid datetime | `validate:"datetime=2006-01-02"` |

You can combine multiple validation tags by separating them with commas:

```go
Email string `validate:"required,email"`
```

## Using the Validator

To validate an entity, use the `Validator` from the validation package:

```go
import "github.com/gooferOrm/goofer/pkg/validation"

// Create a validator
validator := validation.NewValidator()

// Validate an entity
user := &User{
    Name:  "John Doe",
    Email: "invalid-email",
    Age:   150,
}

// Validate the entity
errors, err := validator.ValidateEntity(user)
if err != nil {
    log.Fatalf("Validation error: %v", err)
}

// Check for validation errors
if len(errors) > 0 {
    for _, e := range errors {
        fmt.Printf("Field %s: %s\n", e.Field, e.Message)
    }
    return
}

// If no errors, save the entity
if err := userRepo.Save(user); err != nil {
    log.Fatalf("Failed to save user: %v", err)
}
```

## Validation Hooks

You can implement custom validation logic by adding a `Validate` method to your entity:

```go
// Validate implements custom validation logic
func (u *User) Validate() error {
    // Custom validation logic
    if strings.Contains(u.Name, "admin") && u.Age < 18 {
        return errors.New("admin users must be at least 18 years old")
    }
    return nil
}
```

Entities that implement the `ValidatableEntity` interface will have their `Validate` method called automatically during validation:

```go
// ValidatableEntity is an interface for entities that can validate themselves
type ValidatableEntity interface {
    schema.Entity
    Validate() error
}
```

## Automatic Validation

Goofer ORM can automatically validate entities before saving them by using the `ValidateHook`:

```go
// Create a validate hook
validateHook := validation.NewValidateHook()

// Register the hook with the repository
userRepo.RegisterHook(validateHook)

// Now, validation will be performed automatically before saving
if err := userRepo.Save(user); err != nil {
    // Check if it's a validation error
    if validationErr, ok := err.(validation.ValidationError); ok {
        fmt.Printf("Validation failed: %v\n", validationErr)
        return
    }
    log.Fatalf("Failed to save user: %v", err)
}
```

## Validation Error Handling

Validation errors are returned as a slice of `ValidationError` structs:

```go
// ValidationError represents a validation error
type ValidationError struct {
    Field   string
    Message string
}
```

You can handle these errors in your application to provide user-friendly error messages:

```go
errors, err := validator.ValidateEntity(user)
if err != nil {
    log.Fatalf("Validation error: %v", err)
}

if len(errors) > 0 {
    // Create a map of field errors for your API response
    fieldErrors := make(map[string]string)
    for _, e := range errors {
        fieldErrors[e.Field] = e.Message
    }
    
    // Return the errors to the client
    response := map[string]interface{}{
        "success": false,
        "errors":  fieldErrors,
    }
    
    // Convert to JSON and send in your HTTP response
    jsonResponse, _ := json.Marshal(response)
    fmt.Println(string(jsonResponse))
    return
}
```

## Custom Validators

You can register custom validators with the validator:

```go
// Create a validator
validator := validation.NewValidator()

// Register a custom validator
validator.RegisterValidation("is_admin_email", func(fl validator.FieldLevel) bool {
    return strings.HasSuffix(fl.Field().String(), "@admin.com")
})

// Use the custom validator in your entity
type Admin struct {
    ID    uint   `orm:"primaryKey;autoIncrement"`
    Email string `orm:"unique;type:varchar(255);notnull" validate:"required,email,is_admin_email"`
}
```

## Validation in Transactions

Validation works seamlessly with transactions:

```go
err := userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Save the user (validation will happen automatically if the hook is registered)
    if err := txRepo.Save(user); err != nil {
        return err
    }
    
    return nil
})

if err != nil {
    log.Fatalf("Transaction failed: %v", err)
}
```

## Best Practices

### Use Appropriate Validation Rules

Choose validation rules that match your business requirements:

```go
type User struct {
    Username string `validate:"required,min=3,max=50"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"gte=18,lte=130"`
    Role     string `validate:"required,oneof=admin user guest"`
}
```

### Combine ORM and Validation Tags

Use both ORM and validation tags to ensure data integrity:

```go
type User struct {
    Email string `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
}
```

### Custom Validation Methods

Implement custom validation methods for complex validation logic:

```go
func (u *User) Validate() error {
    // Custom validation logic
    return nil
}
```

### Validation Error Messages

Provide clear and helpful validation error messages:

```go
if len(errors) > 0 {
    for _, e := range errors {
        switch e.Field {
        case "Email":
            fmt.Println("Please enter a valid email address")
        case "Password":
            fmt.Println("Password must be at least 8 characters long")
        default:
            fmt.Printf("%s: %s\n", e.Field, e.Message)
        }
    }
}
```

## Next Steps

- Learn about [Hooks](./hooks) to understand how to add lifecycle events to your entities
- Explore the [Repository Pattern](./repository-pattern) to see how validation integrates with CRUD operations
- Check out the [Examples](../examples/validation) section for more examples of validation