// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package store

import (
	"context"
	"fmt"
)

// Schedule represents a background cron job.
type Schedule struct {
	ID        string `json:"id"`
	TaskName  string `json:"task_name"`
	CronExpr  string `json:"cron_expr"`
	Command   string `json:"command"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
}

// ListSchedules returns all schedules.
func (db *DB) ListSchedules(ctx context.Context) ([]Schedule, error) {
	query := `SELECT id, task_name, cron_expr, command, enabled, created_at FROM schedules ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("store: list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.ID, &s.TaskName, &s.CronExpr, &s.Command, &s.Enabled, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// CreateSchedule adds a new schedule.
func (db *DB) CreateSchedule(ctx context.Context, s *Schedule) error {
	query := `INSERT INTO schedules (id, task_name, cron_expr, command, enabled) VALUES (?, ?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, s.ID, s.TaskName, s.CronExpr, s.Command, s.Enabled)
	if err != nil {
		return fmt.Errorf("store: create schedule: %w", err)
	}
	return nil
}

// DeleteSchedule removes a schedule.
func (db *DB) DeleteSchedule(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM schedules WHERE id = ?", id)
	return err
}
