// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/store"
)

// ScheduleHandler handles schedule management.
type ScheduleHandler struct {
	db     *store.DB
	logger *slog.Logger
}

// NewScheduleHandler creates a new schedule handler.
func NewScheduleHandler(db *store.DB, logger *slog.Logger) *ScheduleHandler {
	return &ScheduleHandler{db: db, logger: logger}
}

// ListSchedules handles GET /api/v1/schedules.
func (h *ScheduleHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.db.ListSchedules(r.Context())
	if err != nil {
		h.logger.Error("failed to list schedules", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if schedules == nil {
		schedules = []store.Schedule{}
	}
	writeJSON(w, http.StatusOK, schedules)
}

// CreateSchedule handles POST /api/v1/schedules.
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var s store.Schedule
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := h.db.CreateSchedule(r.Context(), &s); err != nil {
		h.logger.Error("failed to create schedule", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to create schedule")
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

// DeleteSchedule handles DELETE /api/v1/schedules/{id}.
func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.db.DeleteSchedule(r.Context(), id); err != nil {
		h.logger.Error("failed to delete schedule", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to delete schedule")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
