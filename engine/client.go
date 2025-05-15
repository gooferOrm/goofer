package engine

import (
    "database/sql"
    "fmt"

    "github.com/gooferOrm/goofer/dialect"
    "github.com/gooferOrm/goofer/schema"
)

// Client is your one stop Goofer engine.
// It implements the RepositoryProvider interface.
type Client struct {
    db      *sql.DB
    dialect dialect.Dialect
}

// Ensure Client implements RepositoryProvider
var _ RepositoryProvider = (*Client)(nil)

// NewClient creates a new Goofer client with the provided database connection and dialect.
// It can optionally register and auto-migrate the provided entities.
//
// Example:
//   db, _ := sql.Open("sqlite3", "test.db")
//   client, err := NewClient(db, &dialect.SQLite{}, &User{}, &Product{})
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer client.Close()
func NewClient(
    db *sql.DB,
    d dialect.Dialect,
    entities ...schema.Entity,
) (*Client, error) {
    client := &Client{db: db, dialect: d}
    
    if len(entities) > 0 {
        if err := client.RegisterEntities(entities...); err != nil {
            return nil, fmt.Errorf("failed to register entities: %w", err)
        }
    }
    
    return client, nil
}

// Close closes the underlying database connection
func (c *Client) Close() error {
    return c.db.Close()
}

// RegisterEntities registers multiple entities with the schema registry and optionally auto-migrates them
func (c *Client) RegisterEntities(entities ...schema.Entity) error {
    // Register entities
    for _, e := range entities {
        if err := schema.Registry.RegisterEntity(e); err != nil {
            return fmt.Errorf("register %T: %w", e, err)
        }
    }

    // Auto-migrate
    for _, e := range entities {
        meta, ok := schema.Registry.GetEntityMetadata(schema.GetEntityType(e))
        if !ok {
            return fmt.Errorf("no metadata for %T", e)
        }
        ddl := c.dialect.CreateTableSQL(meta)
        if _, err := c.db.Exec(ddl); err != nil {
            return fmt.Errorf("migrate %s: %w", meta.TableName, err)
        }
    }
    return nil
}
