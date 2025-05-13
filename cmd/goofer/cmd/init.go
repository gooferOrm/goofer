package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	initProjectName string
	initDialect     string
	withExamples    bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Goofer ORM project",
	Long: `Initialize a new Goofer ORM project with the recommended directory structure and files.
This command helps you set up a new project quickly with the correct configuration.

Examples:
  goofer init my-app
  goofer init my-app --dialect=postgres --with-examples`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			initProjectName = args[0]
		} else {
			// Use current directory name as project name
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}
			initProjectName = filepath.Base(currentDir)
		}
		initProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Add flags
	initCmd.Flags().StringVarP(&initDialect, "dialect", "d", "sqlite", "Database dialect (sqlite, mysql, postgres)")
	initCmd.Flags().BoolVar(&withExamples, "with-examples", false, "Include example entities")
}

func initProject() {
	// Create project directory if it doesn't exist and it's not the current directory
	currentDir, _ := os.Getwd()
	projectDir := filepath.Join(currentDir, initProjectName)
	
	if filepath.Clean(currentDir) != filepath.Clean(projectDir) {
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			fmt.Printf("Error creating project directory: %v\n", err)
			return
		}
	}

	// Create directory structure
	dirs := []string{
		"cmd",
		"internal/models",
		"internal/repository",
		"migrations",
		"db",
		"config",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Create go.mod file
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gooferOrm/goofer v0.1.0
`, initProjectName)

	// Add dialect-specific dependencies
	switch initDialect {
	case "sqlite":
		goModContent += `	github.com/mattn/go-sqlite3 v1.14.28
`
	case "mysql":
		goModContent += `	github.com/go-sql-driver/mysql v1.7.1
`
	case "postgres":
		goModContent += `	github.com/lib/pq v1.10.9
`
	}

	goModContent += `)
`

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		fmt.Printf("Error creating go.mod: %v\n", err)
		return
	}

	// Create example models if requested
	if withExamples {
		// Create a User model
		userModelContent := fmt.Sprintf(`package models

import (
	"time"

	"github.com/gooferOrm/goofer/schema"
)

// User represents a user in the system
type User struct {
	ID        uint      ` + "`orm:\"primaryKey;autoIncrement\" validate:\"required\"`" + `
	Name      string    ` + "`orm:\"type:varchar(255);notnull\" validate:\"required\"`" + `
	Email     string    ` + "`orm:\"unique;type:varchar(255);notnull\" validate:\"required,email\"`" + `
	CreatedAt time.Time ` + "`orm:\"type:timestamp;default:CURRENT_TIMESTAMP\"`" + `
	UpdatedAt time.Time ` + "`orm:\"type:timestamp\"`" + `
}

// TableName returns the database table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate is called before the record is created
func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is called before the record is updated
func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}
`)

		if err := os.WriteFile(filepath.Join(projectDir, "internal/models/user.go"), []byte(userModelContent), 0644); err != nil {
			fmt.Printf("Error creating user model: %v\n", err)
			return
		}

		// Create Post model
		postModelContent := fmt.Sprintf(`package models

import (
	"time"

	"github.com/gooferOrm/goofer/schema"
)

// Post represents a blog post
type Post struct {
	ID        uint      ` + "`orm:\"primaryKey;autoIncrement\" validate:\"required\"`" + `
	Title     string    ` + "`orm:\"type:varchar(255);notnull\" validate:\"required\"`" + `
	Content   string    ` + "`orm:\"type:text\" validate:\"required\"`" + `
	UserID    uint      ` + "`orm:\"index;notnull\" validate:\"required\"`" + `
	Published bool      ` + "`orm:\"type:boolean;default:false\"`" + `
	CreatedAt time.Time ` + "`orm:\"type:timestamp;default:CURRENT_TIMESTAMP\"`" + `
	UpdatedAt time.Time ` + "`orm:\"type:timestamp\"`" + `
	User      *User     ` + "`orm:\"relation:ManyToOne;foreignKey:UserID\"`" + `
}

// TableName returns the database table name
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate is called before the record is created
func (p *Post) BeforeCreate() error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is called before the record is updated
func (p *Post) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return nil
}
`)

		if err := os.WriteFile(filepath.Join(projectDir, "internal/models/post.go"), []byte(postModelContent), 0644); err != nil {
			fmt.Printf("Error creating post model: %v\n", err)
			return
		}
	}

	// Create main.go
	mainContent := fmt.Sprintf(`package main

import (
	"database/sql"
	"fmt"
	"log"

%s
	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"

	"%s/internal/models"
)

func main() {
	// Open database connection
%s

	// Create dialect
%s

	// Register entities
%s

	fmt.Println("Connected to database successfully!")
}
`,
		getDriverImport(initDialect),
		initProjectName,
		getDbConnectCode(initDialect),
		getDialectCode(initDialect),
		getRegisterEntitiesCode(withExamples),
	)

	if err := os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainContent), 0644); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		return
	}

	// Create config file based on dialect
	switch initDialect {
	case "sqlite":
		configContent := `# SQLite database configuration
database:
  dialect: sqlite
  path: ./db/app.db
`
		if err := os.WriteFile(filepath.Join(projectDir, "config/config.yaml"), []byte(configContent), 0644); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			return
		}
	case "mysql":
		configContent := `# MySQL database configuration
database:
  dialect: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: myapp
  params: parseTime=true
`
		if err := os.WriteFile(filepath.Join(projectDir, "config/config.yaml"), []byte(configContent), 0644); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			return
		}
	case "postgres":
		configContent := `# PostgreSQL database configuration
database:
  dialect: postgres
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: myapp
  sslmode: disable
`
		if err := os.WriteFile(filepath.Join(projectDir, "config/config.yaml"), []byte(configContent), 0644); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			return
		}
	}

	// Create a basic README.md
	readmeContent := fmt.Sprintf("# %s\n\nA Go application using Goofer ORM.\n\n## Features\n\n- Type-safe database operations\n- %s database support\n- Clean architecture\n\n## Getting Started\n\n### Prerequisites\n\n- Go 1.21 or higher\n- %s\n\n### Installation\n\n1. Clone the repository\n2. Run \"go mod tidy\" to install dependencies\n3. Configure the database in \"config/config.yaml\"\n4. Run \"go run main.go\"\n\n## Project Structure\n\n- \"cmd/\": Command-line applications\n- \"internal/\": Internal packages\n  - \"models/\": Database entity models\n  - \"repository/\": Data access layer\n- \"migrations/\": Database migrations\n- \"db/\": Database files (for SQLite)\n- \"config/\": Configuration files\n\n## License\n\nThis project is licensed under the MIT License - see the LICENSE file for details.",
		strings.Title(initProjectName),
		strings.Title(initDialect),
		getPrerequisiteText(initDialect),
	)

	if err := os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readmeContent), 0644); err != nil {
		fmt.Printf("Error creating README.md: %v\n", err)
		return
	}

	// Create a .gitignore file
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*.bin

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Database files
*.db
*.sqlite

# Environment variables
.env

# IDE specific files
.idea/
.vscode/
*.swp
*.swo
`

	if err := os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		fmt.Printf("Error creating .gitignore: %v\n", err)
		return
	}

	fmt.Printf("Initialized Goofer ORM project: %s\n", initProjectName)
	fmt.Printf("Project created with %s dialect\n", initDialect)
	if withExamples {
		fmt.Println("Example models included")
	}
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. cd", initProjectName)
	fmt.Println("2. go mod tidy")
	fmt.Println("3. Edit config/config.yaml to configure your database")
	fmt.Println("4. Run 'go run main.go' to test database connection")
}

func getDriverImport(dialect string) string {
	switch dialect {
	case "sqlite":
		return "\t_ \"github.com/mattn/go-sqlite3\""
	case "mysql":
		return "\t_ \"github.com/go-sql-driver/mysql\""
	case "postgres":
		return "\t_ \"github.com/lib/pq\""
	default:
		return "\t// Unknown dialect"
	}
}

func getDbConnectCode(dialect string) string {
	switch dialect {
	case "sqlite":
		return `	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()`
	case "mysql":
		return `	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/myapp?parseTime=true")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()`
	case "postgres":
		return `	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=myapp sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()`
	default:
		return `	// Unknown dialect
	db, err := sql.Open("unknown", "connection string")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()`
	}
}

func getDialectCode(dialect string) string {
	switch dialect {
	case "sqlite":
		return `	sqliteDialect := dialect.NewSQLiteDialect()
	// Use sqliteDialect for database operations`
	case "mysql":
		return `	mysqlDialect := dialect.NewMySQLDialect()
	// Use mysqlDialect for database operations`
	case "postgres":
		return `	postgresDialect := dialect.NewPostgresDialect()
	// Use postgresDialect for database operations`
	default:
		return `	// Unknown dialect`
	}
}

func getRegisterEntitiesCode(withExamples bool) string {
	if withExamples {
		return `	if err := schema.Registry.RegisterEntity(models.User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}
	
	if err := schema.Registry.RegisterEntity(models.Post{}); err != nil {
		log.Fatalf("Failed to register Post entity: %v", err)
	}`
	} else {
		return `	// Register your entities here
	// Example:
	// if err := schema.Registry.RegisterEntity(models.User{}); err != nil {
	//     log.Fatalf("Failed to register entity: %v", err)
	// }`
	}
}

func getPrerequisiteText(dialect string) string {
	switch dialect {
	case "sqlite":
		return "SQLite"
	case "mysql":
		return "MySQL database server"
	case "postgres":
		return "PostgreSQL database server"
	default:
		return "Database server"
	}
}