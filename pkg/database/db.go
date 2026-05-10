package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps sql.DB with helpers.
type DB struct {
	*sql.DB
}


func New(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to Create Database: %w", err)
	}

	schemaCreate := `
	PRAGMA journal_mode=WAL;
	PRAGMA foreign_keys=ON;

	CREATE TABLE IF NOT EXISTS users (
		id           TEXT PRIMARY KEY,
		username     TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS manga (
		id             TEXT PRIMARY KEY,
		title          TEXT NOT NULL,
		author         TEXT NOT NULL DEFAULT '',
		genres         TEXT NOT NULL DEFAULT '[]',
		status         TEXT NOT NULL DEFAULT 'ongoing',
		total_chapters INTEGER NOT NULL DEFAULT 0,
		description    TEXT NOT NULL DEFAULT '',
		cover_url      TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		user_id         TEXT NOT NULL,
		manga_id        TEXT NOT NULL,
		current_chapter INTEGER NOT NULL DEFAULT 0,
		status          TEXT NOT NULL DEFAULT 'plan-to-read',
		rating          INTEGER NOT NULL DEFAULT 0,
		notes           TEXT NOT NULL DEFAULT '',
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (manga_id) REFERENCES manga(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title);
	CREATE INDEX IF NOT EXISTS idx_manga_status ON manga(status);
	CREATE INDEX IF NOT EXISTS idx_progress_user ON user_progress(user_id);
	CREATE INDEX IF NOT EXISTS idx_progress_status ON user_progress(status);
	`
	_, err = sqlDB.Exec(schemaCreate)
	if err != nil {
		return nil, fmt.Errorf("Failed to Create Database: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to Connect Database: %w", err)
	}
	
	return sqlDB, nil
	// // Single writer to avoid SQLite locking under concurrent HTTP requests.
	// sqlDB.SetMaxOpenConns(1)
	
	// db := &DB{sqlDB}
	// if err := db.Migrate(); err != nil {
	// 	sqlDB.Close()
	// 	return nil, fmt.Errorf("migrate: %w", err)
	// }
}

// Migrate applies the schema DDL idempotently.
func (db *DB) Migrate() error {
	_, err := db.Exec(schema)
	return err
}