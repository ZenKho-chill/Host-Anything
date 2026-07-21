// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SetSetting sets a key-value setting in the database.
func (db *DB) SetSetting(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.ExecContext(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("store: set setting: %w", err)
	}
	return nil
}

// GetSetting retrieves a setting by key. Returns empty string if not found.
func (db *DB) GetSetting(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = ?`
	var val string
	err := db.QueryRowContext(ctx, query, key).Scan(&val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("store: get setting: %w", err)
	}
	return val, nil
}
