# Working with Relationships

Goofer ORM makes it easy to work with entity relationships. This example demonstrates how to model and work with different types of relationships: one-to-one, one-to-many, and many-to-many.

## Schema Definition

First, let's define our entities with various relationships:

```go
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

// User entity with multiple relationship types
type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Name      string    `orm:"type:varchar(255);notnull"`
	Email     string    `orm:"unique;type:varchar(255);notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	
	// One-to-One relationship: User has one Profile
	Profile   *Profile  `orm:"relation:OneToOne;foreignKey:UserID"`
	
	// One-to-Many relationship: User has many Posts
	Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
	
	// Many-to-Many relationship: User has many Roles through UserRoles
	Roles     []Role    `orm:"relation:ManyToMany;joinTable:user_roles;foreignKey:UserID;referenceKey:RoleID"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// Profile entity for One-to-One relationship
type Profile struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	UserID    uint      `orm:"unique;notnull"`
	Bio       string    `orm:"type:text"`
	Avatar    string    `orm:"type:varchar(255)"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	
	// Reference back to User
	User      *User     `orm:"relation:OneToOne;foreignKey:UserID"`
}

// TableName returns the table name for the Profile entity
func (Profile) TableName() string {
	return "profiles"
}

// Post entity for One-to-Many relationship
type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Title     string    `orm:"type:varchar(255);notnull"`
	Content   string    `orm:"type:text"`
	UserID    uint      `orm:"index;notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	
	// Reference back to User
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
	
	// One-to-Many relationship: Post has many Comments
	Comments  []Comment `orm:"relation:OneToMany;foreignKey:PostID"`
}

// TableName returns the table name for the Post entity
func (Post) TableName() string {
	return "posts"
}

// Comment entity for One-to-Many relationship with Post
type Comment struct {
	ID        uint      `orm:"primaryKey;autoIncrement"`
	Content   string    `orm:"type:text;notnull"`
	PostID    uint      `orm:"index;notnull"`
	UserID    uint      `orm:"index;notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	
	// Reference back to Post
	Post      *Post     `orm:"relation:ManyToOne;foreignKey:PostID"`
	
	// Reference back to User
	User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}

// TableName returns the table name for the Comment entity
func (Comment) TableName() string {
	return "comments"
}

// Role entity for Many-to-Many relationship
type Role struct {
	ID          uint      `orm:"primaryKey;autoIncrement"`
	Name        string    `orm:"type:varchar(50);unique;notnull"`
	Description string    `orm:"type:text"`
	CreatedAt   time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	
	// Many-to-Many relationship: Role has many Users through UserRoles
	Users       []User    `orm:"relation:ManyToMany;joinTable:user_roles;foreignKey:RoleID;referenceKey:UserID"`
}

// TableName returns the table name for the Role entity
func (Role) TableName() string {
	return "roles"
}

// UserRole entity for Many-to-Many join table
type UserRole struct {
	UserID    uint      `orm:"primaryKey;notnull"`
	RoleID    uint      `orm:"primaryKey;notnull"`
	CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the UserRole entity
func (UserRole) TableName() string {
	return "user_roles"
}
```

## Setup and Table Creation

Let's set up our database and create tables for all entities:

```go
func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./relationships.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Profile{}); err != nil {
		log.Fatalf("Failed to register Profile entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Post{}); err != nil {
		log.Fatalf("Failed to register Post entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Comment{}); err != nil {
		log.Fatalf("Failed to register Comment entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Role{}); err != nil {
		log.Fatalf("Failed to register Role entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(UserRole{}); err != nil {
		log.Fatalf("Failed to register UserRole entity: %v", err)
	}

	// Create tables
	createTables(db, sqliteDialect)

	// Create repositories
	userRepo := repository.NewRepository[User](db, sqliteDialect)
	profileRepo := repository.NewRepository[Profile](db, sqliteDialect)
	postRepo := repository.NewRepository[Post](db, sqliteDialect)
	commentRepo := repository.NewRepository[Comment](db, sqliteDialect)
	roleRepo := repository.NewRepository[Role](db, sqliteDialect)
	userRoleRepo := repository.NewRepository[UserRole](db, sqliteDialect)
	
	// Rest of the code...
}

func createTables(db *sql.DB, sqliteDialect *dialect.SQLiteDialect) {
	// Get entity metadata
	userMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
	profileMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Profile{}))
	postMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Post{}))
	commentMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Comment{}))
	roleMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Role{}))
	userRoleMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(UserRole{}))

	// Create tables
	_, err := db.Exec(sqliteDialect.CreateTableSQL(userMeta))
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(profileMeta))
	if err != nil {
		log.Fatalf("Failed to create profiles table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(postMeta))
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(commentMeta))
	if err != nil {
		log.Fatalf("Failed to create comments table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(roleMeta))
	if err != nil {
		log.Fatalf("Failed to create roles table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(userRoleMeta))
	if err != nil {
		log.Fatalf("Failed to create user_roles table: %v", err)
	}
}
```

## Working with Relationships

### Creating a User with One-to-One Relationship (Profile)

```go
// Create a user
user := &User{
	Name:  "John Doe",
	Email: "john@example.com",
}

// Save the user
if err := userRepo.Save(user); err != nil {
	log.Fatalf("Failed to save user: %v", err)
}

fmt.Printf("Created user with ID: %d\n", user.ID)

// Create a profile (One-to-One relationship)
profile := &Profile{
	UserID: user.ID,
	Bio:    "Software developer and tech enthusiast",
	Avatar: "avatar.jpg",
}

// Save the profile
if err := profileRepo.Save(profile); err != nil {
	log.Fatalf("Failed to save profile: %v", err)
}

fmt.Printf("Created profile with ID: %d for user ID: %d\n", profile.ID, profile.UserID)
```

### Creating a Post (One-to-Many relationship)

```go
// Create a post (One-to-Many relationship)
post := &Post{
	Title:   "Understanding ORM Relationships",
	Content: "This post explains different types of ORM relationships...",
	UserID:  user.ID,
}

// Save the post
if err := postRepo.Save(post); err != nil {
	log.Fatalf("Failed to save post: %v", err)
}

fmt.Printf("Created post with ID: %d by user ID: %d\n", post.ID, post.UserID)
```

### Creating a Comment (One-to-Many relationship with Post)

```go
// Create a comment (One-to-Many relationship with Post)
comment := &Comment{
	Content: "Great post! Very informative.",
	PostID:  post.ID,
	UserID:  user.ID, // Self-comment for simplicity
}

// Save the comment
if err := commentRepo.Save(comment); err != nil {
	log.Fatalf("Failed to save comment: %v", err)
}

fmt.Printf("Created comment with ID: %d for post ID: %d\n", comment.ID, comment.PostID)
```

### Creating Many-to-Many Relationship (Users and Roles)

```go
// Create roles (for Many-to-Many relationship)
adminRole := &Role{
	Name:        "Admin",
	Description: "Administrator with full access",
}

// Save the admin role
if err := roleRepo.Save(adminRole); err != nil {
	log.Fatalf("Failed to save admin role: %v", err)
}

userRole := &Role{
	Name:        "User",
	Description: "Regular user with limited access",
}

// Save the user role
if err := roleRepo.Save(userRole); err != nil {
	log.Fatalf("Failed to save user role: %v", err)
}

fmt.Printf("Created roles: Admin (ID: %d), User (ID: %d)\n", adminRole.ID, userRole.ID)

// Assign roles to user (Many-to-Many relationship)
userAdminRole := &UserRole{
	UserID: user.ID,
	RoleID: adminRole.ID,
}

// Save the user-admin role assignment
if err := userRoleRepo.Save(userAdminRole); err != nil {
	log.Fatalf("Failed to assign admin role to user: %v", err)
}

userUserRole := &UserRole{
	UserID: user.ID,
	RoleID: userRole.ID,
}

// Save the user-user role assignment
if err := userRoleRepo.Save(userUserRole); err != nil {
	log.Fatalf("Failed to assign user role to user: %v", err)
}

fmt.Printf("Assigned roles to user ID: %d\n", user.ID)
```

## Querying Relationships

### Finding a User with Profile (One-to-One)

```go
// Find user with ID
foundUser, err := userRepo.FindByID(user.ID)
if err != nil {
	log.Fatalf("Failed to find user: %v", err)
}

fmt.Printf("\nFound user: %s (%s)\n", foundUser.Name, foundUser.Email)

// Find profile by user ID (One-to-One)
foundProfile, err := profileRepo.Find().
	Where("user_id = ?", user.ID).
	One()
if err != nil {
	log.Fatalf("Failed to find profile: %v", err)
}

fmt.Printf("User's profile: Bio: %s, Avatar: %s\n", foundProfile.Bio, foundProfile.Avatar)
```

### Finding Posts by User (One-to-Many)

```go
// Find posts by user ID (One-to-Many)
posts, err := postRepo.Find().
	Where("user_id = ?", user.ID).
	All()
if err != nil {
	log.Fatalf("Failed to find posts: %v", err)
}

fmt.Printf("User has %d posts:\n", len(posts))
for _, p := range posts {
	fmt.Printf("- %s: %s\n", p.Title, p.Content)

	// Find comments for this post (One-to-Many)
	comments, err := commentRepo.Find().Where("post_id = ?", p.ID).All()
	if err != nil {
		log.Fatalf("Failed to find comments: %v", err)
	}

	fmt.Printf("  Post has %d comments:\n", len(comments))
	for _, c := range comments {
		fmt.Printf("  - %s\n", c.Content)
	}
}
```

### Finding Roles for a User (Many-to-Many)

```go
// Find roles for user (Many-to-Many)
userRoles, err := userRoleRepo.Find().
	Where("user_id = ?", user.ID).
	All()
if err != nil {
	log.Fatalf("Failed to find user roles: %v", err)
}

fmt.Printf("User has %d roles:\n", len(userRoles))
for _, ur := range userRoles {
	role, err := roleRepo.FindByID(ur.RoleID)
	if err != nil {
		log.Fatalf("Failed to find role: %v", err)
	}
	fmt.Printf("- %s: %s\n", role.Name, role.Description)
}
```

## Filtering by Relationships

You can also use relationships to filter entities:

```go
// Find all posts with at least one comment
postsWithComments, err := postRepo.Find().
	Where("id IN (SELECT post_id FROM comments)").
	All()
if err != nil {
	log.Fatalf("Failed to find posts with comments: %v", err)
}

fmt.Printf("Found %d posts with comments\n", len(postsWithComments))

// Find all users with the Admin role
adminUsers, err := userRepo.Find().
	Where("id IN (SELECT user_id FROM user_roles WHERE role_id = ?)", adminRole.ID).
	All()
if err != nil {
	log.Fatalf("Failed to find admin users: %v", err)
}

fmt.Printf("Found %d users with Admin role\n", len(adminUsers))
```

## Updating Relationships

### Updating a One-to-One Relationship

```go
// Update profile
foundProfile.Bio = "Updated bio information"
if err := profileRepo.Save(foundProfile); err != nil {
	log.Fatalf("Failed to update profile: %v", err)
}

// Verify update
updatedProfile, err := profileRepo.FindByID(foundProfile.ID)
if err != nil {
	log.Fatalf("Failed to find updated profile: %v", err)
}

fmt.Printf("Updated profile bio: %s\n", updatedProfile.Bio)
```

### Adding a New Post to a User (One-to-Many)

```go
// Add a new post for the user
newPost := &Post{
	Title:   "Another Post",
	Content: "More content about ORM relationships...",
	UserID:  user.ID,
}

if err := postRepo.Save(newPost); err != nil {
	log.Fatalf("Failed to save new post: %v", err)
}

fmt.Printf("Added new post with ID: %d to user ID: %d\n", newPost.ID, newPost.UserID)
```

### Removing a Role from a User (Many-to-Many)

```go
// Remove user role from user
userRoleToDelete, err := userRoleRepo.Find().
	Where("user_id = ? AND role_id = ?", user.ID, userRole.ID).
	One()
if err != nil {
	log.Fatalf("Failed to find user role: %v", err)
}

if err := userRoleRepo.Delete(userRoleToDelete); err != nil {
	log.Fatalf("Failed to delete user role: %v", err)
}

fmt.Printf("Removed 'User' role from user ID: %d\n", user.ID)
```

## Conclusion

This example demonstrates how to model and work with different types of relationships in Goofer ORM:

1. **One-to-One**: User and Profile
2. **One-to-Many**: User and Posts, Post and Comments
3. **Many-to-Many**: Users and Roles through UserRoles join table

Goofer ORM makes it easy to define these relationships using struct tags and provides a fluent API for querying and manipulating related entities.