package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "reflect"
    "strings"
    "time"
)

// Entity interface for model metadata
type Entity interface {
    TableName() string
}

// ORM tag parser constants
const (
    TagName          = "orm"
    PrimaryKeyOption = "primaryKey"
    AutoIncrementOpt = "autoIncrement"
    UniqueOption     = "unique"
    IndexOption      = "index"
    RelationOption   = "relation"
    ForeignKeyOption = "foreignKey"
)

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
    Relation      *RelationMetadata
}

// RelationMetadata describes entity relationships
type RelationMetadata struct {
    Type       RelationType
    Entity     reflect.Type
    ForeignKey string
}

// RelationType defines relationship types
type RelationType int

const (
    HasOne RelationType = iota
    HasMany
    BelongsTo
    ManyToMany
)

// SchemaRegistry maintains entity metadata
type SchemaRegistry struct {
    entities map[reflect.Type]*EntityMetadata
}

// EntityMetadata contains complete entity schema
type EntityMetadata struct {
    TableName   string
    Fields      []FieldMetadata
    PrimaryKey  *FieldMetadata
    Relations   []RelationMetadata
    Indexes     []IndexMetadata
}

// IndexMetadata describes database indexes
type IndexMetadata struct {
    Name    string
    Columns []string
    Unique  bool
}

// Repository provides type-safe database operations
type Repository[T Entity] struct {
    db      *sql.DB
    dialect Dialect
    metadata *EntityMetadata
    ctx     context.Context
}

// QueryBuilder enables fluent query construction
type QueryBuilder[T Entity] struct {
    repo       *Repository[T]
    conditions []string
    args       []interface{}
    includes   []string
    order      string
    limit      int
    offset     int
}

// Dialect interface for database-specific implementations
type Dialect interface {
    Placeholder(int) string
    QuoteIdentifier(string) string
    DataType(field FieldMetadata) string
    CreateTableSQL(*EntityMetadata) string
}

// Initialize schema registry
var registry = &SchemaRegistry{
    entities: make(map[reflect.Type]*EntityMetadata),
}

// RegisterEntity analyzes and registers entity schema
func RegisterEntity(entity Entity) error {
    entityType := reflect.TypeOf(entity)
    if entityType.Kind() == reflect.Ptr {
        entityType = entityType.Elem()
    }

    meta := &EntityMetadata{
        TableName: entity.TableName(),
    }

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
        
        if fieldMeta.IsPrimaryKey {
            meta.PrimaryKey = fieldMeta
        }
        
        if fieldMeta.Relation != nil {
            meta.Relations = append(meta.Relations, *fieldMeta.Relation)
        }
    }

    registry.entities[entityType] = meta
    return nil
}

// parseFieldTag converts ORM tags to metadata
func parseFieldTag(field reflect.StructField, tag string) (*FieldMetadata, error) {
    options := parseTagOptions(tag)
    meta := &FieldMetadata{
        Name:   field.Name,
        DBName: snakeCase(field.Name),
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
        case strings.HasPrefix(opt, "type:"):
            meta.Type = strings.TrimPrefix(opt, "type:")
        case strings.HasPrefix(opt, RelationOption):
            rel, err := parseRelationOption(opt)
            if err != nil {
                return nil, err
            }
            meta.Relation = rel
        case strings.HasPrefix(opt, ForeignKeyOption):
            // Handle foreign key relationships
        }
    }

    return meta, nil
}

func parseTagOptions(tag string) []string {
    return strings.Split(tag, ";")
}

// NewRepository creates type-safe repository
func NewRepository[T Entity](db *sql.DB, dialect Dialect) *Repository[T] {
    var entity T
    entityType := reflect.TypeOf(entity)
    if entityType.Kind() == reflect.Ptr {
        entityType = entityType.Elem()
    }

    meta, exists := registry.entities[entityType]
    if !exists {
        panic(fmt.Sprintf("entity %s not registered", entityType.Name()))
    }

    return &Repository[T]{
        db:       db,
        dialect:  dialect,
        metadata: meta,
        ctx:      context.Background(),
    }
}

// Find initiates a query builder
func (r *Repository[T]) Find() *QueryBuilder[T] {
    return &QueryBuilder[T]{repo: r}
}

// Where adds condition to query
func (qb *QueryBuilder[T]) Where(cond string, args ...interface{}) *QueryBuilder[T] {
    qb.conditions = append(qb.conditions, cond)
    qb.args = append(qb.args, args...)
    return qb
}

// Include specifies relations to preload
func (qb *QueryBuilder[T]) Include(relations ...string) *QueryBuilder[T] {
    qb.includes = append(qb.includes, relations...)
    return qb
}

// Exec executes the built query
func (qb *QueryBuilder[T]) Exec() ([]T, error) {
    query := qb.buildSelectQuery()
    rows, err := qb.repo.db.QueryContext(qb.repo.ctx, query, qb.args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return qb.scanRows(rows)
}

// buildSelectQuery constructs the SQL query
func (qb *QueryBuilder[T]) buildSelectQuery() string {
    var selects []string
    for _, field := range qb.repo.metadata.Fields {
        selects = append(selects, qb.repo.dialect.QuoteIdentifier(field.DBName))
    }

    query := fmt.Sprintf("SELECT %s FROM %s",
        strings.Join(selects, ", "),
        qb.repo.dialect.QuoteIdentifier(qb.repo.metadata.TableName),
    )

    if len(qb.conditions) > 0 {
        query += " WHERE " + strings.Join(qb.conditions, " AND ")
    }

    return query
}

// Save handles insert/update operations
func (r *Repository[T]) Save(entity *T) error {
    meta := r.metadata
    if meta.PrimaryKey == nil {
        return errors.New("entity missing primary key")
    }

    val := reflect.ValueOf(entity).Elem()
    pkValue := val.FieldByName(meta.PrimaryKey.Name)

    if pkValue.IsZero() {
        return r.insert(entity)
    }
    return r.update(entity)
}

// Transaction executes a database transaction
func (r *Repository[T]) Transaction(fn func(*Repository[T]) error) error {
    tx, err := r.db.BeginTx(r.ctx, nil)
    if err != nil {
        return err
    }

    txRepo := &Repository[T]{
        db:       r.db,
        dialect:  r.dialect,
        metadata: r.metadata,
        ctx:      r.ctx,
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    err = fn(txRepo)
    return err
}

// Example entity definitions
type User struct {
    ID       uint      `orm:"primaryKey;autoIncrement"`
    Name     string    `orm:"type:varchar(255);notnull"`
    Email    string    `orm:"unique;type:varchar(255)"`
    Posts    []Post    `orm:"relation:HasMany;foreignKey:UserID"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string { return "users" }

type Post struct {
    ID      uint   `orm:"primaryKey"`
    Title   string `orm:"type:text;notnull"`
    UserID  uint   `orm:"index"`
    Author  *User  `orm:"relation:BelongsTo;foreignKey:UserID"`
}

func (Post) TableName() string { return "posts" }

// Example usage
func main() {
    // Initialize database
    db, _ := sql.Open("mysql", "user:pass@/dbname")
    
    // Register entities
    RegisterEntity(User{})
    RegisterEntity(Post{})

    // Create repository
    userRepo := NewRepository[User](db, MySQLDialect{})
    postRepo := NewRepository[Post](db, MySQLDialect{})

    // Create user with transaction
    userRepo.Transaction(func(r *Repository[User]) error {
        newUser := &User{
            Name:  "Tach",
            Email: "tach@dev.com",
        }
        
        if err := r.Save(newUser); err != nil {
            return err
        }

        // Create related post
        post := &Post{
            Title:  "Hello ORM",
            UserID: newUser.ID,
        }
        return postRepo.Save(post)
    })

    // Query with includes
    users, _ := userRepo.Find().
        Where("email LIKE ?", "%@dev.com").
        Include("Posts").
        Exec()

	fmt.Println(users)
}