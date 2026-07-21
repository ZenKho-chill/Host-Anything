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

	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

// TemplateListHandler returns an HTTP handler that lists all available templates.
func TemplateListHandler(reg *template.Registry, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summaries, err := reg.List()
		if err != nil {
			logger.Error("failed to list templates", "error", err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to retrieve templates")
			return
		}

		if summaries == nil {
			summaries = make([]types.TemplateSummary, 0) // Ensure empty array instead of null in JSON
		}

		writeJSON(w, http.StatusOK, summaries)
	}
}
