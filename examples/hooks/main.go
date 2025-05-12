package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// User entity with lifecycle hooks
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp"`
	UpdatedAt time.Time `orm:"type:timestamp"`
	LastLogin *time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// BeforeCreate is called before creating a new record
func (u *User) BeforeCreate() error {
	fmt.Println("BeforeCreate hook called")
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// AfterCreate is called after creating a new record
func (u *User) AfterCreate() error {
	fmt.Println("AfterCreate hook called")
	fmt.Printf("User %s created at %v\n", u.Name, u.CreatedAt)
	return nil
}

// BeforeUpdate is called before updating a record
func (u *User) BeforeUpdate() error {
	fmt.Println("BeforeUpdate hook called")
	u.UpdatedAt = time.Now()
	return nil
}

// AfterUpdate is called after updating a record
func (u *User) AfterUpdate() error {
	fmt.Println("AfterUpdate hook called")
	fmt.Printf("User %s updated at %v\n", u.Name, u.UpdatedAt)
	return nil
}

// BeforeSave is called before saving (create or update) a record
func (u *User) BeforeSave() error {
	fmt.Println("BeforeSave hook called")
	
	// Normalize email
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	
	// Normalize name
	u.Name = strings.TrimSpace(u.Name)
	
	return nil
}

// AfterSave is called after saving (create or update) a record
func (u *User) AfterSave() error {
	fmt.Println("AfterSave hook called")
	return nil
}

// BeforeDelete is called before deleting a record
func (u *User) BeforeDelete() error {
	fmt.Println("BeforeDelete hook called")
	fmt.Printf("About to delete user %s (ID: %d)\n", u.Name, u.ID)
	return nil
}

// AfterDelete is called after deleting a record
func (u *User) AfterDelete() error {
	fmt.Println("AfterDelete hook called")
	return nil
}

// Activity entity with auto-logging
type Activity struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	UserID    uint      `orm:"index;notnull"`
	Action    string    `orm:"type:varchar(50);notnull"`
	Details   string    `orm:"type:text"`
	Timestamp time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the Activity entity
func (Activity) TableName() string {
	return "activities"
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./hooks.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Activity{}); err != nil {
		log.Fatalf("Failed to register Activity entity: %v", err)
	}

	// Get entity metadata
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
	activityMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Activity{}))

	// Create tables
	userSQL := sqliteDialect.CreateTableSQL(userMeta)
	activitySQL := sqliteDialect.CreateTableSQL(activityMeta)

	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(activitySQL)
	if err != nil {
		log.Fatalf("Failed to create activities table: %v", err)
	}

	// Create repositories
	userRepo := repository.NewRepository[User](db, sqliteDialect)
	activityRepo := repository.NewRepository[Activity](db, sqliteDialect)

	fmt.Println("=== Testing Lifecycle Hooks ===")

	// Create a user
	user := &User{
		Name:  "John Doe",
		Email: "JOHN@example.com", // Will be normalized in BeforeSave
	}

	fmt.Println("\n--- Creating User ---")
	if err := userRepo.Save(user); err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	// Log activity manually (in a real app, this would be in the AfterCreate hook)
	activity := &Activity{
		UserID:    user.ID,
		Action:    "create",
		Details:   fmt.Sprintf("User created: %s", user.Name),
		Timestamp: time.Now(),
	}
	if err := activityRepo.Save(activity); err != nil {
		log.Fatalf("Failed to save activity: %v", err)
	}

	fmt.Printf("\nCreated user with ID: %d\n", user.ID)
	fmt.Printf("Notice that the email was normalized: %s\n", user.Email)

	// Update the user
	fmt.Println("\n--- Updating User ---")
	user.Name = "Jane Doe"
	if err := userRepo.Save(user); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	// Log activity
	activity = &Activity{
		UserID:    user.ID,
		Action:    "update",
		Details:   fmt.Sprintf("User updated: %s", user.Name),
		Timestamp: time.Now(),
	}
	if err := activityRepo.Save(activity); err != nil {
		log.Fatalf("Failed to save activity: %v", err)
	}

	fmt.Printf("\nUpdated user name to: %s\n", user.Name)
	fmt.Printf("Notice that UpdatedAt was set automatically: %v\n", user.UpdatedAt)

	// Fetch the user to see the changes
	updatedUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("\nFetched user: %s (%s)\n", updatedUser.Name, updatedUser.Email)
	fmt.Printf("Created At: %v\n", updatedUser.CreatedAt)
	fmt.Printf("Updated At: %v\n", updatedUser.UpdatedAt)

	// Delete the user
	fmt.Println("\n--- Deleting User ---")
	if err := userRepo.Delete(updatedUser); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	// Log activity
	activity = &Activity{
		UserID:    updatedUser.ID,
		Action:    "delete",
		Details:   fmt.Sprintf("User deleted: %s", updatedUser.Name),
		Timestamp: time.Now(),
	}
	if err := activityRepo.Save(activity); err != nil {
		log.Fatalf("Failed to save activity: %v", err)
	}

	fmt.Println("\nDeleted user")

	// Display activity log
	activities, err := activityRepo.Find().OrderBy("timestamp ASC").All()
	if err != nil {
		log.Fatalf("Failed to find activities: %v", err)
	}

	fmt.Printf("\n=== Activity Log (%d entries) ===\n", len(activities))
	for i, a := range activities {
		fmt.Printf("%d. [%v] User %d - %s: %s\n", 
			i+1, a.Timestamp.Format("2006-01-02 15:04:05"), 
			a.UserID, a.Action, a.Details)
	}
}