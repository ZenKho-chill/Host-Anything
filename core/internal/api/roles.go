// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package api

import (
	"log/slog"
	"net/http"

	"github.com/host-anything/hostanything/internal/store"
)

// RoleHandler handles role management.
type RoleHandler struct {
	db     *store.DB
	logger *slog.Logger
}

// NewRoleHandler creates a new role handler.
func NewRoleHandler(db *store.DB, logger *slog.Logger) *RoleHandler {
	return &RoleHandler{db: db, logger: logger}
}

// ListRoles handles GET /api/v1/roles.
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []store.Role{})
}
