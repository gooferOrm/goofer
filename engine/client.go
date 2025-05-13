package engine

import (
    "database/sql"
    "fmt"

    "github.com/gooferOrm/goofer/dialect"
    "github.com/gooferOrm/goofer/repository"
    "github.com/gooferOrm/goofer/schema"
)

// Client is your one stop Goofer engine.
type Client struct {
    db      *sql.DB
    dialect dialect.Dialect
}

// NewClient will:
// 1. register all your entities
// 2. auto-create their tables
// 3. give you a Repo[T] accessor
func NewClient(
    db *sql.DB,
    d dialect.Dialect,
    entities ...interface{},
) (*Client, error) {
    //register entities
    for _, e := range entities {
        if err := schema.Registry.RegisterEntity(e); err != nil {
            return nil, fmt.Errorf("register %T: %w", e, err)
        }
    }

    //auto-migrate
    for _, e := range entities {
        meta, ok := schema.Registry.GetEntityMetadata(schema.GetEntityType(e))
        if !ok {
            return nil, fmt.Errorf("no metadata for %T", e)
        }
        ddl := d.CreateTableSQL(meta)
        if _, err := db.Exec(ddl); err != nil {
            return nil, fmt.Errorf("migrate %s: %w", meta.TableName, err)
        }
    }

    return &Client{db: db, dialect: d}, nil
}

// Repo[T] gives you a fully wired Repository[T].
func (c *Client) Repo[T any]() *repository.Repository[T] {
    return repository.NewRepository[T](c.db, c.dialect)
}
