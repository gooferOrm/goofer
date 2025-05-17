package main

import (
	"fmt"
	"log"
	"with_engine_client/goofer"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the application
	client, err := goofer.Init()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer goofer.Close(client)

	// Get a repository for the User entity
	userRepo := client.Repository(&goofer.User{}).(*goofer.UserRepo)

	// Create a new user
	user := &goofer.User{
		Name:  "Tach",
		Email: "tach@example.com",
	}

	if err := userRepo.Create(user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created user with ID: %d", user.ID)

	// Find the user by ID
	foundUser2, err := userRepo.FindByID(int64(user.ID))
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	log.Printf("Found user: %s", foundUser2.Name)

	// Get a repository for the Post entity
	postRepo := client.Repository(&goofer.Post{}).(*goofer.PostRepo)

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
	// TODO: Implement post repository methods

	// Example of finding a user by ID
	foundUser2, err = userRepo.FindByID(int64(user.ID))
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("\nFound user: %+v\n", foundUser2)
}
