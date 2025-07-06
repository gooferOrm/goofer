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
	Age       int       `orm:"type:int"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./advanced_queries.db")
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

	// Create table
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
	userSQL := sqliteDialect.CreateTableSQL(userMeta)
	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create repository
	userRepo := repository.NewRepository[User](db, sqliteDialect)

	// Create sample data
	users := createSampleData(userRepo)

	fmt.Println("=== Enhanced Query Builder Examples ===\n")

	// Example 1: WHERE IN clause
	fmt.Println("1. WHERE IN clause:")
	userIDs := []interface{}{users[0].ID, users[1].ID}
	specificUsers, err := userRepo.Find().
		WhereIn("id", userIDs).
		All()
	if err != nil {
		log.Printf("Error finding specific users: %v", err)
	} else {
		fmt.Printf("Found %d specific users\n", len(specificUsers))
	}

	// Example 2: BETWEEN clause
	fmt.Println("\n2. BETWEEN clause:")
	recentUsers, err := userRepo.Find().
		WhereBetween("created_at", time.Now().AddDate(0, 0, -7), time.Now()).
		All()
	if err != nil {
		log.Printf("Error finding recent users: %v", err)
	} else {
		fmt.Printf("Found %d recent users\n", len(recentUsers))
	}

	// Example 3: LIKE conditions
	fmt.Println("\n3. LIKE conditions:")
	usersByName, err := userRepo.Find().
		WhereLike("name", "%Doe%").
		All()
	if err != nil {
		log.Printf("Error finding users by name: %v", err)
	} else {
		fmt.Printf("Found %d users with 'Doe' in name\n", len(usersByName))
	}

	// Example 4: IS NULL / IS NOT NULL
	fmt.Println("\n4. IS NULL / IS NOT NULL:")
	usersWithAge, err := userRepo.Find().
		WhereNotNull("age").
		All()
	if err != nil {
		log.Printf("Error finding users with age: %v", err)
	} else {
		fmt.Printf("Found %d users with age set\n", len(usersWithAge))
	}

	// Example 5: OR conditions
	fmt.Println("\n5. OR conditions:")
	usersByNameOrEmail, err := userRepo.Find().
		Where("name = ?", "John Doe").
		OrWhere("email = ?", "jane@example.com").
		All()
	if err != nil {
		log.Printf("Error finding users by name or email: %v", err)
	} else {
		fmt.Printf("Found %d users by name or email\n", len(usersByNameOrEmail))
	}

	// Example 6: DISTINCT
	fmt.Println("\n6. DISTINCT:")
	uniqueAges, err := userRepo.Find().
		Distinct().
		OrderBy("age ASC").
		All()
	if err != nil {
		log.Printf("Error finding unique ages: %v", err)
	} else {
		fmt.Printf("Found users with %d different ages\n", len(uniqueAges))
	}

	// Example 7: Advanced filtering
	fmt.Println("\n7. Advanced filtering:")
	advancedUsers, err := userRepo.Find().
		Where("age > ?", 20).
		Where("age < ?", 50).
		WhereLike("name", "%Doe%").
		OrderBy("age DESC").
		Limit(5).
		All()
	if err != nil {
		log.Printf("Error finding advanced users: %v", err)
	} else {
		fmt.Printf("Found %d advanced filtered users\n", len(advancedUsers))
	}

	// Example 8: Count with conditions
	fmt.Println("\n8. Count with conditions:")
	count, err := userRepo.Find().
		Where("age >= ?", 25).
		Count()
	if err != nil {
		log.Printf("Error counting users: %v", err)
	} else {
		fmt.Printf("Count of users aged 25+: %d\n", count)
	}

	fmt.Println("\n=== New Query Builder Features ===")
	fmt.Println("✅ WHERE IN clauses")
	fmt.Println("✅ BETWEEN clauses")
	fmt.Println("✅ LIKE conditions")
	fmt.Println("✅ IS NULL / IS NOT NULL")
	fmt.Println("✅ OR conditions")
	fmt.Println("✅ DISTINCT")
	fmt.Println("✅ Advanced filtering combinations")
	fmt.Println("✅ Enhanced COUNT queries")
}

// createSampleData creates sample data for testing
func createSampleData(userRepo *repository.Repository[User]) []User {
	users := []User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
		{Name: "Alice Brown", Email: "alice@example.com", Age: 28},
	}

	for i := range users {
		if err := userRepo.Save(&users[i]); err != nil {
			log.Printf("Failed to save user: %v", err)
		}
	}

	return users
}
