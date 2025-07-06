package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/introspection"
)

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./introspection.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Create some sample tables for introspection
	createSampleTables(db)

	// Create introspector
	introspector := introspection.NewIntrospector(db, sqliteDialect)

	fmt.Println("=== Database Introspection Example ===\n")

	// Introspect all tables
	fmt.Println("1. Introspecting all tables:")
	tables, err := introspector.IntrospectAllTables()
	if err != nil {
		log.Fatalf("Failed to introspect tables: %v", err)
	}

	for _, table := range tables {
		fmt.Printf("Table: %s\n", table.Name)
		fmt.Printf("  Columns: %d\n", len(table.Columns))
		fmt.Printf("  Primary Key: %s\n", table.PrimaryKey)
		fmt.Printf("  Indexes: %d\n", len(table.Indexes))
		fmt.Printf("  Foreign Keys: %d\n", len(table.ForeignKeys))
		fmt.Println()
	}

	// Introspect a specific table
	fmt.Println("2. Introspecting specific table (users):")
	userTable, err := introspector.IntrospectTable("users")
	if err != nil {
		log.Fatalf("Failed to introspect users table: %v", err)
	}

	fmt.Printf("Table: %s\n", userTable.Name)
	fmt.Println("Columns:")
	for _, col := range userTable.Columns {
		fmt.Printf("  - %s (%s)", col.Name, col.Type)
		if col.IsPrimaryKey {
			fmt.Print(" [PRIMARY KEY]")
		}
		if !col.IsNullable {
			fmt.Print(" [NOT NULL]")
		}
		if col.IsUnique {
			fmt.Print(" [UNIQUE]")
		}
		if col.DefaultValue != nil {
			fmt.Printf(" [DEFAULT: %s]", *col.DefaultValue)
		}
		fmt.Println()
	}

	// Generate Go entities
	fmt.Println("\n3. Generating Go entities:")
	entities, err := introspector.GenerateEntities()
	if err != nil {
		log.Fatalf("Failed to generate entities: %v", err)
	}

	fmt.Println("Generated Go code:")
	fmt.Println("```go")
	fmt.Println(entities)
	fmt.Println("```")

	// Generate entity for a specific table
	fmt.Println("\n4. Generating entity for users table:")
	userEntity, err := introspector.GenerateEntity(userTable)
	if err != nil {
		log.Fatalf("Failed to generate user entity: %v", err)
	}

	fmt.Println("Generated User entity:")
	fmt.Println("```go")
	fmt.Println(userEntity)
	fmt.Println("```")

	fmt.Println("\n=== Introspection Features ===")
	fmt.Println("✅ Database schema analysis")
	fmt.Println("✅ Table information extraction")
	fmt.Println("✅ Column metadata parsing")
	fmt.Println("✅ Primary key detection")
	fmt.Println("✅ Index information")
	fmt.Println("✅ Foreign key relationships")
	fmt.Println("✅ Go struct generation")
	fmt.Println("✅ ORM tag generation")
	fmt.Println("✅ Type mapping (SQL to Go)")
}

// createSampleTables creates sample tables for introspection
func createSampleTables(db *sql.DB) {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			age INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			user_id INTEGER NOT NULL,
			status VARCHAR(50) DEFAULT 'draft',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	// Create comments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create comments table: %v", err)
	}

	// Insert sample data
	_, err = db.Exec(`
		INSERT INTO users (name, email, age) VALUES 
		('John Doe', 'john@example.com', 30),
		('Jane Smith', 'jane@example.com', 25)
	`)
	if err != nil {
		log.Printf("Failed to insert sample data: %v", err)
	}
}
