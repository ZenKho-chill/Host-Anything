// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package store

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB represents a connection to the SQLite database.
type DB struct {
	*sql.DB
	logger *slog.Logger
}

// Open initializes the SQLite database at the given data directory.
func Open(dataDir string, logger *slog.Logger) (*DB, error) {
	dbPath := filepath.Join(dataDir, "hostanything.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("store: open db: %w", err)
	}

	// Optimize SQLite for concurrent access
	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	db.Exec("PRAGMA foreign_keys=ON;")

	dbWrapper := &DB{DB: db, logger: logger}

	if err := dbWrapper.migrate(); err != nil {
		return nil, fmt.Errorf("store: migrate: %w", err)
	}

	return dbWrapper, nil
}

// migrate creates tables if they don't exist.
func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			permissions_json TEXT NOT NULL DEFAULT '[]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE RESTRICT
		);`,
		`CREATE TABLE IF NOT EXISTS schedules (
			id TEXT PRIMARY KEY,
			task_name TEXT NOT NULL,
			cron_expr TEXT NOT NULL,
			command TEXT NOT NULL,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
