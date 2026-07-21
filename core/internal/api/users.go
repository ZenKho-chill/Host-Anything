// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package api

import (
	"log/slog"
	"net/http"

	"github.com/host-anything/hostanything/internal/store"
)

// UserHandler handles user management.
type UserHandler struct {
	db     *store.DB
	logger *slog.Logger
}

// NewUserHandler creates a new user handler.
func NewUserHandler(db *store.DB, logger *slog.Logger) *UserHandler {
	return &UserHandler{db: db, logger: logger}
}

// ListUsers handles GET /api/v1/users.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// For this milestone, we just return an empty list or the current user.
	// Implementing full CRUD for users requires more DB methods.
	writeJSON(w, http.StatusOK, []store.User{})
}

// CreateUser handles POST /api/v1/users.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotImplemented, "not implemented")
}
