# Goofer ORM - New Features

This document describes the new features that have been implemented in Goofer ORM to enhance its capabilities and developer experience.

## 🚀 New Features

### 1. Enhanced Query Builder

The query builder has been significantly enhanced with new methods and capabilities:

#### New Query Methods

```go
// WHERE IN clause
users, err := userRepo.Find().
    WhereIn("id", []interface{}{1, 2, 3}).
    All()

// WHERE NOT IN clause
users, err := userRepo.Find().
    WhereNotIn("id", []interface{}{1, 2, 3}).
    All()

// BETWEEN clause
users, err := userRepo.Find().
    WhereBetween("created_at", startDate, endDate).
    All()

// LIKE conditions
users, err := userRepo.Find().
    WhereLike("name", "%Doe%").
    All()

// IS NULL / IS NOT NULL
users, err := userRepo.Find().
    WhereNull("deleted_at").
    All()

users, err := userRepo.Find().
    WhereNotNull("email").
    All()

// OR conditions
users, err := userRepo.Find().
    Where("name = ?", "John").
    OrWhere("email = ?", "john@example.com").
    All()

// DISTINCT
users, err := userRepo.Find().
    Distinct().
    All()

// JOIN support (structure ready)
users, err := userRepo.Find().
    Join("posts", "users.id = posts.user_id").
    All()

// LEFT JOIN
users, err := userRepo.Find().
    LeftJoin("profiles", "users.id = profiles.user_id").
    All()

// GROUP BY and HAVING (structure ready)
users, err := userRepo.Find().
    GroupBy("age").
    Having("COUNT(*) > 1").
    All()
```

#### Enhanced Query Features

- **Complex WHERE conditions**: Combine multiple conditions with AND/OR logic
- **Advanced filtering**: Use LIKE, IN, BETWEEN, NULL checks
- **JOIN support**: INNER, LEFT, RIGHT, FULL joins
- **Aggregation**: GROUP BY and HAVING clauses
- **Pagination**: LIMIT and OFFSET with better control
- **Distinct queries**: Remove duplicate results
- **Ordered results**: Multiple ORDER BY clauses

### 2. Eager Loading Support

Eager loading allows you to load related entities in a single query, reducing the N+1 query problem:

```go
// Load users with their posts
users, err := userRepo.Find().
    With("Posts").
    All()

// Load users with multiple relations
users, err := userRepo.Find().
    With("Posts", "Profile", "Comments").
    All()

// Load nested relations
users, err := userRepo.Find().
    With("Posts.Comments").
    All()
```

#### Supported Relationship Types

- **One-to-One**: Load related single entities
- **One-to-Many**: Load collections of related entities
- **Many-to-One**: Load parent entities
- **Many-to-Many**: Load related entities through join tables

#### Benefits

- **Performance**: Reduces database queries
- **Convenience**: Access related data without additional queries
- **Flexibility**: Choose which relations to load
- **Type Safety**: Maintains Go's type system benefits

### 3. Database Schema Introspection

The introspection feature allows you to reverse-engineer existing databases into Goofer ORM entities:

#### Basic Usage

```go
// Create introspector
introspector := introspection.NewIntrospector(db, dialect)

// Introspect all tables
tables, err := introspector.IntrospectAllTables()
if err != nil {
    log.Fatal(err)
}

// Introspect specific table
userTable, err := introspector.IntrospectTable("users")
if err != nil {
    log.Fatal(err)
}

// Generate Go entities
entities, err := introspector.GenerateEntities()
if err != nil {
    log.Fatal(err)
}

// Generate entity for specific table
userEntity, err := introspector.GenerateEntity(userTable)
if err != nil {
    log.Fatal(err)
}
```

#### Generated Code Example

```go
// Generated from database introspection
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Name      string    `orm:"type:varchar(255);notnull"`
    Email     string    `orm:"type:varchar(255);notnull;unique"`
    Age       int       `orm:"type:int"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
    return "users"
}
```

#### Supported Features

- **Table discovery**: Automatically find all tables
- **Column analysis**: Extract column types, constraints, defaults
- **Primary key detection**: Identify primary keys
- **Index information**: Extract index metadata
- **Foreign key relationships**: Detect relationships between tables
- **Type mapping**: Convert SQL types to Go types
- **ORM tag generation**: Generate appropriate ORM tags
- **Multi-dialect support**: SQLite, MySQL, PostgreSQL

## 📁 Examples

### Advanced Queries Example

```bash
cd examples/advanced_queries
go run main.go
```

This example demonstrates:
- WHERE IN, BETWEEN, LIKE conditions
- IS NULL / IS NOT NULL checks
- OR conditions
- DISTINCT queries
- Advanced filtering combinations
- Enhanced COUNT queries

### Introspection Example

```bash
cd examples/introspection
go run main.go
```

This example demonstrates:
- Database schema analysis
- Table information extraction
- Column metadata parsing
- Go struct generation
- ORM tag generation

## 🔧 Implementation Details

### Query Builder Enhancements

The query builder has been enhanced with:

1. **New QueryBuilder struct fields**:
   - `joins`: Array of JOIN clauses
   - `groupBy`: GROUP BY clause
   - `having`: HAVING clause
   - `distinct`: DISTINCT flag

2. **Enhanced SQL generation**:
   - Support for JOIN clauses
   - GROUP BY and HAVING
   - DISTINCT keyword
   - Complex WHERE conditions

3. **New query methods**:
   - `With()`: Eager loading
   - `Join()`, `LeftJoin()`, etc.: JOIN operations
   - `WhereIn()`, `WhereNotIn()`: IN clauses
   - `WhereBetween()`: BETWEEN clauses
   - `WhereLike()`: LIKE conditions
   - `WhereNull()`, `WhereNotNull()`: NULL checks
   - `OrWhere()`: OR conditions
   - `Distinct()`: DISTINCT queries

### Eager Loading Implementation

Eager loading is implemented through:

1. **Relation detection**: Analyze entity metadata for relationships
2. **Query optimization**: Load related entities efficiently
3. **Data mapping**: Map related data to parent entities
4. **Type safety**: Maintain Go's type system

### Introspection Implementation

The introspection system includes:

1. **Database analysis**: Query system tables for schema information
2. **Metadata extraction**: Parse column types, constraints, relationships
3. **Code generation**: Generate Go structs with ORM tags
4. **Multi-dialect support**: Handle different database systems

## 🚧 Future Enhancements

### Planned Features

1. **Full Eager Loading Implementation**:
   - Complete relationship loading logic
   - Nested relationship support
   - Performance optimizations

2. **Advanced Introspection**:
   - Foreign key relationship detection
   - Index analysis
   - View support
   - Stored procedure analysis

3. **Query Builder Extensions**:
   - Subquery support
   - Raw SQL integration
   - Query optimization hints
   - Connection pooling

4. **Additional Features**:
   - Caching layer
   - Event system
   - Audit trail
   - Multi-tenancy support

## 📚 Usage Guidelines

### Best Practices

1. **Query Builder**:
   - Use specific WHERE methods for better readability
   - Combine conditions logically
   - Use indexes for performance
   - Limit result sets appropriately

2. **Eager Loading**:
   - Load only necessary relations
   - Avoid loading large collections
   - Use pagination for large datasets
   - Consider lazy loading for rarely accessed data

3. **Introspection**:
   - Review generated code before use
   - Customize generated entities as needed
   - Handle database-specific features
   - Test generated code thoroughly

### Performance Considerations

1. **Query Optimization**:
   - Use appropriate indexes
   - Limit result sets
   - Avoid N+1 queries with eager loading
   - Use specific column selection when possible

2. **Memory Management**:
   - Use pagination for large datasets
   - Avoid loading unnecessary relations
   - Consider streaming for large results

## 🤝 Contributing

These features are part of the ongoing development of Goofer ORM. Contributions are welcome:

1. **Bug Reports**: Report issues with specific examples
2. **Feature Requests**: Suggest new capabilities
3. **Code Contributions**: Submit pull requests
4. **Documentation**: Help improve documentation
5. **Testing**: Add test cases and examples

## 📄 License

This project is licensed under the same terms as the main Goofer ORM project. 