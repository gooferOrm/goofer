package introspection

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/dialect"
)

// Introspector provides database schema introspection capabilities
type Introspector struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewIntrospector creates a new introspector for the given database and dialect
func NewIntrospector(db *sql.DB, d dialect.Dialect) *Introspector {
	return &Introspector{
		db:      db,
		dialect: d,
	}
}

// TableInfo represents information about a database table
type TableInfo struct {
	Name        string
	Columns     []ColumnInfo
	PrimaryKey  string
	Indexes     []IndexInfo
	ForeignKeys []ForeignKeyInfo
}

// ColumnInfo represents information about a database column
type ColumnInfo struct {
	Name         string
	Type         string
	IsNullable   bool
	IsPrimaryKey bool
	IsUnique     bool
	DefaultValue *string
	Comment      string
}

// IndexInfo represents information about a database index
type IndexInfo struct {
	Name     string
	Columns  []string
	IsUnique bool
}

// ForeignKeyInfo represents information about a foreign key constraint
type ForeignKeyInfo struct {
	Name             string
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

// IntrospectTable introspects a single table and returns its information
func (i *Introspector) IntrospectTable(tableName string) (*TableInfo, error) {
	info := &TableInfo{
		Name: tableName,
	}

	// Get column information
	columns, err := i.getColumns(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
	}
	info.Columns = columns

	// Get primary key information
	pk, err := i.getPrimaryKey(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary key for table %s: %w", tableName, err)
	}
	info.PrimaryKey = pk

	// Get index information
	indexes, err := i.getIndexes(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes for table %s: %w", tableName, err)
	}
	info.Indexes = indexes

	// Get foreign key information
	foreignKeys, err := i.getForeignKeys(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign keys for table %s: %w", tableName, err)
	}
	info.ForeignKeys = foreignKeys

	return info, nil
}

// IntrospectAllTables introspects all tables in the database
func (i *Introspector) IntrospectAllTables() ([]*TableInfo, error) {
	tables, err := i.getTableNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get table names: %w", err)
	}

	var tableInfos []*TableInfo
	for _, table := range tables {
		info, err := i.IntrospectTable(table)
		if err != nil {
			return nil, err
		}
		tableInfos = append(tableInfos, info)
	}

	return tableInfos, nil
}

// GenerateEntity generates a Go struct from table information
func (i *Introspector) GenerateEntity(tableInfo *TableInfo) (string, error) {
	var builder strings.Builder

	// Generate struct name (convert table name to PascalCase)
	structName := toPascalCase(tableInfo.Name)

	builder.WriteString(fmt.Sprintf("// %s represents the %s table\n", structName, tableInfo.Name))
	builder.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	// Generate fields
	for _, column := range tableInfo.Columns {
		fieldName := toPascalCase(column.Name)
		goType := i.mapSQLTypeToGoType(column.Type)

		// Build ORM tags
		tags := i.buildORMTags(column, tableInfo)

		builder.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, goType, tags))
	}

	builder.WriteString("}\n\n")

	// Generate TableName method
	builder.WriteString(fmt.Sprintf("// TableName returns the table name for the %s entity\n", structName))
	builder.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	builder.WriteString(fmt.Sprintf("\treturn \"%s\"\n", tableInfo.Name))
	builder.WriteString("}\n")

	return builder.String(), nil
}

// GenerateEntities generates Go structs for all tables
func (i *Introspector) GenerateEntities() (string, error) {
	tables, err := i.IntrospectAllTables()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("package models\n\n")
	builder.WriteString("import \"time\"\n\n")

	for _, table := range tables {
		entity, err := i.GenerateEntity(table)
		if err != nil {
			return "", err
		}
		builder.WriteString(entity)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

// getTableNames retrieves all table names from the database
func (i *Introspector) getTableNames() ([]string, error) {
	var query string
	switch i.dialect.Name() {
	case "sqlite":
		query = "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
	case "mysql":
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()"
	case "postgres":
		query = "SELECT tablename FROM pg_tables WHERE schemaname = 'public'"
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", i.dialect.Name())
	}

	rows, err := i.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// getColumns retrieves column information for a table
func (i *Introspector) getColumns(tableName string) ([]ColumnInfo, error) {
	var query string
	switch i.dialect.Name() {
	case "sqlite":
		query = "PRAGMA table_info(" + i.dialect.QuoteIdentifier(tableName) + ")"
	case "mysql":
		query = `
			SELECT 
				column_name, 
				data_type, 
				is_nullable = 'YES' as is_nullable,
				column_key = 'PRI' as is_primary_key,
				column_default,
				column_comment
			FROM information_schema.columns 
			WHERE table_schema = DATABASE() AND table_name = ?
		`
	case "postgres":
		query = `
			SELECT 
				column_name,
				data_type,
				is_nullable = 'YES' as is_nullable,
				column_default,
				col_description((table_schema||'.'||table_name)::regclass, ordinal_position) as comment
			FROM information_schema.columns 
			WHERE table_schema = 'public' AND table_name = ?
		`
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", i.dialect.Name())
	}

	var rows *sql.Rows
	var err error

	if i.dialect.Name() == "sqlite" {
		rows, err = i.db.Query(query)
	} else {
		rows, err = i.db.Query(query, tableName)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var isNullable string
		var isPrimaryKey string
		var defaultValue sql.NullString
		var comment sql.NullString

		switch i.dialect.Name() {
		case "sqlite":
			var cid int
			var notNull int
			var pk int
			err = rows.Scan(&cid, &col.Name, &col.Type, &notNull, &pk, &defaultValue)
			col.IsNullable = notNull == 0
			col.IsPrimaryKey = pk == 1
		case "mysql":
			err = rows.Scan(&col.Name, &col.Type, &isNullable, &isPrimaryKey, &defaultValue, &comment)
			col.IsNullable = isNullable == "YES"
			col.IsPrimaryKey = isPrimaryKey == "PRI"
		case "postgres":
			err = rows.Scan(&col.Name, &col.Type, &isNullable, &defaultValue, &comment)
			col.IsNullable = isNullable == "YES"
		}

		if err != nil {
			return nil, err
		}

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
		}
		if comment.Valid {
			col.Comment = comment.String
		}

		columns = append(columns, col)
	}

	return columns, nil
}

// getPrimaryKey retrieves primary key information for a table
func (i *Introspector) getPrimaryKey(tableName string) (string, error) {
	// For now, we'll get this from the columns query
	// In a more complete implementation, you'd query the database's system tables
	columns, err := i.getColumns(tableName)
	if err != nil {
		return "", err
	}

	for _, col := range columns {
		if col.IsPrimaryKey {
			return col.Name, nil
		}
	}

	return "", nil
}

// getIndexes retrieves index information for a table
func (i *Introspector) getIndexes(tableName string) ([]IndexInfo, error) {
	// This is a simplified implementation
	// In a complete implementation, you'd query the database's system tables
	return []IndexInfo{}, nil
}

// getForeignKeys retrieves foreign key information for a table
func (i *Introspector) getForeignKeys(tableName string) ([]ForeignKeyInfo, error) {
	// This is a simplified implementation
	// In a complete implementation, you'd query the database's system tables
	return []ForeignKeyInfo{}, nil
}

// mapSQLTypeToGoType maps SQL types to Go types
func (i *Introspector) mapSQLTypeToGoType(sqlType string) string {
	sqlType = strings.ToLower(sqlType)

	switch {
	case strings.Contains(sqlType, "int"):
		if strings.Contains(sqlType, "bigint") {
			return "int64"
		}
		return "int"
	case strings.Contains(sqlType, "varchar"), strings.Contains(sqlType, "text"), strings.Contains(sqlType, "char"):
		return "string"
	case strings.Contains(sqlType, "decimal"), strings.Contains(sqlType, "numeric"), strings.Contains(sqlType, "float"), strings.Contains(sqlType, "double"):
		return "float64"
	case strings.Contains(sqlType, "bool"):
		return "bool"
	case strings.Contains(sqlType, "date"), strings.Contains(sqlType, "time"), strings.Contains(sqlType, "timestamp"):
		return "time.Time"
	case strings.Contains(sqlType, "json"):
		return "string" // Could be map[string]interface{} or custom type
	case strings.Contains(sqlType, "blob"), strings.Contains(sqlType, "binary"):
		return "[]byte"
	default:
		return "string"
	}
}

// buildORMTags builds ORM tags for a column
func (i *Introspector) buildORMTags(column ColumnInfo, tableInfo *TableInfo) string {
	var tags []string

	// Add type
	tags = append(tags, fmt.Sprintf("type:%s", column.Type))

	// Add primary key
	if column.IsPrimaryKey {
		tags = append(tags, "primaryKey")
		if strings.Contains(strings.ToLower(column.Type), "int") {
			tags = append(tags, "autoIncrement")
		}
	}

	// Add nullable
	if !column.IsNullable {
		tags = append(tags, "notnull")
	}

	// Add unique
	if column.IsUnique {
		tags = append(tags, "unique")
	}

	// Add default value
	if column.DefaultValue != nil {
		tags = append(tags, fmt.Sprintf("default:%s", *column.DefaultValue))
	}

	return fmt.Sprintf(`orm:"%s"`, strings.Join(tags, ";"))
}

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}
