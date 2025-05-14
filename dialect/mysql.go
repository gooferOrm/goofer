package dialect

import (
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/schema"
)

// MySQLDialect implements the Dialect interface for MySQL
type MySQLDialect struct {
	*BaseDialect
}

func NewMySQLDialect() *MySQLDialect {
	return &MySQLDialect{
		BaseDialect: &BaseDialect{},
	}
}

// Name returns the name of the dialect
func (d *MySQLDialect) Name() string {
	return "mysql"
}

// Placeholder returns the placeholder for a parameter at the given index
func (d *MySQLDialect) Placeholder(int) string {
	return "?"
}

// QuoteIdentifier quotes an identifier with backticks
func (d *MySQLDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", name)
}

// DataType maps a field metadata to a MySQL-specific type
func (d *MySQLDialect) DataType(field schema.FieldMetadata) string {
	if field.Type != "" {
		return field.Type
	}

	// Default type mapping
	switch {
	case strings.HasPrefix(field.Type, "varchar"):
		return field.Type
	case strings.HasPrefix(field.Type, "int"):
		return field.Type
	case field.IsAutoIncr:
		return "INT AUTO_INCREMENT"
	case strings.EqualFold(field.Type, "text"):
		return "TEXT"
	case strings.EqualFold(field.Type, "boolean"):
		return "TINYINT(1)"
	case strings.EqualFold(field.Type, "datetime"):
		return "DATETIME"
	case strings.EqualFold(field.Type, "timestamp"):
		return "TIMESTAMP"
	case strings.EqualFold(field.Type, "float"):
		return "FLOAT"
	case strings.EqualFold(field.Type, "double"):
		return "DOUBLE"
	case strings.EqualFold(field.Type, "decimal"):
		return "DECIMAL(10,2)"
	case strings.EqualFold(field.Type, "json"):
		return "JSON"
	case strings.EqualFold(field.Type, "blob"):
		return "BLOB"
	default:
		return "VARCHAR(255)"
	}
}

// CreateTableSQL generates SQL to create a table for the entity
func (d *MySQLDialect) CreateTableSQL(meta *schema.EntityMetadata) string {
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
		
		if field.IsAutoIncr && !strings.Contains(strings.ToUpper(d.DataType(field)), "AUTO_INCREMENT") {
			column += " AUTO_INCREMENT"
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
	builder.WriteString("\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")
	
	// Add indexes
	for _, field := range meta.Fields {
		if field.IsIndexed && !field.IsPrimaryKey && !field.IsUnique {
			indexName := fmt.Sprintf("idx_%s_%s", meta.TableName, field.DBName)
			indexSQL := fmt.Sprintf("\nCREATE INDEX %s ON %s (%s);",
				d.QuoteIdentifier(indexName),
				d.QuoteIdentifier(meta.TableName),
				d.QuoteIdentifier(field.DBName))
			builder.WriteString(indexSQL)
		}
	}
	
	return builder.String()
}