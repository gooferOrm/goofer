# Schema Parser

The Schema Parser is a core component of Goofer ORM that uses Go's reflection capabilities to analyze your entity structs at runtime. It extracts metadata from struct tags, maps Go types to database types, and builds a complete schema registry for your application.

## How It Works

The Schema Parser works by:

1. Analyzing entity structs using reflection
2. Parsing ORM tags on struct fields
3. Extracting metadata about fields, types, and relationships
4. Building a schema registry that the ORM can use for database operations

## Parsing Process

When you register an entity with Goofer ORM, the Schema Parser performs the following steps:

```go
// Register an entity
if err := schema.Registry.RegisterEntity(User{}); err != nil {
    log.Fatalf("Failed to register User entity: %v", err)
}
```

### 1. Type Analysis

First, the parser analyzes the entity's type using reflection:

```go
entityType := reflect.TypeOf(entity)
if entityType.Kind() == reflect.Ptr {
    entityType = entityType.Elem()
}
```

### 2. Table Name Resolution

Next, it calls the `TableName()` method to get the database table name:

```go
meta := &EntityMetadata{
    TableName: entity.TableName(),
}
```

### 3. Field Analysis

Then, it iterates through each field in the struct:

```go
for i := 0; i < entityType.NumField(); i++ {
    field := entityType.Field(i)
    tag := field.Tag.Get(TagName)
    if tag == "" || tag == "-" {
        continue
    }

    fieldMeta, err := parseFieldTag(field, tag)
    if err != nil {
        return err
    }

    meta.Fields = append(meta.Fields, *fieldMeta)
}
```

### 4. Tag Parsing

For each field, it parses the ORM tag to extract metadata:

```go
func parseFieldTag(field reflect.StructField, tag string) (*FieldMetadata, error) {
    options := parseTagOptions(tag)
    meta := &FieldMetadata{
        Name:       field.Name,
        DBName:     snakeCase(field.Name),
        IsNullable: true, // Default to nullable
    }

    for _, opt := range options {
        switch {
        case opt == PrimaryKeyOption:
            meta.IsPrimaryKey = true
        case opt == AutoIncrementOpt:
            meta.IsAutoIncr = true
        case opt == UniqueOption:
            meta.IsUnique = true
        case opt == IndexOption:
            meta.IsIndexed = true
        case opt == NotNullOption:
            meta.IsNullable = false
        case strings.HasPrefix(opt, TypeOption+":"):
            meta.Type = strings.TrimPrefix(opt, TypeOption+":")
        case strings.HasPrefix(opt, DefaultOption+":"):
            meta.Default = strings.TrimPrefix(opt, DefaultOption+":")
        case strings.HasPrefix(opt, RelationOption+":"):
            relType := strings.TrimPrefix(opt, RelationOption+":")
            meta.Relation = &RelationMetadata{
                Type: RelationType(relType),
            }
        case strings.HasPrefix(opt, ForeignKeyOption+":"):
            if meta.Relation != nil {
                meta.Relation.ForeignKey = strings.TrimPrefix(opt, ForeignKeyOption+":")
            }
        }
    }

    // Infer type from Go type if not specified
    if meta.Type == "" {
        meta.Type = inferSQLType(field.Type)
    }

    return meta, nil
}
```

### 5. Type Inference

If a field doesn't have an explicit type specified, the parser infers the SQL type from the Go type:

```go
func inferSQLType(t reflect.Type) string {
    switch t.Kind() {
    case reflect.String:
        return "VARCHAR(255)"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return "INTEGER"
    case reflect.Float32, reflect.Float64:
        return "FLOAT"
    case reflect.Bool:
        return "BOOLEAN"
    case reflect.Struct:
        if t.String() == "time.Time" {
            return "TIMESTAMP"
        }
    case reflect.Slice:
        if t.Elem().Kind() == reflect.Uint8 {
            return "BLOB"
        }
    }
    return "TEXT"
}
```

### 6. Relationship Analysis

For fields that represent relationships, the parser extracts relationship metadata:

```go
if meta.Relation != nil {
    meta.Relations = append(meta.Relations, *fieldMeta.Relation)
}
```

### 7. Schema Registry

Finally, the parser stores the entity metadata in the global schema registry:

```go
r.entities[entityType] = meta
```

## Schema Registry

The Schema Registry is a global repository of entity metadata that the ORM uses for database operations:

```go
// Global registry instance
var Registry = NewSchemaRegistry()

// SchemaRegistry maintains entity metadata
type SchemaRegistry struct {
    entities map[reflect.Type]*EntityMetadata
}
```

You can access the registry to get metadata for an entity:

```go
userMeta, exists := schema.Registry.GetEntityMetadata(schema.GetEntityType(User{}))
if !exists {
    log.Fatalf("User entity not registered")
}
```

## Entity Metadata

The entity metadata contains all the information the ORM needs to work with an entity:

```go
// EntityMetadata contains complete entity schema
type EntityMetadata struct {
    TableName   string
    Fields      []FieldMetadata
    PrimaryKey  *FieldMetadata
    Relations   []RelationMetadata
    Indexes     []IndexMetadata
}
```

## Field Metadata

Each field in an entity has its own metadata:

```go
// FieldMetadata contains parsed ORM tag information
type FieldMetadata struct {
    Name          string
    DBName        string
    Type          string
    IsPrimaryKey  bool
    IsAutoIncr    bool
    IsUnique      bool
    IsIndexed     bool
    IsNullable    bool
    Default       interface{}
    Relation      *RelationMetadata
}
```

## Relation Metadata

Relationship fields have additional metadata:

```go
// RelationMetadata describes entity relationships
type RelationMetadata struct {
    Type       RelationType
    Entity     reflect.Type
    ForeignKey string
}
```

## Best Practices

- Register all entities at application startup
- Use explicit types in ORM tags for clarity
- Keep entity definitions clean and focused
- Use meaningful names for fields and relationships
- Validate the schema registry after registration

## Next Steps

- Learn about [Entity System](./entity-system) to understand how to define entities
- Explore [Relation Mapping](./relation-mapping) to see how to define and work with entity relationships
- Check out the [Migration Engine](./migration-engine) to learn how to generate database schemas from entity metadata