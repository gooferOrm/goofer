package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/engine"
	"github.com/gooferOrm/goofer/repository"
)

// Task entity for a simple todo CLI
type Task struct {
	ID          uint      `orm:"primaryKey;autoIncrement"`
	Title       string    `orm:"type:varchar(255);notnull"`
	Description string    `orm:"type:text"`
	Completed   bool      `orm:"type:boolean;default:false"`
	CreatedAt   time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Task) TableName() string {
	return "tasks"
}

// Global variables
var (
	db       *sql.DB
	client   *engine.Client
	taskRepo *repository.Repository[Task]
)

func main() {
	// Initialize database
	initDatabase()
	defer db.Close()

	fmt.Println("=== Simple Task Manager CLI ===")
	fmt.Println("Commands: add, list, complete, delete, quit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "add":
			handleAdd(args)
		case "list":
			handleList()
		case "complete":
			handleComplete(args)
		case "delete":
			handleDelete(args)
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			showHelp()
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
		fmt.Println()
	}
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Create client with auto-migration
	client, err = engine.NewClient(db, sqliteDialect, &Task{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create repository
	taskRepo = repository.NewRepository[Task](db, sqliteDialect)
}

func handleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: add <title> [description]")
		return
	}

	title := args[0]
	description := ""
	if len(args) > 1 {
		description = strings.Join(args[1:], " ")
	}

	task := &Task{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := taskRepo.Save(task); err != nil {
		fmt.Printf("Error creating task: %v\n", err)
		return
	}

	fmt.Printf("Created task with ID: %d\n", task.ID)
}

func handleList() {
	tasks, err := taskRepo.Find().All()
	if err != nil {
		fmt.Printf("Error listing tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("Tasks:")
	fmt.Printf("%-5s %-8s %-30s %s\n", "ID", "Status", "Title", "Description")
	fmt.Println(strings.Repeat("-", 80))

	for _, task := range tasks {
		status := "Pending"
		if task.Completed {
			status = "Done"
		}

		// Truncate description if too long
		desc := task.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		fmt.Printf("%-5d %-8s %-30s %s\n", task.ID, status, task.Title, desc)
	}
}

func handleComplete(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: complete <task-id>")
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	task, err := taskRepo.FindByID(uint(id))
	if err != nil {
		fmt.Printf("Task not found: %v\n", err)
		return
	}

	task.Completed = true
	if err := taskRepo.Save(task); err != nil {
		fmt.Printf("Error updating task: %v\n", err)
		return
	}

	fmt.Printf("Marked task %d as completed\n", task.ID)
}

func handleDelete(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: delete <task-id>")
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	task, err := taskRepo.FindByID(uint(id))
	if err != nil {
		fmt.Printf("Task not found: %v\n", err)
		return
	}

	if err := taskRepo.Delete(task); err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		return
	}

	fmt.Printf("Deleted task %d\n", task.ID)
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  add <title> [description]  - Add a new task")
	fmt.Println("  list                       - List all tasks")
	fmt.Println("  complete <task-id>         - Mark a task as completed")
	fmt.Println("  delete <task-id>           - Delete a task")
	fmt.Println("  help                       - Show this help")
	fmt.Println("  quit                       - Exit the application")
}

// Example usage of advanced Goofer ORM features
func demonstrateAdvancedFeatures() {
	fmt.Println("\n=== Advanced Features Demo ===")

	// 1. Find with conditions
	fmt.Println("1. Finding incomplete tasks:")
	incompleteTasks, err := taskRepo.Find().Where("completed = ?", false).All()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d incomplete tasks\n", len(incompleteTasks))
	}

	// 2. Count records
	fmt.Println("2. Counting tasks:")
	totalCount, err := taskRepo.Find().Count()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Total tasks: %d\n", totalCount)
	}

	// 3. Find with limit and offset
	fmt.Println("3. Finding first 3 tasks:")
	firstThree, err := taskRepo.Find().Limit(3).All()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("First 3 tasks: %d found\n", len(firstThree))
	}

	// 4. Transaction example
	fmt.Println("4. Transaction example:")
	err = taskRepo.Transaction(func(txRepo *repository.Repository[Task]) error {
		// Create multiple tasks in a transaction
		task1 := &Task{Title: "Transaction Task 1", Description: "Created in transaction"}
		task2 := &Task{Title: "Transaction Task 2", Description: "Also in transaction"}

		if err := txRepo.Save(task1); err != nil {
			return err
		}
		if err := txRepo.Save(task2); err != nil {
			return err
		}

		fmt.Printf("Created tasks %d and %d in transaction\n", task1.ID, task2.ID)
		return nil
	})

	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}
}
