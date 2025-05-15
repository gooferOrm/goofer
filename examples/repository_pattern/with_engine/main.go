package main

import (
	"fmt"
	"log"

	"github.com/gooferOrm/goofer/engine"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        uint   `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string `orm:"size:255;not null" json:"name"`
	Email     string `orm:"size:255;not null;uniqueIndex" json:"email"`
}

func (User) TableName() string { return "users" }

type Post struct {
	ID      uint   `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Title   string `orm:"size:255;not null" json:"title"`
	Content string `orm:"type:text;not null" json:"content"`
	UserID  uint   `orm:"index;not null" json:"user_id"`
}

func (Post) TableName() string { return "posts" }

func main() {
	// Initialize the Goofer client with SQLite
	client, err := engine.Connect("sqlite3", "./with_engine.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	// Register your entities with the client
	err = client.RegisterEntities(&User{}, &Post{})
	if err != nil {
		log.Fatalf("failed to register entities: %v", err)
	}

	// Get a repository for the User entity
	userRepo := engine.Repo[User](client)

	// Create a new user
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if err := userRepo.Save(user); err != nil {
		log.Fatalf("failed to save user: %v", err)
	}
	log.Printf("Created user with ID: %d", user.ID)

	// Get a repository for the Post entity
	postRepo := engine.Repo[Post](client)

	// Create a new post
	post := &Post{
		Title:   "Hello, Goofer!",
		Content: "This is my first post using Goofer ORM with engine",
		UserID:  user.ID,
	}

	if err := postRepo.Save(post); err != nil {
		log.Fatalf("failed to save post: %v", err)
	}
	log.Printf("Created post with ID: %d", post.ID)

	// Query posts
	posts, err := postRepo.Find().All()
	if err != nil {
		log.Fatalf("failed to fetch posts: %v", err)
	}

	// Print all posts
	fmt.Println("\nAll posts:")
	for _, p := range posts {
		fmt.Printf("ID: %d, Title: %s, Content: %s, UserID: %d\n",
			p.ID, p.Title, p.Content, p.UserID)
	}

	// Example of using the repository directly
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("failed to find user: %v", err)
	}
	fmt.Printf("\nFound user: %+v\n", foundUser)
}
