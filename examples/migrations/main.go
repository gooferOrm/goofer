package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/migration"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// Initial User entity (v1)
type UserV1 struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (UserV1) TableName() string {
	return "users"
}

// Updated User entity (v2) - Added Age field
type UserV2 struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	Age       int       `orm:"type:int;default:0" validate:"gte=0,lte=120"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (UserV2) TableName() string {
	return "users"
}

// Final User entity (v3) - Added Address field and renamed Age to YearsOld
type UserV3 struct {
	ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
	Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
	Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
	YearsOld  int       `orm:"type:int;default:0" validate:"gte=0,lte=120"`
	Address   string    `orm:"type:varchar(255)" validate:"omitempty"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the User entity
func (UserV3) TableName() string {
	return "users"
}

// Migration from v1 to v2 - Add Age column
func migrateV1ToV2(db *sql.DB, dialect dialect.Dialect) error {
	// Add Age column to users table
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s DEFAULT 0",
		dialect.QuoteIdentifier("users"),
		dialect.QuoteIdentifier("age"),
		"INTEGER")

	_, err := db.Exec(sql)
	return err
}

// Migration from v2 to v3 - Rename Age to YearsOld and add Address column
func migrateV2ToV3(db *sql.DB, dialect dialect.Dialect) error {
	// For SQLite, we need to create a new table, copy data, and drop the old table
	// since SQLite doesn't support renaming columns directly

	// 1. Create a temporary table with the new schema
	tempTableSQL := `CREATE TABLE users_new (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		years_old INTEGER DEFAULT 0,
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(tempTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create temporary table: %v", err)
	}

	// 2. Copy data from the old table to the new table
	copySQL := `INSERT INTO users_new (id, name, email, years_old, created_at)
		SELECT id, name, email, age, created_at FROM users`

	_, err = db.Exec(copySQL)
	if err != nil {
		return fmt.Errorf("failed to copy data: %v", err)
	}

	// 3. Drop the old table
	_, err = db.Exec("DROP TABLE users")
	if err != nil {
		return fmt.Errorf("failed to drop old table: %v", err)
	}

	// 4. Rename the new table to the original name
	_, err = db.Exec("ALTER TABLE users_new RENAME TO users")
	if err != nil {
		return fmt.Errorf("failed to rename table: %v", err)
	}

	return nil
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./migrations.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Create migration manager
	migrationManager := migration.NewManager(db, sqliteDialect)

	// Register migrations
	migrationManager.RegisterMigration("v1_to_v2", migrateV1ToV2)
	migrationManager.RegisterMigration("v2_to_v3", migrateV2ToV3)

	// Step 1: Create initial schema (v1)
	fmt.Println("Step 1: Creating initial schema (v1)")
	
	// Register v1 entity
	if err := schema.Registry.RegisterEntity(UserV1{}); err != nil {
		log.Fatalf("Failed to register UserV1 entity: %v", err)
	}

	// Get entity metadata
	userV1Meta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(UserV1{}))

	// Create table
	userV1SQL := sqliteDialect.CreateTableSQL(userV1Meta)
	fmt.Println("User table SQL (v1):")
	fmt.Println(userV1SQL)

	_, err = db.Exec(userV1SQL)
	if err != nil {
		log.Fatalf("Failed to create users table (v1): %v", err)
	}

	// Create repository for v1
	userV1Repo := repository.NewRepository[UserV1](db, sqliteDialect)

	// Insert some test data
	user1 := &UserV1{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if err := userV1Repo.Save(user1); err != nil {
		log.Fatalf("Failed to save user (v1): %v", err)
	}

	user2 := &UserV1{
		Name:  "Jane Smith",
		Email: "jane@example.com",
	}

	if err := userV1Repo.Save(user2); err != nil {
		log.Fatalf("Failed to save user (v1): %v", err)
	}

	fmt.Printf("Created users with IDs: %d, %d\n", user1.ID, user2.ID)

	// Step 2: Migrate to v2 (add Age column)
	fmt.Println("\nStep 2: Migrating to v2 (adding Age column)")
	
	if err := migrationManager.RunMigration("v1_to_v2"); err != nil {
		log.Fatalf("Failed to run migration v1_to_v2: %v", err)
	}

	// Clear the registry and register v2 entity
	schema.Registry = schema.NewSchemaRegistry()
	if err := schema.Registry.RegisterEntity(UserV2{}); err != nil {
		log.Fatalf("Failed to register UserV2 entity: %v", err)
	}

	// Create repository for v2
	userV2Repo := repository.NewRepository[UserV2](db, sqliteDialect)

	// Find users and update their ages
	usersV2, err := userV2Repo.Find().All()
	if err != nil {
		log.Fatalf("Failed to find users (v2): %v", err)
	}

	fmt.Printf("Found %d users after v1_to_v2 migration:\n", len(usersV2))
	for i, u := range usersV2 {
		fmt.Printf("- %s (%s), Age: %d\n", u.Name, u.Email, u.Age)
		
		// Update ages
		u.Age = 30 + i
		if err := userV2Repo.Save(&u); err != nil {
			log.Fatalf("Failed to update user age (v2): %v", err)
		}
	}

	// Step 3: Migrate to v3 (rename Age to YearsOld and add Address)
	fmt.Println("\nStep 3: Migrating to v3 (renaming Age to YearsOld and adding Address)")
	
	if err := migrationManager.RunMigration("v2_to_v3"); err != nil {
		log.Fatalf("Failed to run migration v2_to_v3: %v", err)
	}

	// Clear the registry and register v3 entity
	schema.Registry = schema.NewSchemaRegistry()
	if err := schema.Registry.RegisterEntity(UserV3{}); err != nil {
		log.Fatalf("Failed to register UserV3 entity: %v", err)
	}

	// Create repository for v3
	userV3Repo := repository.NewRepository[UserV3](db, sqliteDialect)

	// Find users and update their addresses
	usersV3, err := userV3Repo.Find().All()
	if err != nil {
		log.Fatalf("Failed to find users (v3): %v", err)
	}

	fmt.Printf("Found %d users after v2_to_v3 migration:\n", len(usersV3))
	for i, u := range usersV3 {
		fmt.Printf("- %s (%s), YearsOld: %d, Address: %s\n", u.Name, u.Email, u.YearsOld, u.Address)
		
		// Update addresses
		u.Address = fmt.Sprintf("%d Main St, City %d", 100+i, i)
		if err := userV3Repo.Save(&u); err != nil {
			log.Fatalf("Failed to update user address (v3): %v", err)
		}
	}

	// Final check
	fmt.Println("\nFinal state after all migrations:")
	finalUsers, err := userV3Repo.Find().All()
	if err != nil {
		log.Fatalf("Failed to find users (final): %v", err)
	}

	for _, u := range finalUsers {
		fmt.Printf("- %s (%s), YearsOld: %d, Address: %s\n", u.Name, u.Email, u.YearsOld, u.Address)
	}

	fmt.Println("\nMigrations completed successfully!")
}