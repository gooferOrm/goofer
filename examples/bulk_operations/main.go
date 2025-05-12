package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// BulkUser entity for bulk operation example
type BulkUser struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Username  string    `orm:"type:varchar(50);notnull;unique"`
	Email     string    `orm:"type:varchar(100);notnull;unique"`
	Active    bool      `orm:"type:boolean;default:true"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the BulkUser entity
func (BulkUser) TableName() string {
	return "bulk_users"
}

// BulkRepository extends Repository to add bulk operations
type BulkRepository[T schema.Entity] struct {
	*repository.Repository[T]
	db       *sql.DB
	dialect  dialect.Dialect
	metadata *schema.EntityMetadata
}

// NewBulkRepository creates a new repository with bulk operation support
func NewBulkRepository[T schema.Entity](db *sql.DB, dialect dialect.Dialect) *BulkRepository[T] {
	var entity T
	entityType := schema.GetEntityType(entity)

	meta, exists := schema.Registry.GetEntityMetadata(entityType)
	if !exists {
		panic(fmt.Sprintf("entity %s not registered", entityType.Name()))
	}

	return &BulkRepository[T]{
		Repository: repository.NewRepository[T](db, dialect),
		db:         db,
		dialect:    dialect,
		metadata:   meta,
	}
}

// BulkInsert efficiently inserts multiple entities in a single operation
func (r *BulkRepository[T]) BulkInsert(entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare column names and placeholders
	var columnNames []string
	var placeholderGroups []string
	var allValues []interface{}

	// Skip primary key for auto-increment fields
	skipPK := false
	if r.metadata.PrimaryKey != nil && r.metadata.PrimaryKey.IsAutoIncr {
		skipPK = true
	}

	// Get column names from fields
	for _, field := range r.metadata.Fields {
		// Skip primary key if auto-increment and relation fields
		if (skipPK && field.IsPrimaryKey) || field.Relation != nil {
			continue
		}
		columnNames = append(columnNames, r.dialect.QuoteIdentifier(field.DBName))
	}

	// Prepare values for each entity
	for i, entity := range entities {
		values, placeholders := r.extractValues(entity, skipPK)
		placeholderGroups = append(placeholderGroups, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		allValues = append(allValues, values...)

		// Batch inserts in chunks of 1000 to avoid too many parameters
		if (i+1)%1000 == 0 || i == len(entities)-1 {
			query := fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES %s",
				r.dialect.QuoteIdentifier(r.metadata.TableName),
				strings.Join(columnNames, ", "),
				strings.Join(placeholderGroups, ", "),
			)

			_, err = tx.Exec(query, allValues...)
			if err != nil {
				return err
			}

			// Reset for next batch if needed
			placeholderGroups = placeholderGroups[:0]
			allValues = allValues[:0]
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// extractValues extracts values from an entity for bulk insert
func (r *BulkRepository[T]) extractValues(entity T, skipPK bool) ([]interface{}, []string) {
	var values []interface{}
	var placeholders []string

	entityVal := schema.GetEntityValue(entity)
	placeholderIndex := 0

	for _, field := range r.metadata.Fields {
		// Skip primary key if auto-increment and relation fields
		if (skipPK && field.IsPrimaryKey) || field.Relation != nil {
			continue
		}

		fieldVal := entityVal.FieldByName(field.Name)
		values = append(values, fieldVal.Interface())
		placeholders = append(placeholders, r.dialect.Placeholder(placeholderIndex))
		placeholderIndex++
	}

	return values, placeholders
}

// BulkUpdate efficiently updates multiple entities in a single transaction
func (r *BulkRepository[T]) BulkUpdate(entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the update statement
	var setColumns []string
	for _, field := range r.metadata.Fields {
		if field.IsPrimaryKey || field.Relation != nil {
			continue
		}
		setColumns = append(setColumns, fmt.Sprintf("%s = ?", r.dialect.QuoteIdentifier(field.DBName)))
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		strings.Join(setColumns, ", "),
		r.dialect.QuoteIdentifier(r.metadata.PrimaryKey.DBName),
	)

	// Prepare the statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute update for each entity
	for _, entity := range entities {
		entityVal := schema.GetEntityValue(entity)
		var values []interface{}

		// Collect values for SET clause
		for _, field := range r.metadata.Fields {
			if field.IsPrimaryKey || field.Relation != nil {
				continue
			}
			fieldVal := entityVal.FieldByName(field.Name)
			values = append(values, fieldVal.Interface())
		}

		// Add ID for WHERE clause
		idVal := entityVal.FieldByName(r.metadata.PrimaryKey.Name)
		values = append(values, idVal.Interface())

		// Execute update
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// BulkDelete efficiently deletes multiple entities by their IDs
func (r *BulkRepository[T]) BulkDelete(ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create placeholders for the IN clause
	var placeholders []string
	for i := 0; i < len(ids); i++ {
		placeholders = append(placeholders, r.dialect.Placeholder(i))
	}

	// Build the delete query
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s IN (%s)",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		r.dialect.QuoteIdentifier(r.metadata.PrimaryKey.DBName),
		strings.Join(placeholders, ", "),
	)

	// Execute the delete
	_, err = tx.Exec(query, ids...)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GenerateRandomUser creates a random user for testing
func GenerateRandomUser(index int) BulkUser {
	return BulkUser{
		Username:  fmt.Sprintf("user%d", index),
		Email:     fmt.Sprintf("user%d@example.com", index),
		Active:    rand.Intn(2) == 1,
		CreatedAt: time.Now(),
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Open SQLite database
	db, err := sql.Open("sqlite3", "./bulk_operations.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(BulkUser{}); err != nil {
		log.Fatalf("Failed to register BulkUser entity: %v", err)
	}

	// Get entity metadata
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(BulkUser{}))

	// Create table
	userSQL := sqliteDialect.CreateTableSQL(userMeta)
	_, err = db.Exec(userSQL)
	if err != nil {
		log.Fatalf("Failed to create bulk_users table: %v", err)
	}

	// Create bulk repository
	userRepo := NewBulkRepository[BulkUser](db, sqliteDialect)

	fmt.Println("=== Bulk Operations Example ===")

	// Generate users for bulk insert
	const userCount = 10000
	var users []BulkUser
	for i := 0; i < userCount; i++ {
		users = append(users, GenerateRandomUser(i))
	}

	// Measure time for bulk insert
	fmt.Printf("Inserting %d users in bulk...\n", userCount)
	startTime := time.Now()
	if err := userRepo.BulkInsert(users); err != nil {
		log.Fatalf("Bulk insert failed: %v", err)
	}
	insertDuration := time.Since(startTime)
	fmt.Printf("Bulk insert completed in %v\n", insertDuration)

	// Count users to verify
	count, err := userRepo.Find().Count()
	if err != nil {
		log.Fatalf("Count query failed: %v", err)
	}
	fmt.Printf("Total users in database: %d\n", count)

	// Fetch some users for bulk update
	fmt.Println("\nFetching users for bulk update...")
	fetchedUsers, err := userRepo.Find().Limit(100).All()
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	// Modify users
	for i := range fetchedUsers {
		fetchedUsers[i].Username = fmt.Sprintf("updated_user%d", i)
		fetchedUsers[i].Email = fmt.Sprintf("updated_user%d@example.com", i)
		fetchedUsers[i].Active = !fetchedUsers[i].Active
	}

	// Perform bulk update
	fmt.Printf("Updating %d users in bulk...\n", len(fetchedUsers))
	startTime = time.Now()
	if err := userRepo.BulkUpdate(fetchedUsers); err != nil {
		log.Fatalf("Bulk update failed: %v", err)
	}
	updateDuration := time.Since(startTime)
	fmt.Printf("Bulk update completed in %v\n", updateDuration)

	// Verify a few updated users
	verifyUser, err := userRepo.FindByID(1)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("\nVerified updated user: %s (%s)\n", verifyUser.Username, verifyUser.Email)

	// Demonstrate bulk delete
	var idsToDelete []interface{}
	for i := 5000; i < 6000; i++ {
		idsToDelete = append(idsToDelete, uint(i))
	}

	fmt.Printf("\nDeleting %d users in bulk...\n", len(idsToDelete))
	startTime = time.Now()
	if err := userRepo.BulkDelete(idsToDelete); err != nil {
		log.Fatalf("Bulk delete failed: %v", err)
	}
	deleteDuration := time.Since(startTime)
	fmt.Printf("Bulk delete completed in %v\n", deleteDuration)

	// Count users after delete to verify
	countAfterDelete, err := userRepo.Find().Count()
	if err != nil {
		log.Fatalf("Count query failed: %v", err)
	}
	fmt.Printf("Users remaining in database: %d (deleted %d users)\n", 
		countAfterDelete, count-countAfterDelete)

	// Compare with individual operations
	fmt.Println("\n=== Performance Comparison ===")
	fmt.Printf("Bulk insert of %d users: %v\n", userCount, insertDuration)
	fmt.Printf("Bulk update of %d users: %v\n", len(fetchedUsers), updateDuration)
	fmt.Printf("Bulk delete of %d users: %v\n", len(idsToDelete), deleteDuration)

	// Show estimate for individual operations
	estInsertTime := time.Duration(float64(insertDuration) * 50) // Estimate: bulk is ~50x faster
	estUpdateTime := time.Duration(float64(updateDuration) * 20) // Estimate: bulk is ~20x faster
	estDeleteTime := time.Duration(float64(deleteDuration) * 30) // Estimate: bulk is ~30x faster

	fmt.Printf("\nEstimated time for individual operations:\n")
	fmt.Printf("Individual insert of %d users: ~%v (est.)\n", userCount, estInsertTime)
	fmt.Printf("Individual update of %d users: ~%v (est.)\n", len(fetchedUsers), estUpdateTime)
	fmt.Printf("Individual delete of %d users: ~%v (est.)\n", len(idsToDelete), estDeleteTime)
}