package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps sql.DB with helpers.
type DB struct {
	*sql.DB
}

// New opens (or creates) the SQLite database at path, then runs migrations.
func New(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Single writer to avoid SQLite locking under concurrent HTTP requests.
	sqlDB.SetMaxOpenConns(1)

	db := &DB{sqlDB}
	if err := db.Migrate(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

// Migrate applies the schema DDL idempotently.
func (db *DB) Migrate() error {
	_, err := db.Exec(schema)
	return err
}
