// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// User represents a system user.
type User struct {
	ID           string
	Username     string
	PasswordHash string
	RoleID       string
}

// Role represents a collection of permissions.
type Role struct {
	ID          string
	Name        string
	Permissions string // JSON array of permission strings
}

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// GetUserByUsername retrieves a user by their username.
func (db *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password_hash, role_id FROM users WHERE username = ?`
	var u User
	err := db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("store: get user: %w", err)
	}
	return &u, nil
}

// CreateUser inserts a new user into the database.
func (db *DB) CreateUser(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, username, password_hash, role_id) VALUES (?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, u.ID, u.Username, u.PasswordHash, u.RoleID)
	if err != nil {
		return fmt.Errorf("store: create user: %w", err)
	}
	return nil
}

// CreateRole inserts a new role into the database.
func (db *DB) CreateRole(ctx context.Context, r *Role) error {
	query := `INSERT INTO roles (id, name, permissions_json) VALUES (?, ?, ?)`
	_, err := db.ExecContext(ctx, query, r.ID, r.Name, r.Permissions)
	if err != nil {
		return fmt.Errorf("store: create role: %w", err)
	}
	return nil
}

// SeedAdmin seeds the initial super admin account if the database is empty.
func (db *DB) SeedAdmin(ctx context.Context, username, hash string) error {
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	role := &Role{
		ID:          "role-admin",
		Name:        "Super Admin",
		Permissions: `["*"]`,
	}
	if err := db.CreateRole(ctx, role); err != nil {
		return err
	}

	user := &User{
		ID:           "user-admin",
		Username:     username,
		PasswordHash: hash,
		RoleID:       "role-admin",
	}
	return db.CreateUser(ctx, user)
}
