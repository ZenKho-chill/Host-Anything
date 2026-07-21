// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/host-anything/hostanything/internal/store"
)

// SettingsHandler handles settings management.
type SettingsHandler struct {
	db     *store.DB
	logger *slog.Logger
}

// NewSettingsHandler creates a new settings handler.
func NewSettingsHandler(db *store.DB, logger *slog.Logger) *SettingsHandler {
	return &SettingsHandler{db: db, logger: logger}
}

// GetSettings handles GET /api/v1/settings.
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	timezone, _ := h.db.GetSetting(r.Context(), "timezone")

	settings := map[string]string{
		"timezone": timezone,
	}
	writeJSON(w, http.StatusOK, settings)
}

// UpdateSettings handles POST /api/v1/settings.
func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	for k, v := range req {
		if err := h.db.SetSetting(r.Context(), k, v); err != nil {
			h.logger.Error("failed to set setting", "key", k, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update settings")
			return
		}
	}

	writeJSON(w, http.StatusOK, req)
}
