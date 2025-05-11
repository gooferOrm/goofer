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
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// Post entity
type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Title     string    `orm:"type:varchar(255);notnull"`
	Content   string    `orm:"type:text"`
	UserID    uint      `orm:"index;notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}

// TableName returns the table name for the Post entity
func (Post) TableName() string {
	return "posts"
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./db.db")
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
	if err := schema.Registry.RegisterEntity(Post{}); err != nil {
		log.Fatalf("Failed to register Post entity: %v", err)
	}

	// Get entity metadata
	userMeta, ok := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
	if !ok {
		log.Fatal("Failed to get User entity metadata")
	}

	postMeta, ok := schema.Registry.GetEntityMetadata(schema.GetEntityType(Post{}))
	if !ok {
		log.Fatal("Failed to get Post entity metadata")
	}

	// Print the SQL for table creation
	userSQL := sqliteDialect.CreateTableSQL(userMeta)
	postSQL := sqliteDialect.CreateTableSQL(postMeta)

	fmt.Println("User table SQL:")
	fmt.Println(userSQL)
	fmt.Println("Post table SQL:")
	fmt.Println(postSQL)

	// Create tables
	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(postSQL)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	// Create repositories
	userRepo := repository.NewRepository[User](db, sqliteDialect)
	postRepo := repository.NewRepository[Post](db, sqliteDialect)

	// Create a user
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Save the user
	if err := userRepo.Save(user); err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Create a post
	post := &Post{
		Title:   "Hello, World!",
		Content: "This is my first post.",
		UserID:  user.ID,
	}

	// Save the post
	if err := postRepo.Save(post); err != nil {
		log.Fatalf("Failed to save post: %v", err)
	}

	fmt.Printf("Created post with ID: %d\n", post.ID)

	// Find the user by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)

	// Find posts by user ID
	posts, err := postRepo.Find().Where("user_id = ?", user.ID).All()
	if err != nil {
		log.Fatalf("Failed to find posts: %v", err)
	}

	fmt.Printf("Found %d posts by user %s:\n", len(posts), foundUser.Name)
	for _, p := range posts {
		fmt.Printf("- %s: %s\n", p.Title, p.Content)
	}

	// Update the user
	foundUser.Name = "Jane Doe"
	if err := userRepo.Save(foundUser); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	fmt.Printf("Updated user name to: %s\n", foundUser.Name)

	// Delete the post
	if err := postRepo.Delete(post); err != nil {
		log.Fatalf("Failed to delete post: %v", err)
	}

	fmt.Println("Deleted post")

	// Transaction example
	err = userRepo.Transaction(func(txRepo *repository.Repository[User]) error {
		// Create a new user in the transaction
		newUser := &User{
			Name:  "Transaction User",
			Email: "tx@example.com",
		}

		// Save the user in the transaction
		if err := txRepo.Save(newUser); err != nil {
			return err
		}

		fmt.Printf("Created user in transaction with ID: %d\n", newUser.ID)
		// Simulate an error to rollback the transaction
		// return errors.New("simulated error")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		fmt.Println("Transaction committed successfully")
	}
}
