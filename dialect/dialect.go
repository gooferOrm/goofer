package dialect

import (
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/schema"
)

// Dialect interface for database-specific implementations
type Dialect interface {
	// Placeholder returns the placeholder for a parameter at the given index
	Placeholder(int) string
	
	// QuoteIdentifier quotes an identifier (table name, column name)
	QuoteIdentifier(string) string
	
	// DataType maps a field metadata to a database-specific type
	DataType(field schema.FieldMetadata) string
	
	// CreateTableSQL generates SQL to create a table for the entity
	CreateTableSQL(*schema.EntityMetadata) string
	
	// Name returns the name of the dialect
	Name() string
}

// BaseDialect provides common functionality for dialects
type BaseDialect struct{
	Dialect
}

// QuoteIdentifier quotes an identifier with double quotes
func (d *BaseDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

// DataType provides a default implementation that can be overridden by specific dialects
func (d *BaseDialect) DataType(field schema.FieldMetadata) string {
	switch field.Type {
	case "varchar":
		return "TEXT"
	case "int", "integer":
		return "INTEGER"
	case "float", "double":
		return "REAL"
	case "boolean":
		return "INTEGER"
	case "datetime", "timestamp":
		return "TEXT"
	case "json":
		return "TEXT"
	case "blob":
		return "BLOB"
	default:
		return "TEXT"
	}
}

// CreateTableSQL generates SQL to create a table for the entity
func (d *BaseDialect) CreateTableSQL(meta *schema.EntityMetadata) string {
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