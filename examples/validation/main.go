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
	"github.com/gooferOrm/goofer/validation"
)

// User entity with validation tags
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required,min=3,max=50"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	Age       int       `orm:"type:int;default:0" validate:"gte=0,lte=120"`
	Role      string    `orm:"type:varchar(20);notnull" validate:"required,oneof=admin user guest"`
	Active    bool      `orm:"type:boolean;default:true"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// Custom validation method
func (u *User) Validate() error {
	// Create a validator
	validator := validation.NewValidator()
	
	// Validate using struct tags
	errors, err := validator.ValidateEntity(u)
	if err != nil {
		return err
	}
	
	// Check if there are validation errors
	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errors {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
		return fmt.Errorf("validation failed")
	}
	
	// Custom validation logic
	if u.Role == "admin" && u.Age < 18 {
		return fmt.Errorf("admin users must be at least 18 years old")
	}
	
	return nil
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./validation.db")
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
	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create repository
	userRepo := repository.NewRepository[User](db, sqliteDialect)

	// Create validator
	validator := validation.NewValidator()

	fmt.Println("=== Testing validation ===")

	// Example 1: Valid user
	validUser := &User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Role:   "admin",
		Active: true,
	}

	// Validate before saving
	errors, err := validator.ValidateEntity(validUser)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors found:")
		for _, e := range errors {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	} else {
		fmt.Println("Valid user passed validation")
		if err := userRepo.Save(validUser); err != nil {
			log.Fatalf("Failed to save valid user: %v", err)
		}
		fmt.Printf("Created valid user with ID: %d\n", validUser.ID)
	}

	// Example 2: Invalid user - missing required fields
	invalidUser1 := &User{
		// Name is missing
		Email:  "invalid@example.com",
		Age:    -5, // Invalid age
		Role:   "superuser", // Invalid role
		Active: true,
	}

	// Validate before saving
	errors, err = validator.ValidateEntity(invalidUser1)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	if len(errors) > 0 {
		fmt.Println("\nValidation errors found (as expected):")
		for _, e := range errors {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	} else {
		// This should not happen
		fmt.Println("No validation errors found (unexpected)")
		if err := userRepo.Save(invalidUser1); err != nil {
			log.Fatalf("Failed to save invalid user: %v", err)
		}
	}

	// Example 3: Invalid user - custom validation rule
	invalidUser2 := &User{
		Name:   "Young Admin",
		Email:  "young@example.com",
		Age:    16, // Too young for admin
		Role:   "admin",
		Active: true,
	}

	// Use custom validation method
	if err := invalidUser2.Validate(); err != nil {
		fmt.Printf("\nCustom validation error: %v\n", err)
	} else {
		// This should not happen
		fmt.Println("No custom validation errors found (unexpected)")
		if err := userRepo.Save(invalidUser2); err != nil {
			log.Fatalf("Failed to save invalid user: %v", err)
		}
	}

	// Example 4: Valid user with custom validation
	validUser2 := &User{
		Name:   "Adult Admin",
		Email:  "adult@example.com",
		Age:    25,
		Role:   "admin",
		Active: true,
	}

	// Use custom validation method
	if err := validUser2.Validate(); err != nil {
		fmt.Printf("Custom validation error: %v\n", err)
	} else {
		fmt.Println("\nValid user passed custom validation")
		if err := userRepo.Save(validUser2); err != nil {
			log.Fatalf("Failed to save valid user: %v", err)
		}
		fmt.Printf("Created valid user with ID: %d\n", validUser2.ID)
	}

	// Fetch and display all saved users
	users, err := userRepo.Find().All()
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	fmt.Printf("\n=== Saved users (%d) ===\n", len(users))
	for _, u := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d, Role: %s, Active: %t\n",
			u.ID, u.Name, u.Email, u.Age, u.Role, u.Active)
	}
}