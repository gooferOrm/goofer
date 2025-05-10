package migration

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gooferOrm/goofer/pkg/repository"
	"github.com/gooferOrm/goofer/pkg/schema"
)

// Migration represents a database migration
type Migration struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	AppliedAt time.Time `db:"applied_at"`
	Script    string    `db:"script"`
	Checksum  string    `db:"checksum"`
}

// MigrationScript represents a migration script
type MigrationScript struct {
	Up   string
	Down string
}

// Migrator handles database migrations
type Migrator struct {
	db      *sql.DB
	dialect repository.Dialect
	outPath string
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, dialect repository.Dialect, outPath string) *Migrator {
	return &Migrator{
		db:      db,
		dialect: dialect,
		outPath: outPath,
	}
}

// ensureMigrationTable creates the migration table if it doesn't exist
func (m *Migrator) ensureMigrationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL,
		script TEXT NOT NULL,
		checksum VARCHAR(32) NOT NULL
	);`

	_, err := m.db.Exec(query)
	return err
}

// Up runs pending migrations
func (m *Migrator) Up() error {
	if err := m.ensureMigrationTable(); err != nil {
		return err
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	// Get available migrations
	available, err := m.getAvailableMigrations()
	if err != nil {
		return err
	}

	// Find pending migrations
	pending := m.getPendingMigrations(applied, available)
	if len(pending) == 0 {
		fmt.Println("No pending migrations")
		return nil
	}

	// Sort migrations by ID
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].ID < pending[j].ID
	})

	// Run pending migrations
	for _, migration := range pending {
		fmt.Printf("Running migration: %s\n", migration.Name)

		// Begin transaction
		tx, err := m.db.Begin()
		if err != nil {
			return err
		}

		// Execute migration script
		_, err = tx.Exec(migration.Script)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %w", migration.ID, err)
		}

		// Record migration
		_, err = tx.Exec(
			"INSERT INTO migrations (id, name, applied_at, script, checksum) VALUES (?, ?, ?, ?, ?)",
			migration.ID,
			migration.Name,
			time.Now(),
			migration.Script,
			migration.Checksum,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", migration.ID, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %w", migration.ID, err)
		}

		fmt.Printf("Migration applied: %s\n", migration.Name)
	}

	return nil
}

// Down reverts the last migration
func (m *Migrator) Down() error {
	if err := m.ensureMigrationTable(); err != nil {
		return err
	}

	// Get the last applied migration
	var migration Migration
	err := m.db.QueryRow(`
		SELECT id, name, applied_at, script, checksum
		FROM migrations
		ORDER BY applied_at DESC
		LIMIT 1
	`).Scan(&migration.ID, &migration.Name, &migration.AppliedAt, &migration.Script, &migration.Checksum)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no migrations to revert")
		}
		return err
	}

	// Get the down script
	downScript, err := m.getDownScript(migration.ID)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	// Execute down script
	_, err = tx.Exec(downScript)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing down migration %s: %w", migration.ID, err)
	}

	// Delete migration record
	_, err = tx.Exec("DELETE FROM migrations WHERE id = ?", migration.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting migration record %s: %w", migration.ID, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing down migration %s: %w", migration.ID, err)
	}

	fmt.Printf("Migration reverted: %s\n", migration.Name)
	return nil
}

// Status shows the migration status
func (m *Migrator) Status() ([]Migration, error) {
	if err := m.ensureMigrationTable(); err != nil {
		return nil, err
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	// Get available migrations
	available, err := m.getAvailableMigrations()
	if err != nil {
		return nil, err
	}

	// Find pending migrations
	pending := m.getPendingMigrations(applied, available)

	// Print status
	fmt.Println("Applied migrations:")
	for _, migration := range applied {
		fmt.Printf("  %s - %s (applied at %s)\n", migration.ID, migration.Name, migration.AppliedAt)
	}

	fmt.Println("\nPending migrations:")
	for _, migration := range pending {
		fmt.Printf("  %s - %s\n", migration.ID, migration.Name)
	}

	return applied, nil
}

// getAppliedMigrations returns the list of applied migrations
func (m *Migrator) getAppliedMigrations() ([]Migration, error) {
	rows, err := m.db.Query(`
		SELECT id, name, applied_at, script, checksum
		FROM migrations
		ORDER BY applied_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.ID, &migration.Name, &migration.AppliedAt, &migration.Script, &migration.Checksum)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, rows.Err()
}

// getAvailableMigrations returns the list of available migrations
func (m *Migrator) getAvailableMigrations() ([]Migration, error) {
	files, err := ioutil.ReadDir(m.outPath)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		id := strings.Split(file.Name(), "_")[0]
		name := strings.TrimSuffix(file.Name(), ".up.sql")
		name = strings.Replace(name, id+"_", "", 1)

		script, err := ioutil.ReadFile(filepath.Join(m.outPath, file.Name()))
		if err != nil {
			return nil, err
		}

		checksum := md5.Sum(script)

		migrations = append(migrations, Migration{
			ID:       id,
			Name:     name,
			Script:   string(script),
			Checksum: hex.EncodeToString(checksum[:]),
		})
	}

	return migrations, nil
}

// getPendingMigrations returns the list of pending migrations
func (m *Migrator) getPendingMigrations(applied, available []Migration) []Migration {
	appliedMap := make(map[string]bool)
	for _, migration := range applied {
		appliedMap[migration.ID] = true
	}

	var pending []Migration
	for _, migration := range available {
		if !appliedMap[migration.ID] {
			pending = append(pending, migration)
		}
	}

	return pending
}

// getDownScript returns the down script for a migration
func (m *Migrator) getDownScript(id string) (string, error) {
	files, err := ioutil.ReadDir(m.outPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), id) || !strings.HasSuffix(file.Name(), ".down.sql") {
			continue
		}

		script, err := ioutil.ReadFile(filepath.Join(m.outPath, file.Name()))
		if err != nil {
			return "", err
		}

		return string(script), nil
	}

	return "", fmt.Errorf("down script not found for migration %s", id)
}

// MigrationGenerator generates migration files
type MigrationGenerator struct {
	Registry  *schema.SchemaRegistry
	Dialect   repository.Dialect
	OutPath   string
}

// Generate creates a new migration file
func (g *MigrationGenerator) Generate(name string) error {
	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(g.OutPath, 0755); err != nil {
		return err
	}

	// Generate timestamp for migration ID
	timestamp := time.Now().Format("20060102150405")

	// Generate migration scripts
	script, err := g.generateMigrationScript()
	if err != nil {
		return err
	}

	// Write up script
	upFilename := filepath.Join(g.OutPath, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	if err := ioutil.WriteFile(upFilename, []byte(script.Up), 0644); err != nil {
		return err
	}

	// Write down script
	downFilename := filepath.Join(g.OutPath, fmt.Sprintf("%s_%s.down.sql", timestamp, name))
	if err := ioutil.WriteFile(downFilename, []byte(script.Down), 0644); err != nil {
		return err
	}

	fmt.Printf("Generated migration: %s\n", name)
	return nil
}

// generateMigrationScript generates migration scripts from entity metadata
func (g *MigrationGenerator) generateMigrationScript() (*MigrationScript, error) {
	var upBuilder strings.Builder
	var downBuilder strings.Builder

	// Get all entity metadata
	for _, meta := range g.Registry.GetAllEntities() {
		// Generate CREATE TABLE statement
		createTable := g.Dialect.CreateTableSQL(meta)
		upBuilder.WriteString(createTable)
		upBuilder.WriteString("\n\n")

		// Generate DROP TABLE statement
		dropTable := fmt.Sprintf("DROP TABLE IF EXISTS %s;", g.Dialect.QuoteIdentifier(meta.TableName))
		downBuilder.WriteString(dropTable)
		downBuilder.WriteString("\n\n")
	}

	return &MigrationScript{
		Up:   upBuilder.String(),
		Down: downBuilder.String(),
	}, nil
}
