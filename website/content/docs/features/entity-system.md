# Entity System

The Entity System is the foundation of Goofer ORM. It allows you to define your database schema using Go structs with tags for metadata, providing a type-safe and intuitive way to work with your database.

## Defining Entities

An entity in Goofer ORM is a Go struct that represents a database table. Each field in the struct corresponds to a column in the table.

```go
// User entity
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"unique;type:varchar(255);notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
    return "users"
}
```

## Entity Interface

All entities must implement the `Entity` interface, which requires a `TableName()` method to specify the database table name:

```go
type Entity interface {
    TableName() string
}
```

This method allows you to customize the table name for each entity, which is especially useful when working with existing databases or when you want to use a naming convention different from the default.

## ORM Tags

Goofer ORM uses struct tags to define metadata for each field. The tag format is:

```
`orm:"option1;option2;option3"`
```

### Available Tag Options

| Option | Description | Example |
|--------|-------------|---------|
| `primaryKey` | Marks the field as the primary key | `orm:"primaryKey"` |
| `autoIncrement` | Enables auto-increment for the field | `orm:"autoIncrement"` |
| `type:TYPE` | Specifies the database column type | `orm:"type:varchar(255)"` |
| `notnull` | Makes the field non-nullable | `orm:"notnull"` |
| `unique` | Creates a unique constraint | `orm:"unique"` |
| `index` | Creates an index on the field | `orm:"index"` |
| `default:VALUE` | Sets a default value | `orm:"default:CURRENT_TIMESTAMP"` |
| `relation:TYPE` | Defines a relationship type | `orm:"relation:OneToMany"` |
| `foreignKey:FIELD` | Specifies the foreign key field | `orm:"foreignKey:UserID"` |
| `joinTable:TABLE` | Specifies the join table for many-to-many relationships | `orm:"joinTable:user_roles"` |
| `referenceKey:FIELD` | Specifies the reference key for many-to-many relationships | `orm:"referenceKey:RoleID"` |

## Entity Registration

Before using an entity with Goofer ORM, you need to register it with the schema registry:

```go
if err := schema.Registry.RegisterEntity(User{}); err != nil {
    log.Fatalf("Failed to register User entity: %v", err)
}
```

Registration analyzes the entity using reflection and stores its metadata in the registry, making it available for the ORM to use when generating SQL and performing database operations.

## Entity Metadata

After registration, you can access the entity's metadata:

```go
userMeta, exists := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
if !exists {
    log.Fatalf("User entity not registered")
}

fmt.Printf("Table name: %s\n", userMeta.TableName)
fmt.Printf("Number of fields: %d\n", len(userMeta.Fields))
```

The metadata includes information about the table name, fields, primary key, relationships, and more.

## Validation Tags

In addition to ORM tags, you can also use validation tags from the `go-playground/validator` package:

```go
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
    Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
    Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
    Age       int       `orm:"type:int" validate:"gte=0,lte=130"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
```

These validation tags are used by the validation system to ensure that your data meets your requirements before it's saved to the database.

## Best Practices

- Use meaningful names for your entities and fields
- Implement the `TableName()` method for all entities
- Use appropriate ORM tags to define your schema
- Add validation tags for data integrity
- Keep entities focused on a single responsibility
- Use relationships to model connections between entities

## Next Steps

- Learn about the [Schema Parser](./schema-parser) to understand how entity metadata is extracted
- Explore [Relation Mapping](./relation-mapping) to see how to define and work with entity relationships
- Check out the [Repository Pattern](./repository-pattern) to learn how to perform CRUD operations on entities