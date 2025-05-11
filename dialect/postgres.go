package dialect

import (
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/schema"
)

// PostgresDialect implements the Dialect interface for PostgreSQL
type PostgresDialect struct {
	BaseDialect
}

// Name returns the name of the dialect
func (d *PostgresDialect) Name() string {
	return "postgres"
}

// Placeholder returns the placeholder for a parameter at the given index
func (d *PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}

// QuoteIdentifier quotes an identifier with double quotes
func (d *PostgresDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

// DataType maps a field metadata to a PostgreSQL-specific type
func (d *PostgresDialect) DataType(field schema.FieldMetadata) string {
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
		return "SERIAL"
	case strings.EqualFold(field.Type, "text"):
		return "TEXT"
	case strings.EqualFold(field.Type, "boolean"):
		return "BOOLEAN"
	case strings.EqualFold(field.Type, "datetime"):
		return "TIMESTAMP"
	case strings.EqualFold(field.Type, "timestamp"):
		return "TIMESTAMP"
	case strings.EqualFold(field.Type, "float"):
		return "REAL"
	case strings.EqualFold(field.Type, "double"):
		return "DOUBLE PRECISION"
	case strings.EqualFold(field.Type, "decimal"):
		return "NUMERIC(10,2)"
	case strings.EqualFold(field.Type, "json"):
		return "JSONB"
	case strings.EqualFold(field.Type, "blob"):
		return "BYTEA"
	default:
		return "VARCHAR(255)"
	}
}

// CreateTableSQL generates SQL to create a table for the entity
func (d *PostgresDialect) CreateTableSQL(meta *schema.EntityMetadata) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", d.QuoteIdentifier(meta.TableName)))
	
	var columns []string
	for _, field := range meta.Fields {
		// Skip relation fields
		if field.Relation != nil {
			continue
		}
		
		var column string
		
		// Handle auto-increment primary key specially for PostgreSQL
		if field.IsPrimaryKey && field.IsAutoIncr {
			if strings.EqualFold(field.Type, "int") || field.Type == "" {
				column = fmt.Sprintf("  %s SERIAL PRIMARY KEY", d.QuoteIdentifier(field.DBName))
			} else if strings.EqualFold(field.Type, "bigint") {
				column = fmt.Sprintf("  %s BIGSERIAL PRIMARY KEY", d.QuoteIdentifier(field.DBName))
			} else {
				column = fmt.Sprintf("  %s %s PRIMARY KEY", d.QuoteIdentifier(field.DBName), d.DataType(field))
			}
		} else {
			column = fmt.Sprintf("  %s %s", d.QuoteIdentifier(field.DBName), d.DataType(field))
			
			if field.IsPrimaryKey {
				column += " PRIMARY KEY"
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