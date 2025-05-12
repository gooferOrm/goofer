# Relation Mapping

Relation Mapping in Goofer ORM makes it easy to work with related entities. It supports various relationship types and provides a clean API for defining and querying relationships.

## Supported Relationship Types

Goofer ORM supports four main types of relationships:

1. **One-to-One**: A relationship where each record in one table is associated with exactly one record in another table.
2. **One-to-Many**: A relationship where each record in one table can be associated with multiple records in another table.
3. **Many-to-One**: The inverse of a one-to-many relationship, where multiple records in one table can be associated with a single record in another table.
4. **Many-to-Many**: A relationship where multiple records in one table can be associated with multiple records in another table, typically using a join table.

## Defining Relationships

Relationships are defined using the `orm` tag on struct fields:

### One-to-One Relationship

```go
// User entity with a one-to-one relationship to Profile
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    
    // One-to-One relationship: User has one Profile
    Profile   *Profile  `orm:"relation:OneToOne;foreignKey:UserID"`
}

// Profile entity with a one-to-one relationship to User
type Profile struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    UserID    uint      `orm:"unique;notnull"` // Foreign key to User
    Bio       string    `orm:"type:text"`
    Avatar    string    `orm:"type:varchar(255)"`
    
    // Reference back to User
    User      *User     `orm:"relation:OneToOne;foreignKey:UserID"`
}
```

### One-to-Many Relationship

```go
// User entity with a one-to-many relationship to Post
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    
    // One-to-Many relationship: User has many Posts
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

// Post entity with a many-to-one relationship to User
type Post struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Title     string    `orm:"type:varchar(255);notnull"`
    Content   string    `orm:"type:text"`
    UserID    uint      `orm:"index;notnull"` // Foreign key to User
    
    // Many-to-One relationship: Post belongs to User
    User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}
```

### Many-to-Many Relationship

```go
// User entity with a many-to-many relationship to Role
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    
    // Many-to-Many relationship: User has many Roles through UserRoles
    Roles     []Role    `orm:"relation:ManyToMany;joinTable:user_roles;foreignKey:UserID;referenceKey:RoleID"`
}

// Role entity with a many-to-many relationship to User
type Role struct {
    ID          uint      `orm:"primaryKey;autoIncrement"`
    Name        string    `orm:"type:varchar(50);unique;notnull"`
    Description string    `orm:"type:text"`
    
    // Many-to-Many relationship: Role has many Users through UserRoles
    Users       []User    `orm:"relation:ManyToMany;joinTable:user_roles;foreignKey:RoleID;referenceKey:UserID"`
}

// UserRole entity for the join table
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

## Relationship Tag Options

When defining relationships, you can use the following tag options:

| Option | Description | Example |
|--------|-------------|---------|
| `relation:TYPE` | Defines the relationship type | `orm:"relation:OneToMany"` |
| `foreignKey:FIELD` | Specifies the foreign key field | `orm:"foreignKey:UserID"` |
| `joinTable:TABLE` | Specifies the join table for many-to-many relationships | `orm:"joinTable:user_roles"` |
| `referenceKey:FIELD` | Specifies the reference key for many-to-many relationships | `orm:"referenceKey:RoleID"` |

## Working with Relationships

### Creating Related Entities

To create related entities, you first create the parent entity, then create the child entities with the appropriate foreign key:

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
```

### Querying Related Entities

To query related entities, you can use the repository's `Find()` method with appropriate conditions:

#### Finding One-to-One Relationships

```go
// Find user with ID
foundUser, err := userRepo.FindByID(user.ID)
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Find profile by user ID (One-to-One)
foundProfile, err := profileRepo.Find().
    Where("user_id = ?", user.ID).
    One()
if err != nil {
    log.Fatalf("Failed to find profile: %v", err)
}

fmt.Printf("User's profile: Bio: %s, Avatar: %s\n", foundProfile.Bio, foundProfile.Avatar)
```

#### Finding One-to-Many Relationships

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
}
```

#### Finding Many-to-Many Relationships

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

## Eager Loading vs. Lazy Loading

Goofer ORM currently uses a lazy loading approach for relationships, meaning that related entities are not automatically loaded when you query an entity. You need to explicitly query for related entities as shown in the examples above.

Future versions of Goofer ORM may support eager loading, which would allow you to automatically load related entities in a single query.

## Best Practices

- Define relationships on both sides (parent and child) for clarity
- Use appropriate indexes on foreign key fields for better performance
- Use meaningful names for relationship fields
- Consider using transactions when creating or updating related entities
- Be mindful of N+1 query problems when working with relationships

## Next Steps

- Learn about the [Repository Pattern](./repository-pattern) to see how to perform CRUD operations on entities
- Explore [Transactions](./transactions) to understand how to ensure data integrity when working with related entities
- Check out the [Examples](../examples/relationships) section for more detailed examples of working with relationships