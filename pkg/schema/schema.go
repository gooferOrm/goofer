package schema

import (
	"fmt"
	"reflect"
)

// Schema defines the database schema
// It's similar to Prisma's schema.prisma file
// but implemented in Go
const (
	// Field types
	TypeString   = "string"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeBoolean  = "boolean"
	TypeDateTime = "datetime"
	TypeEnum     = "enum"
	TypeJson     = "json"
	TypeBytes    = "bytes"

	// Field constraints
	ConstraintUnique = "unique"
	ConstraintIndex  = "index"
	ConstraintPK     = "primaryKey"
)

// Model represents a database model
// Similar to Prisma's model definition
// @model User {
//   id        Int      @id @default(autoincrement())
//   email     String   @unique
//   name      String
//   posts     Post[]   @relation("UserToPost")
// }
type Model struct {
	Name        string
	Fields      []Field
	Relations   []Relation
	Indexes     []Index
	UniqueKeys  []UniqueKey
}

// Field represents a model field
// Similar to Prisma's field definition
type Field struct {
	Name        string
	Type        string
	IsNullable  bool
	IsUnique    bool
	IsIndexed   bool
	IsPK        bool
	Default     interface{}
	EnumValues  []string
}

// Relation represents a model relation
// Similar to Prisma's relation definition
type Relation struct {
	Name        string
	Type        RelationType
	Model       string
	ForeignKey  string
	Through     string
}

// RelationType defines the type of relation
type RelationType string

const (
	OneToOne     RelationType = "OneToOne"
	OneToMany    RelationType = "OneToMany"
	ManyToOne    RelationType = "ManyToOne"
	ManyToMany   RelationType = "ManyToMany"
)

// Index represents a database index
type Index struct {
	Name    string
	Fields  []string
	Unique  bool
}

// UniqueKey represents a unique constraint
type UniqueKey struct {
	Name    string
	Fields  []string
}

// ParseModel parses a struct into a Model definition
// This is similar to Prisma's schema parsing
func ParseModel(entity interface{}) (*Model, error) {
	model := &Model{
		Name: reflect.TypeOf(entity).Name(),
	}

	// Parse struct fields
	return model, nil
}

// ValidateModel validates the model schema
func (m *Model) Validate() error {
	// Similar to Prisma's schema validation
	return nil
}
