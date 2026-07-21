// Copyright 2026 Host Anything Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/template"
)

// TemplateGetHandler returns an HTTP handler that fetches a specific template by name.
func TemplateGetHandler(reg *template.Registry, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if name == "" {
			writeJSONError(w, http.StatusBadRequest, "Template name is required")
			return
		}

		// Optional version query parameter (e.g. ?version=1.0.0)
		version := r.URL.Query().Get("version")

		tmpl, err := reg.Get(name, version)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, http.StatusNotFound, err.Error())
				return
			}
			logger.Error("failed to get template", "name", name, "version", version, "error", err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to retrieve template")
			return
		}

		// Security: Mask default values for secret variables before returning to the UI
		for i := range tmpl.Config {
			if tmpl.Config[i].Type == template.ConfigTypeSecret && tmpl.Config[i].Default != nil {
				tmpl.Config[i].Default = "***"
			}
		}

		writeJSON(w, http.StatusOK, tmpl)
	}
}
