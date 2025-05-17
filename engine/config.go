package engine

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/schema"
)

// Config holds the database configuration
type Config struct {
	Driver   string
	DSN      string
	LogLevel string // "debug", "info", "error"
	// RegisterEntities func(entities []schema.Entity)
}

func (c Config) RegisterEntities(entities []schema.Entity)error {
	return nil
}

// NewConfig creates a new database configuration with sensible defaults
func NewConfig(driver, dsn string) *Config {
	return &Config{
		Driver:   driver,
		DSN:      dsn,
		LogLevel: "error",
	}
}

// WithLogLevel sets the log level
func (c *Config) WithLogLevel(level string) *Config {
	c.LogLevel = level
	return c
}

// Connect creates a new database connection with the given configuration
func (c *Config) Connect() (*Client, error) {
	db, err := sql.Open(c.Driver, c.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Create appropriate dialect based on driver
	var d dialect.Dialect
	switch strings.ToLower(c.Driver) {
	case "sqlite3":
		d = &dialect.SQLiteDialect{}
	case "postgres":
		d = &dialect.PostgresDialect{}
	case "mysql":
		d = &dialect.MySQLDialect{}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", c.Driver)
	}
	return &Client{db: db, dialect: d}, nil
}

// Connect is a convenience function for quick database connection
func Connect(driver, dsn string) (*Client, error) {
	return NewConfig(driver, dsn).Connect()
}
