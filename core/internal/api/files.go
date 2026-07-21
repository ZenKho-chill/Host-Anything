// Copyright 2026 Host Anything Contributors
// Licensed under the Apache License, Version 2.0 (the "License")

package api

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// FileHandler handles file management within the sandbox.
type FileHandler struct {
	baseDir string
	logger  *slog.Logger
}

// NewFileHandler creates a new file handler sandboxed to baseDir.
func NewFileHandler(baseDir string, logger *slog.Logger) *FileHandler {
	return &FileHandler{baseDir: baseDir, logger: logger}
}

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

// ListFiles handles GET /api/v1/files/*
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	subPath := chi.URLParam(r, "*")

	targetPath := filepath.Join(h.baseDir, subPath)

	// Prevent path traversal outside baseDir
	absBase, _ := filepath.Abs(h.baseDir)
	absTarget, _ := filepath.Abs(targetPath)
	if !strings.HasPrefix(absTarget, absBase) {
		writeJSONError(w, http.StatusForbidden, "path traversal blocked")
		return
	}

	stat, err := os.Stat(absTarget)
	if err != nil {
		if os.IsNotExist(err) {
			writeJSONError(w, http.StatusNotFound, "file or directory not found")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "error reading path")
		}
		return
	}

	if !stat.IsDir() {
		// Serve file
		http.ServeFile(w, r, absTarget)
		return
	}

	// List directory
	entries, err := os.ReadDir(absTarget)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "error reading directory")
		return
	}

	var files []FileInfo
	for _, e := range entries {
		info, _ := e.Info()
		files = append(files, FileInfo{
			Name:  e.Name(),
			IsDir: e.IsDir(),
			Size:  info.Size(),
		})
	}

	writeJSON(w, http.StatusOK, files)
}
