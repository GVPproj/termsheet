package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const DBFile = "termsheet.db"

var db *sql.DB

// InitDB initializes the database connection and creates tables if they don't exist
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", DBFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables creates all required tables
func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS provider (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT,
			email TEXT,
			phone TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS client (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			address TEXT,
			email TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS invoice (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_id TEXT NOT NULL,
			client_id TEXT NOT NULL,
			paid BOOLEAN DEFAULT FALSE,
			date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (provider_id) REFERENCES provider (id),
			FOREIGN KEY (client_id) REFERENCES client (id)
		)`,
		`CREATE TABLE IF NOT EXISTS invoice_item (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			item_name TEXT NOT NULL,
			amount REAL NOT NULL,
			cost_per_unit REAL NOT NULL,
			FOREIGN KEY (invoice_id) REFERENCES invoice (id)
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
