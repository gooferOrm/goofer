package main

import (
	"database/sql"
	"fmt"
	"log"


	"github.com/goferOrm/goofer/repository"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    uint   `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name  string `orm:"size:255;not null" json:"name"`
	Email string `orm:"size:255;not null;uniqueIndex" json:"email"`
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
	// Open a database connection
	db, err := sql.Open("sqlite3", "./without_engine.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create a new repository for User
	userRepo := repository.NewRepository[User](db, repository.SQLite)

	// Create a new user
	user := &User{
		Name:  "Jane Smith",
		Email: "jane@example.com",
	}

	if err := userRepo.Save(user); err != nil {
		log.Fatalf("failed to save user: %v", err)
	}
	log.Printf("Created user with ID: %d", user.ID)

	// Create a new repository for Post
	postRepo := repository.NewRepository[Post](db, repository.SQLite)

	// Create a new post
	post := &Post{
		Title:   "Direct Repository Usage",
		Content: "This post was created without using the engine",
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

	// Find user by ID
	foundUser, err := userRepo.FindByID(user.ID)
	if err != nil {
		log.Fatalf("failed to find user: %v", err)
	}
	fmt.Printf("\nFound user: %+v\n", foundUser)
}
