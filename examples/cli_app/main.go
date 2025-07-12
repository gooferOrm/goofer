package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/engine"
	"github.com/gooferOrm/goofer/repository"
)

// User entity
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

// Post entity
type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Title     string    `orm:"type:varchar(255);notnull"`
	Content   string    `orm:"type:text"`
	UserID    uint      `orm:"index;notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}

func (Post) TableName() string {
	return "posts"
}

// Global variables for the CLI
var (
	dbPath   = "./cli_app.db"
	client   *engine.Client
	userRepo *repository.Repository[User]
	postRepo *repository.Repository[Post]
)

var rootCmd = &cobra.Command{
	Use:   "blog-cli",
	Short: "A CLI application for managing a blog with Goofer ORM",
	Long:  `A command-line interface for managing users and posts using Goofer ORM.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database and create tables",
	Run: func(cmd *cobra.Command, args []string) {
		initDatabase()
		fmt.Println("Database initialized successfully!")
	},
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var createUserCmd = &cobra.Command{
	Use:   "create [name] [email]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		email := args[1]

		user := &User{
			Name:  name,
			Email: email,
		}

		if err := userRepo.Save(user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		fmt.Printf("Created user with ID: %d\n", user.ID)
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		users, err := userRepo.Find().All()
		if err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}

		if len(users) == 0 {
			fmt.Println("No users found.")
			return
		}

		fmt.Println("Users:")
		fmt.Printf("%-5s %-20s %-30s %-20s\n", "ID", "Name", "Email", "Created At")
		fmt.Println(string(make([]byte, 80, 80)))
		for _, user := range users {
			fmt.Printf("%-5d %-20s %-30s %-20s\n",
				user.ID,
				user.Name,
				user.Email,
				user.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid user ID: %v", err)
		}

		user, err := userRepo.FindByID(uint(id))
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}

		fmt.Printf("User ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Created At: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Manage posts",
}

var createPostCmd = &cobra.Command{
	Use:   "create [user-id] [title] [content]",
	Short: "Create a new post",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid user ID: %v", err)
		}

		title := args[1]
		content := args[2]

		post := &Post{
			Title:   title,
			Content: content,
			UserID:  uint(userID),
		}

		if err := postRepo.Save(post); err != nil {
			log.Fatalf("Failed to create post: %v", err)
		}

		fmt.Printf("Created post with ID: %d\n", post.ID)
	},
}

var listPostsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all posts",
	Run: func(cmd *cobra.Command, args []string) {
		posts, err := postRepo.Find().All()
		if err != nil {
			log.Fatalf("Failed to list posts: %v", err)
		}

		if len(posts) == 0 {
			fmt.Println("No posts found.")
			return
		}

		fmt.Println("Posts:")
		fmt.Printf("%-5s %-30s %-10s %-20s\n", "ID", "Title", "User ID", "Created At")
		fmt.Println(string(make([]byte, 70, 70)))
		for _, post := range posts {
			fmt.Printf("%-5d %-30s %-10d %-20s\n",
				post.ID,
				post.Title,
				post.UserID,
				post.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var getPostCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a post by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid post ID: %v", err)
		}

		post, err := postRepo.FindByID(uint(id))
		if err != nil {
			log.Fatalf("Failed to get post: %v", err)
		}

		fmt.Printf("Post ID: %d\n", post.ID)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Content: %s\n", post.Content)
		fmt.Printf("User ID: %d\n", post.UserID)
		fmt.Printf("Created At: %s\n", post.CreatedAt.Format("2006-01-02 15:04:05"))
	},
}

var searchPostsCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search posts by title or content",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]

		posts, err := postRepo.Find().
			Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			All()
		if err != nil {
			log.Fatalf("Failed to search posts: %v", err)
		}

		if len(posts) == 0 {
			fmt.Printf("No posts found matching '%s'\n", keyword)
			return
		}

		fmt.Printf("Found %d posts matching '%s':\n", len(posts), keyword)
		for _, post := range posts {
			fmt.Printf("- [%d] %s (User: %d)\n", post.ID, post.Title, post.UserID)
		}
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show database statistics",
	Run: func(cmd *cobra.Command, args []string) {
		userCount, err := userRepo.Find().Count()
		if err != nil {
			log.Fatalf("Failed to count users: %v", err)
		}

		postCount, err := postRepo.Find().Count()
		if err != nil {
			log.Fatalf("Failed to count posts: %v", err)
		}

		fmt.Println("Database Statistics:")
		fmt.Printf("Total Users: %d\n", userCount)
		fmt.Printf("Total Posts: %d\n", postCount)
	},
}

func init() {
	// Add commands to root
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(statsCmd)

	// Add user subcommands
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(listUsersCmd)
	userCmd.AddCommand(getUserCmd)

	// Add post subcommands
	postCmd.AddCommand(createPostCmd)
	postCmd.AddCommand(listPostsCmd)
	postCmd.AddCommand(getPostCmd)
	postCmd.AddCommand(searchPostsCmd)
}

func initDatabase() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Create client with auto-migration
	client, err = engine.NewClient(db, sqliteDialect, &User{}, &Post{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create repositories
	userRepo = repository.NewRepository[User](db, sqliteDialect)
	postRepo = repository.NewRepository[Post](db, sqliteDialect)
}

func main() {
	// Initialize database connection
	initDatabase()
	defer client.Close()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
