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

// SoftDeleteEntity implements soft delete functionality
type SoftDeleteEntity interface {
	schema.Entity
	IsDeleted() bool
	MarkAsDeleted()
	MarkAsActive()
}

// Task entity with soft delete capability
type Task struct {
	ID          uint       `orm:"primaryKey;autoIncrement"`
	Title       string     `orm:"type:varchar(255);notnull"`
	Description string     `orm:"type:text"`
	Completed   bool       `orm:"type:boolean;default:false"`
	DueDate     *time.Time `orm:"type:timestamp"`
	CreatedAt   time.Time  `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `orm:"type:timestamp"`
	DeletedAt   *time.Time `orm:"type:timestamp"`
}

// TableName returns the table name for the Task entity
func (Task) TableName() string {
	return "tasks"
}

// BeforeSave hook to set the UpdatedAt field
func (t *Task) BeforeSave() error {
	t.UpdatedAt = time.Now()
	return nil
}

// IsDeleted returns true if the task is soft-deleted
func (t *Task) IsDeleted() bool {
	return t.DeletedAt != nil
}

// MarkAsDeleted sets the DeletedAt field to the current time
func (t *Task) MarkAsDeleted() {
	now := time.Now()
	t.DeletedAt = &now
}

// MarkAsActive clears the DeletedAt field
func (t *Task) MarkAsActive() {
	t.DeletedAt = nil
}

// SoftDeleteRepository extends the base repository with soft delete functionality
type SoftDeleteRepository[T SoftDeleteEntity] struct {
	*repository.Repository[T]
	IncludeDeleted bool
}

// NewSoftDeleteRepository creates a new repository with soft delete support
func NewSoftDeleteRepository[T SoftDeleteEntity](db *sql.DB, dialect dialect.Dialect) *SoftDeleteRepository[T] {
	return &SoftDeleteRepository[T]{
		Repository:     repository.NewRepository[T](db, dialect),
		IncludeDeleted: false,
	}
}

// WithIncludeDeleted sets whether to include soft-deleted entities in queries
func (r *SoftDeleteRepository[T]) WithIncludeDeleted(include bool) *SoftDeleteRepository[T] {
	r.IncludeDeleted = include
	return r
}

// Find overrides the base Find method to exclude soft-deleted entities
func (r *SoftDeleteRepository[T]) Find() *SoftDeleteQueryBuilder[T] {
	return &SoftDeleteQueryBuilder[T]{
		QueryBuilder:   r.Repository.Find(),
		IncludeDeleted: r.IncludeDeleted,
	}
}

// FindByID overrides the base FindByID method to respect soft-delete
func (r *SoftDeleteRepository[T]) FindByID(id interface{}) (*T, error) {
	entity, err := r.Repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !r.IncludeDeleted && entity.IsDeleted() {
		return nil, sql.ErrNoRows
	}

	return entity, nil
}

// Delete performs a soft delete instead of a hard delete
func (r *SoftDeleteRepository[T]) Delete(entity *T) error {
	(*entity).MarkAsDeleted()
	return r.Repository.Save(entity)
}

// HardDelete actually removes the entity from the database
func (r *SoftDeleteRepository[T]) HardDelete(entity *T) error {
	return r.Repository.Delete(entity)
}

// Restore removes the soft-delete marker
func (r *SoftDeleteRepository[T]) Restore(entity *T) error {
	(*entity).MarkAsActive()
	return r.Repository.Save(entity)
}

// SoftDeleteQueryBuilder extends the query builder for soft delete repositories
type SoftDeleteQueryBuilder[T SoftDeleteEntity] struct {
	*repository.QueryBuilder[T]
	IncludeDeleted bool
}

// All returns all entities, respecting the IncludeDeleted flag
func (qb *SoftDeleteQueryBuilder[T]) All() ([]T, error) {
	if !qb.IncludeDeleted {
		qb.Where("deleted_at IS NULL")
	}
	return qb.QueryBuilder.All()
}

// One returns a single entity, respecting the IncludeDeleted flag
func (qb *SoftDeleteQueryBuilder[T]) One() (*T, error) {
	if !qb.IncludeDeleted {
		qb.Where("deleted_at IS NULL")
	}
	return qb.QueryBuilder.One()
}

// Count returns the count of entities, respecting the IncludeDeleted flag
func (qb *SoftDeleteQueryBuilder[T]) Count() (int64, error) {
	if !qb.IncludeDeleted {
		qb.Where("deleted_at IS NULL")
	}
	return qb.QueryBuilder.Count()
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./soft_delete.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(Task{}); err != nil {
		log.Fatalf("Failed to register Task entity: %v", err)
	}

	// Get entity metadata
	taskMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Task{}))

	// Create table
	taskSQL := sqliteDialect.CreateTableSQL(taskMeta)
	_, err = db.Exec(taskSQL)
	if err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}

	// Create repository with soft delete support
	taskRepo := NewSoftDeleteRepository[Task](db, sqliteDialect)

	fmt.Println("=== Soft Delete Example ===")

	// Create some tasks
	dueDate1 := time.Now().AddDate(0, 0, 7) // 7 days from now
	dueDate2 := time.Now().AddDate(0, 0, 14) // 14 days from now

	tasks := []Task{
		{
			Title:       "Implement soft delete",
			Description: "Add soft delete functionality to the ORM",
			DueDate:     &dueDate1,
		},
		{
			Title:       "Write documentation",
			Description: "Document soft delete feature",
			DueDate:     &dueDate2,
		},
		{
			Title:       "Create examples",
			Description: "Create example code for soft delete",
		},
	}

	// Save tasks
	for i := range tasks {
		if err := taskRepo.Save(&tasks[i]); err != nil {
			log.Fatalf("Failed to save task: %v", err)
		}
		fmt.Printf("Created task: %s (ID: %d)\n", tasks[i].Title, tasks[i].ID)
	}

	// Soft delete one task
	fmt.Println("\n--- Soft Deleting a Task ---")
	taskToDelete := tasks[1] // "Write documentation" task
	if err := taskRepo.Delete(&taskToDelete); err != nil {
		log.Fatalf("Failed to soft delete task: %v", err)
	}
	fmt.Printf("Soft-deleted task: %s (ID: %d)\n", taskToDelete.Title, taskToDelete.ID)
	fmt.Printf("Task DeletedAt: %v\n", *taskToDelete.DeletedAt)

	// List all non-deleted tasks
	fmt.Println("\n--- Active Tasks ---")
	activeTasks, err := taskRepo.Find().All()
	if err != nil {
		log.Fatalf("Failed to fetch active tasks: %v", err)
	}

	fmt.Printf("Found %d active tasks:\n", len(activeTasks))
	for _, t := range activeTasks {
		dueDateStr := "No due date"
		if t.DueDate != nil {
			dueDateStr = t.DueDate.Format("2006-01-02")
		}
		fmt.Printf("- %s: %s (Due: %s)\n", t.Title, t.Description, dueDateStr)
	}

	// Include deleted tasks in query
	fmt.Println("\n--- All Tasks (Including Deleted) ---")
	allTasks, err := taskRepo.WithIncludeDeleted(true).Find().All()
	if err != nil {
		log.Fatalf("Failed to fetch all tasks: %v", err)
	}

	fmt.Printf("Found %d total tasks:\n", len(allTasks))
	for _, t := range allTasks {
		status := "Active"
		if t.IsDeleted() {
			status = "Deleted"
		}
		
		dueDateStr := "No due date"
		if t.DueDate != nil {
			dueDateStr = t.DueDate.Format("2006-01-02")
		}
		
		fmt.Printf("- [%s] %s: %s (Due: %s)\n", status, t.Title, t.Description, dueDateStr)
	}

	// Restore a deleted task
	fmt.Println("\n--- Restoring a Deleted Task ---")
	if err := taskRepo.Restore(&taskToDelete); err != nil {
		log.Fatalf("Failed to restore task: %v", err)
	}
	fmt.Printf("Restored task: %s (ID: %d)\n", taskToDelete.Title, taskToDelete.ID)

	// Verify task is no longer marked as deleted
	restoredTask, err := taskRepo.FindByID(taskToDelete.ID)
	if err != nil {
		log.Fatalf("Failed to find restored task: %v", err)
	}
	fmt.Printf("Task DeletedAt is now: %v\n", restoredTask.DeletedAt)

	// Hard delete a task (actually remove from database)
	fmt.Println("\n--- Hard Deleting a Task ---")
	taskToHardDelete := tasks[2] // "Create examples" task
	if err := taskRepo.HardDelete(&taskToHardDelete); err != nil {
		log.Fatalf("Failed to hard delete task: %v", err)
	}
	fmt.Printf("Hard-deleted task: %s (ID: %d)\n", taskToHardDelete.Title, taskToHardDelete.ID)

	// Try to find the hard-deleted task (should fail)
	_, err = taskRepo.WithIncludeDeleted(true).FindByID(taskToHardDelete.ID)
	if err == nil {
		log.Fatal("Expected error when finding hard-deleted task, but got none")
	}
	fmt.Printf("As expected, couldn't find the hard-deleted task: %v\n", err)

	// Count active vs. all tasks
	activeCount, err := taskRepo.Find().Count()
	if err != nil {
		log.Fatalf("Failed to count active tasks: %v", err)
	}

	totalCount, err := taskRepo.WithIncludeDeleted(true).Find().Count()
	if err != nil {
		log.Fatalf("Failed to count all tasks: %v", err)
	}

	fmt.Printf("\nActive tasks: %d, Total tasks (including deleted): %d\n", activeCount, totalCount)
}