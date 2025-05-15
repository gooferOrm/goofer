package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yourusername/yourproject/goofer"
)

func main() {
	// Initialize the application
	client, err := goofer.Init()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer goofer.Close(client)

	// Get a repository for the User entity
	userRepo := client.GetRepository(&goofer.User{})

	// Create a new user
	user := &goofer.User{
		Name:  "Tach",
		Email: "tach@example.com",
	}

	if err := userRepo.Create(user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created user with ID: %d", user.ID)

	// Get a repository for the Post entity
	postRepo := client.GetRepository(&goofer.Post{})

	// Create a new post
	post := &goofer.Post{
		Title:   "Hello, Goofer!",
		Content: "This is my first post using Goofer ORM",
		UserID:  user.ID,
	}

	if err := postRepo.Create(post); err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	log.Printf("Created post with ID: %d", post.ID)

	// Query posts
	var posts []*goofer.Post
	if err := postRepo.FindAll(&posts); err != nil {
		log.Fatalf("Failed to fetch posts: %v", err)
	}

	// Print all posts
	fmt.Println("\nAll posts:")
	for _, p := range posts {
		fmt.Printf("ID: %d, Title: %s, Content: %s, UserID: %d\n", 
			p.ID, p.Title, p.Content, p.UserID)
	}

	// Example of finding a user by ID
	var foundUser goofer.User
	if err := userRepo.FindByID(user.ID, &foundUser); err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("\nFound user: %+v\n", foundUser)
}
