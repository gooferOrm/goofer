package dialect

import (
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/pkg/schema"
)

// SQLiteDialect implements the Dialect interface for SQLite
type SQLiteDialect struct {
	*BaseDialect
}

// NewSQLiteDialect creates a new SQLite dialect instance
func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{
		BaseDialect: &BaseDialect{},
	}
}

// Name returns the name of the dialect
func (d *SQLiteDialect) Name() string {
	return "sqlite"
}

// Placeholder returns the placeholder for a parameter at the given index
func (d *SQLiteDialect) Placeholder(int) string {
	return "?"
}

// QuoteIdentifier quotes an identifier with double quotes
func (d *SQLiteDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

// DataType maps a field metadata to a SQLite-specific type
func (d *SQLiteDialect) DataType(field schema.FieldMetadata) string {
	// SQLite has a simpler type system
	if field.IsAutoIncr {
		return "INTEGER"
	}

	if field.Type != "" {
		// Check for type prefixes and convert them to SQLite types
		if strings.HasPrefix(field.Type, "varchar") {
			return "TEXT"
		} else if strings.HasPrefix(field.Type, "int") {
			return "INTEGER"
		} else if strings.EqualFold(field.Type, "text") {
			return "TEXT"
		} else if strings.EqualFold(field.Type, "boolean") {
			return "INTEGER"
		} else if strings.EqualFold(field.Type, "datetime") {
			return "TEXT"
		} else if strings.EqualFold(field.Type, "timestamp") {
			return "TEXT"
		} else if strings.EqualFold(field.Type, "float") {
			return "REAL"
		} else if strings.EqualFold(field.Type, "double") {
			return "REAL"
		} else if strings.EqualFold(field.Type, "decimal") {
			return "REAL"
		} else if strings.EqualFold(field.Type, "json") {
			return "TEXT"
		} else if strings.EqualFold(field.Type, "blob") {
			return "BLOB"
		}

		// If no conversion is needed, return the type as is
		return field.Type
	}

	// Default to TEXT for unknown types
	return "TEXT"
}

// CreateTableSQL generates SQL to create a table for the entity
func (d *SQLiteDialect) CreateTableSQL(meta *schema.EntityMetadata) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", d.QuoteIdentifier(meta.TableName)))

	var columns []string
	for _, field := range meta.Fields {
		// Skip relation fields
		if field.Relation != nil {
			continue
		}

		column := fmt.Sprintf("  %s %s", d.QuoteIdentifier(field.DBName), d.DataType(field))

		if field.IsPrimaryKey {
			column += " PRIMARY KEY"
		}

		if field.IsAutoIncr {
			column += " AUTOINCREMENT"
		}

		if !field.IsNullable {
			column += " NOT NULL"
		}

		if field.IsUnique {
			column += " UNIQUE"
		}

		if field.Default != nil {
			column += fmt.Sprintf(" DEFAULT %v", field.Default)
		}

		columns = append(columns, column)
	}

	builder.WriteString(strings.Join(columns, ",\n"))
	builder.WriteString("\n);")

	// Add indexes
	for _, field := range meta.Fields {
		if field.IsIndexed && !field.IsPrimaryKey && !field.IsUnique {
			indexName := fmt.Sprintf("idx_%s_%s", meta.TableName, field.DBName)
			indexSQL := fmt.Sprintf("\nCREATE INDEX IF NOT EXISTS %s ON %s (%s);",
				d.QuoteIdentifier(indexName),
				d.QuoteIdentifier(meta.TableName),
				d.QuoteIdentifier(field.DBName))
			builder.WriteString(indexSQL)
		}
	}

	return builder.String()
}
