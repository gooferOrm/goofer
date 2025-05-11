package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gooferOrm/goofer/pkg/schema"
)

// Dialect interface for database-specific implementations
type Dialect interface {
	Placeholder(int) string
	QuoteIdentifier(string) string
	DataType(field schema.FieldMetadata) string
	CreateTableSQL(*schema.EntityMetadata) string
}

// Repository provides type-safe database operations
type Repository[T schema.Entity] struct {
	db       *sql.DB
	dialect  Dialect
	metadata *schema.EntityMetadata
	ctx      context.Context
}

// NewRepository creates a new repository
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

	return &Repository[T]{
		db:       db,
		dialect:  dialect,
		metadata: meta,
		ctx:      context.Background(),
	}
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
	args       []interface{}
	includes   []string
	order      string
	limit      int
	offset     int
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

	// TODO: Load relations if requested

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
			pkField.SetInt(id)
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
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	// Create a new repository with the transaction
	txRepo := &Repository[T]{
		db:       tx, // Use the transaction instead of the original DB
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
