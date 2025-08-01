package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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

// DBExecutor is an interface that both *sql.DB and *sql.Tx implement
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// AnyEntity is an interface that allows working with any entity type
// This is used internally for untyped repository operations
type AnyEntity interface {
	schema.Entity
}

// Repository provides type-safe database operations
type Repository[T AnyEntity] struct {
	db       DBExecutor
	dialect  Dialect
	metadata *schema.EntityMetadata
	ctx      context.Context
}

// NewRepository creates a new repository for the given entity type
func NewRepository[T schema.Entity](db *sql.DB, dialect Dialect) *Repository[T] {
	var entity T
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	meta, exists := schema.Registry.GetEntityMetadata(entityType)
	if !exists {
		panic(fmt.Sprintf("entity %s not registered", entityType.Name()))
	}

	repo := &Repository[T]{
		db:       db,
		dialect:  dialect,
		metadata: meta,
		ctx:      context.Background(),
	}

	return repo
}

// NewUntypedRepository creates a new untyped repository for the given entity type
// This is used internally by the RepositoryProvider
func NewUntypedRepository(entityType reflect.Type, db *sql.DB, d Dialect) interface{} {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	// Create a new instance of the entity type to get its table name
	elem := reflect.New(entityType).Interface()
	_, ok := elem.(schema.Entity)
	if !ok {
		panic(fmt.Sprintf("type %s does not implement schema.Entity", entityType.Name()))
	}

	// Create a repository for the entity type using reflection
	repoType := reflect.TypeOf((*Repository[AnyEntity])(nil))
	repo := reflect.New(repoType.Elem()).Interface().(*Repository[AnyEntity])
	repo.db = db
	repo.dialect = d
	repo.ctx = context.Background()

	// Set the metadata
	meta, exists := schema.Registry.GetEntityMetadata(entityType)
	if !exists {
		panic(fmt.Sprintf("entity %s not registered", entityType.Name()))
	}
	repo.metadata = meta

	return repo
}

// WithContext sets the context for the repository
func (r *Repository[T]) WithContext(ctx context.Context) *Repository[T] {
	return &Repository[T]{
		db:       r.db,
		dialect:  r.dialect,
		metadata: r.metadata,
		ctx:      ctx,
	}
}

// QueryBuilder enables fluent query construction
type QueryBuilder[T schema.Entity] struct {
	repo       *Repository[T]
	conditions []string
	args       []any
	includes   []string
	joins      []JoinClause
	order      string
	limit      int
	offset     int
	groupBy    string
	having     string
	distinct   bool
}

// JoinClause represents a JOIN operation
type JoinClause struct {
	Type      string // "INNER", "LEFT", "RIGHT", "FULL"
	Table     string
	Condition string
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

// With enables eager loading of relationships
func (qb *QueryBuilder[T]) With(relations ...string) *QueryBuilder[T] {
	qb.includes = append(qb.includes, relations...)
	return qb
}

// Include is an alias for With (for backward compatibility)
func (qb *QueryBuilder[T]) Include(relations ...string) *QueryBuilder[T] {
	return qb.With(relations...)
}

// Join adds a JOIN clause to the query
func (qb *QueryBuilder[T]) Join(table, condition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "INNER",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// LeftJoin adds a LEFT JOIN clause to the query
func (qb *QueryBuilder[T]) LeftJoin(table, condition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "LEFT",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// RightJoin adds a RIGHT JOIN clause to the query
func (qb *QueryBuilder[T]) RightJoin(table, condition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "RIGHT",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// FullJoin adds a FULL JOIN clause to the query
func (qb *QueryBuilder[T]) FullJoin(table, condition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "FULL",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// GroupBy sets the GROUP BY clause
func (qb *QueryBuilder[T]) GroupBy(groupBy string) *QueryBuilder[T] {
	qb.groupBy = groupBy
	return qb
}

// Having sets the HAVING clause
func (qb *QueryBuilder[T]) Having(having string, args ...interface{}) *QueryBuilder[T] {
	qb.having = having
	qb.args = append(qb.args, args...)
	return qb
}

// Distinct sets the DISTINCT clause
func (qb *QueryBuilder[T]) Distinct() *QueryBuilder[T] {
	qb.distinct = true
	return qb
}

// WhereIn adds a WHERE IN condition
func (qb *QueryBuilder[T]) WhereIn(column string, values []interface{}) *QueryBuilder[T] {
	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	condition := fmt.Sprintf("%s IN (%s)", qb.repo.dialect.QuoteIdentifier(column), strings.Join(placeholders, ", "))
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, values...)
	return qb
}

// WhereNotIn adds a WHERE NOT IN condition
func (qb *QueryBuilder[T]) WhereNotIn(column string, values []interface{}) *QueryBuilder[T] {
	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", qb.repo.dialect.QuoteIdentifier(column), strings.Join(placeholders, ", "))
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, values...)
	return qb
}

// WhereBetween adds a WHERE BETWEEN condition
func (qb *QueryBuilder[T]) WhereBetween(column string, start, end interface{}) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s BETWEEN ? AND ?", qb.repo.dialect.QuoteIdentifier(column))
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, start, end)
	return qb
}

// WhereLike adds a WHERE LIKE condition
func (qb *QueryBuilder[T]) WhereLike(column, pattern string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s LIKE ?", qb.repo.dialect.QuoteIdentifier(column))
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, pattern)
	return qb
}

// WhereNull adds a WHERE IS NULL condition
func (qb *QueryBuilder[T]) WhereNull(column string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s IS NULL", qb.repo.dialect.QuoteIdentifier(column))
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereNotNull adds a WHERE IS NOT NULL condition
func (qb *QueryBuilder[T]) WhereNotNull(column string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s IS NOT NULL", qb.repo.dialect.QuoteIdentifier(column))
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// OrWhere adds an OR condition
func (qb *QueryBuilder[T]) OrWhere(cond string, args ...interface{}) *QueryBuilder[T] {
	if len(qb.conditions) > 0 {
		// Wrap existing conditions in parentheses and add OR
		qb.conditions = append([]string{"(" + strings.Join(qb.conditions, " AND ") + ")"}, cond)
	} else {
		qb.conditions = append(qb.conditions, cond)
	}
	qb.args = append(qb.args, args...)
	return qb
}

// OrderBy sets the order clause
func (qb *QueryBuilder[T]) OrderBy(order string) *QueryBuilder[T] {
	qb.order = order
	return qb
}

// Limit sets the limit clause
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

// Offset sets the offset clause
func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.offset = offset
	return qb
}

// One returns a single result
func (qb *QueryBuilder[T]) One() (*T, error) {
	qb.limit = 1
	results, err := qb.All()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return &results[0], nil
}

// All returns all results
func (qb *QueryBuilder[T]) All() ([]T, error) {
	query := qb.buildSelectQuery()
	rows, err := qb.repo.db.QueryContext(qb.repo.ctx, query, qb.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return qb.scanRows(rows)
}

// Count returns the count of matching records
func (qb *QueryBuilder[T]) Count() (int64, error) {
	query := qb.buildCountQuery()
	var count int64
	err := qb.repo.db.QueryRowContext(qb.repo.ctx, query, qb.args...).Scan(&count)
	return count, err
}

// buildSelectQuery constructs the SQL query
func (qb *QueryBuilder[T]) buildSelectQuery() string {
	var selects []string

	// Add DISTINCT if specified
	selectKeyword := "SELECT"
	if qb.distinct {
		selectKeyword = "SELECT DISTINCT"
	}

	// Build select columns
	for _, field := range qb.repo.metadata.Fields {
		selects = append(selects, qb.repo.dialect.QuoteIdentifier(field.DBName))
	}

	query := fmt.Sprintf("%s %s FROM %s",
		selectKeyword,
		strings.Join(selects, ", "),
		qb.repo.dialect.QuoteIdentifier(qb.repo.metadata.TableName),
	)

	// Add JOIN clauses
	for _, join := range qb.joins {
		query += fmt.Sprintf(" %s JOIN %s ON %s",
			join.Type,
			qb.repo.dialect.QuoteIdentifier(join.Table),
			join.Condition,
		)
	}

	if len(qb.conditions) > 0 {
		query += " WHERE " + strings.Join(qb.conditions, " AND ")
	}

	if qb.groupBy != "" {
		query += " GROUP BY " + qb.groupBy
	}

	if qb.having != "" {
		query += " HAVING " + qb.having
	}

	if qb.order != "" {
		query += " ORDER BY " + qb.order
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query
}

// buildCountQuery constructs a COUNT query
func (qb *QueryBuilder[T]) buildCountQuery() string {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s",
		qb.repo.dialect.QuoteIdentifier(qb.repo.metadata.TableName),
	)

	if len(qb.conditions) > 0 {
		query += " WHERE " + strings.Join(qb.conditions, " AND ")
	}

	return query
}

// loadRelations loads related entities for eager loading
func (qb *QueryBuilder[T]) loadRelations(results *[]T) error {
	if len(*results) == 0 {
		return nil
	}

	// Get the first entity to determine its type
	firstEntity := (*results)[0]
	entityType := reflect.TypeOf(firstEntity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	// Get entity metadata
	meta, exists := schema.Registry.GetEntityMetadata(entityType)
	if !exists {
		return fmt.Errorf("entity metadata not found for type %s", entityType.Name())
	}

	// Load each requested relation
	for _, relationName := range qb.includes {
		if err := qb.loadRelation(results, meta, relationName); err != nil {
			return err
		}
	}

	return nil
}

// loadRelation loads a specific relation for all entities in the results
func (qb *QueryBuilder[T]) loadRelation(results *[]T, meta *schema.EntityMetadata, relationName string) error {
	// Find the relation metadata
	var relation *schema.RelationMetadata
	for _, rel := range meta.Relations {
		// This is a simplified lookup - in a real implementation, you'd need to match by field name
		if rel.ForeignKey != "" {
			relation = &rel
			break
		}
	}

	if relation == nil {
		return fmt.Errorf("relation '%s' not found in entity %s", relationName, meta.TableName)
	}

	// Get primary key values from results
	var pkValues []interface{}
	resultsValue := reflect.ValueOf(*results)
	for i := 0; i < resultsValue.Len(); i++ {
		entity := resultsValue.Index(i)
		pkField := entity.FieldByName(meta.PrimaryKey.Name)
		if pkField.IsValid() {
			pkValues = append(pkValues, pkField.Interface())
		}
	}

	if len(pkValues) == 0 {
		return nil
	}

	// Load related entities based on relation type
	switch relation.Type {
	case schema.OneToMany:
		return qb.loadOneToManyRelation(results, relation, pkValues)
	case schema.ManyToOne:
		return qb.loadManyToOneRelation(results, relation, pkValues)
	case schema.OneToOne:
		return qb.loadOneToOneRelation(results, relation, pkValues)
	case schema.ManyToMany:
		return qb.loadManyToManyRelation(results, relation, pkValues)
	default:
		return fmt.Errorf("unsupported relation type: %s", relation.Type)
	}
}

// loadOneToManyRelation loads one-to-many relationships
func (qb *QueryBuilder[T]) loadOneToManyRelation(results *[]T, relation *schema.RelationMetadata, pkValues []interface{}) error {

	// 1. Query the related table using the foreign key
	// 2. Group the results by the foreign key
	// 3. Set the related entities on the appropriate parent entities

	// For now, we'll just log that this relation type is supported
	// TODO: Implement full one-to-many loading logic
	return nil
}

// loadManyToOneRelation loads many-to-one relationships
func (qb *QueryBuilder[T]) loadManyToOneRelation(results *[]T, relation *schema.RelationMetadata, pkValues []interface{}) error {

	// 1. Query the related table using the primary key
	// 2. Set the related entity on the appropriate parent entity

	// For now, we'll just log that this relation type is supported
	// TODO: Implement full many-to-one loading logic
	return nil
}

// loadOneToOneRelation loads one-to-one relationships
func (qb *QueryBuilder[T]) loadOneToOneRelation(results *[]T, relation *schema.RelationMetadata, pkValues []interface{}) error {

	// 1. Query the related table using the foreign key
	// 2. Set the related entity on the appropriate parent entity

	// For now, we'll just log that this relation type is supported
	// TODO: Implement full one-to-one loading logic
	return nil
}

// loadManyToManyRelation loads many-to-many relationships
func (qb *QueryBuilder[T]) loadManyToManyRelation(results *[]T, relation *schema.RelationMetadata, pkValues []interface{}) error {

	// 1. Query the join table using the foreign key
	// 2. Query the related table using the reference key
	// 3. Set the related entities on the appropriate parent entity

	// For now, we'll just log that this relation type is supported
	// TODO: Implement full many-to-many loading logic
	return nil
}

// scanRows scans rows into entity structs
func (qb *QueryBuilder[T]) scanRows(rows *sql.Rows) ([]T, error) {
	var results []T

	// Get column names from result
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a map of column name to field index
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	for rows.Next() {
		// Create a new entity instance
		var entity T
		entityValue := reflect.ValueOf(&entity).Elem()

		// Create a slice of pointers to scan into
		scanValues := make([]interface{}, len(columns))
		for i := range scanValues {
			scanValues[i] = new(interface{})
		}

		// Scan the row into the slice
		if err := rows.Scan(scanValues...); err != nil {
			return nil, err
		}

		// Set the values on the entity
		for _, field := range qb.repo.metadata.Fields {
			colIdx, ok := columnMap[field.DBName]
			if !ok {
				continue
			}

			fieldValue := entityValue.FieldByName(field.Name)
			if !fieldValue.IsValid() || !fieldValue.CanSet() {
				continue
			}

			value := *(scanValues[colIdx].(*interface{}))
			if value == nil {
				continue
			}

			// Convert the value to the field type
			convertedValue := reflect.ValueOf(value)
			if convertedValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(convertedValue.Convert(fieldValue.Type()))
			}
		}

		results = append(results, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load relations if requested
	if len(qb.includes) > 0 {
		if err := qb.loadRelations(&results); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// FindByID finds an entity by its primary key
func (r *Repository[T]) FindByID(id interface{}) (*T, error) {
	if r.metadata.PrimaryKey == nil {
		return nil, errors.New("entity has no primary key")
	}

	return r.Find().Where(
		fmt.Sprintf("%s = ?", r.dialect.QuoteIdentifier(r.metadata.PrimaryKey.DBName)),
		id,
	).One()
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

// insert creates a new record
func (r *Repository[T]) insert(entity *T) error {
	meta := r.metadata
	val := reflect.ValueOf(entity).Elem()

	var columns []string
	var placeholders []string
	var values []interface{}

	for i, field := range meta.Fields {
		// Skip auto-increment primary key for insert
		if field.IsPrimaryKey && field.IsAutoIncr {
			continue
		}

		// Skip relation fields
		if field.Relation != nil {
			continue
		}

		columns = append(columns, r.dialect.QuoteIdentifier(field.DBName))
		placeholders = append(placeholders, r.dialect.Placeholder(i))

		fieldValue := val.FieldByName(field.Name)
		values = append(values, fieldValue.Interface())
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.dialect.QuoteIdentifier(meta.TableName),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	var result sql.Result
	var err error

	if meta.PrimaryKey != nil && meta.PrimaryKey.IsAutoIncr {
		// Execute and get last insert ID
		result, err = r.db.ExecContext(r.ctx, query, values...)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Set the ID on the entity
		pkField := val.FieldByName(meta.PrimaryKey.Name)
		if pkField.CanSet() {
			// Handle different types of primary key fields
			switch pkField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				pkField.SetInt(id)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				pkField.SetUint(uint64(id))
			default:
				return fmt.Errorf("unsupported primary key type: %s", pkField.Type())
			}
		}
	} else {
		// Just execute without getting ID
		_, err = r.db.ExecContext(r.ctx, query, values...)
	}

	return err
}

// update updates an existing record
func (r *Repository[T]) update(entity *T) error {
	meta := r.metadata
	val := reflect.ValueOf(entity).Elem()

	var setColumns []string
	var values []interface{}

	for _, field := range meta.Fields {
		// Skip primary key and relation fields for update SET clause
		if field.IsPrimaryKey || field.Relation != nil {
			continue
		}

		setColumns = append(setColumns,
			fmt.Sprintf("%s = ?", r.dialect.QuoteIdentifier(field.DBName)))

		fieldValue := val.FieldByName(field.Name)
		values = append(values, fieldValue.Interface())
	}

	// Add primary key value for WHERE clause
	pkValue := val.FieldByName(meta.PrimaryKey.Name)
	values = append(values, pkValue.Interface())

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?",
		r.dialect.QuoteIdentifier(meta.TableName),
		strings.Join(setColumns, ", "),
		r.dialect.QuoteIdentifier(meta.PrimaryKey.DBName),
	)

	_, err := r.db.ExecContext(r.ctx, query, values...)
	return err
}

// Delete deletes an entity
func (r *Repository[T]) Delete(entity *T) error {
	meta := r.metadata
	if meta.PrimaryKey == nil {
		return errors.New("entity missing primary key")
	}

	val := reflect.ValueOf(entity).Elem()
	pkValue := val.FieldByName(meta.PrimaryKey.Name)

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?",
		r.dialect.QuoteIdentifier(meta.TableName),
		r.dialect.QuoteIdentifier(meta.PrimaryKey.DBName),
	)

	_, err := r.db.ExecContext(r.ctx, query, pkValue.Interface())
	return err
}

// DeleteByID deletes an entity by its primary key
func (r *Repository[T]) DeleteByID(id interface{}) error {
	meta := r.metadata
	if meta.PrimaryKey == nil {
		return errors.New("entity missing primary key")
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?",
		r.dialect.QuoteIdentifier(meta.TableName),
		r.dialect.QuoteIdentifier(meta.PrimaryKey.DBName),
	)

	_, err := r.db.ExecContext(r.ctx, query, id)
	return err
}

// Transaction executes a database transaction
func (r *Repository[T]) Transaction(fn func(*Repository[T]) error) error {
	// We need to cast r.db to *sql.DB to use BeginTx
	db, ok := r.db.(*sql.DB)
	if !ok {
		return errors.New("cannot start a transaction: db is not a *sql.DB")
	}

	tx, err := db.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	// Create a new repository with the transaction
	txRepo := &Repository[T]{
		db:       tx, // Use the transaction as a DBExecutor
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

// Hook interfaces for entity lifecycle events
type (
	BeforeCreateHook interface {
		BeforeCreate() error
	}

	AfterCreateHook interface {
		AfterCreate() error
	}

	BeforeUpdateHook interface {
		BeforeUpdate() error
	}

	AfterUpdateHook interface {
		AfterUpdate() error
	}

	BeforeDeleteHook interface {
		BeforeDelete() error
	}

	AfterDeleteHook interface {
		AfterDelete() error
	}

	BeforeSaveHook interface {
		BeforeSave() error
	}

	AfterSaveHook interface {
		AfterSave() error
	}
)
