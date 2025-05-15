package goofer

import (
	"database/sql"
	"log"

	"github.com/gooferOrm/goofer/engine"
)

// DB holds the database connection
var DB *sql.DB

// Init initializes the database connection and returns a new client
func Init() (*engine.Client, error) {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return nil, err
	}

	// Create a new Goofer client
	client, err := engine.Connect("sqlite3", "./db.db")
	if err != nil {
		db.Close()
		return nil, err
	}

	// Register your entities
	err = client.RegisterEntities(&User{}, &Post{})
	if err != nil {
		client.Close()
		db.Close()
		return nil, err
	}

	// Set the global DB variable
	DB = db

	log.Println("Database connection established and entities registered")
	return client, nil
}

// Close cleans up the database connection
func Close(client *engine.Client) {
	if client != nil {
		client.Close()
	}
	if DB != nil {
		DB.Close()
	}
}