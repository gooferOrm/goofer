package schema

import (
	"errors"
	"reflect"
	"strings"
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
	NotNullOption    = "notnull"
	RelationOption   = "relation"
	ForeignKeyOption = "foreignKey"
	DefaultOption    = "default"
	TypeOption       = "type"
)

// Field types
const (
	TypeString   = "string"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeBoolean  = "boolean"
	TypeDateTime = "datetime"
	TypeEnum     = "enum"
	TypeJson     = "json"
	TypeBytes    = "bytes"
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
	Default       interface{}
	Relation      *RelationMetadata
}

// RelationMetadata describes entity relationships
type RelationMetadata struct {
	Type       RelationType
	Entity     reflect.Type
	ForeignKey string
}

// RelationType defines relationship types
type RelationType string

const (
	OneToOne     RelationType = "OneToOne"
	OneToMany    RelationType = "OneToMany"
	ManyToOne    RelationType = "ManyToOne"
	ManyToMany   RelationType = "ManyToMany"
)

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

// SchemaRegistry maintains entity metadata
type SchemaRegistry struct {
	entities map[reflect.Type]*EntityMetadata
}

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		entities: make(map[reflect.Type]*EntityMetadata),
	}
}

// Global registry instance
var Registry = NewSchemaRegistry()

// RegisterEntity analyzes and registers entity schema
func (r *SchemaRegistry) RegisterEntity(entity Entity) error {
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

	r.entities[entityType] = meta
	return nil
}

// GetEntityMetadata retrieves metadata for an entity type
func (r *SchemaRegistry) GetEntityMetadata(entityType reflect.Type) (*EntityMetadata, bool) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	meta, exists := r.entities[entityType]
	return meta, exists
}

// parseFieldTag converts ORM tags to metadata
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

// parseTagOptions splits tag string into options
func parseTagOptions(tag string) []string {
	return strings.Split(tag, ";")
}

// inferSQLType maps Go types to SQL types
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

// snakeCase converts CamelCase to snake_case
func snakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ValidateEntityMetadata checks if entity metadata is valid
func ValidateEntityMetadata(meta *EntityMetadata) error {
	if meta.TableName == "" {
		return errors.New("entity must have a table name")
	}

	if len(meta.Fields) == 0 {
		return errors.New("entity must have at least one field")
	}

	if meta.PrimaryKey == nil {
		return errors.New("entity must have a primary key")
	}

	return nil
}

// GetEntityType returns the reflect.Type of an entity
func GetEntityType(entity Entity) reflect.Type {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// GetAllEntities returns all registered entities
func (r *SchemaRegistry) GetAllEntities() []*EntityMetadata {
	var entities []*EntityMetadata
	for _, meta := range r.entities {
		entities = append(entities, meta)
	}
	return entities
}
