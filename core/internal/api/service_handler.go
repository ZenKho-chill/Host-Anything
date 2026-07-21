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
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/host-anything/hostanything/internal/runtime"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/errors"
	"github.com/host-anything/hostanything/pkg/types"
)

// DeployRequest represents the payload for deploying a new service.
type DeployRequest struct {
	TemplateName string            `json:"template_name"`
	Config       map[string]string `json:"config"`
}

// DeployResponse is returned on successful deployment.
type DeployResponse struct {
	ServiceID string `json:"service_id"`
}

// ServiceHandler struct holds the dependencies for service endpoints.
type ServiceHandler struct {
	Manager  *runtime.ServiceManager
	Registry *template.Registry
	Logger   *slog.Logger
	Key      []byte
}

// ListServices handles GET /api/v1/services.
func (h *ServiceHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	services := h.Manager.ListServices()

	// Convert map to a list for JSON
	type Svc struct {
		ID    string             `json:"id"`
		State types.ServiceState `json:"state"`
	}

	var list []Svc
	for id, state := range services {
		list = append(list, Svc{ID: id, State: state})
	}

	writeJSON(w, http.StatusOK, list)
}

// DeployService handles POST /api/v1/services.
func (h *ServiceHandler) DeployService(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if req.TemplateName == "" {
		writeJSONError(w, http.StatusBadRequest, "template_name is required")
		return
	}

	// 1. Fetch Template
	tmpl, err := h.Registry.Get(req.TemplateName, "latest")
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "template not found")
		} else {
			h.Logger.Error("failed to get template", "error", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to get template")
		}
		return
	}

	// 2. Resolve Config
	resolved, err := template.Resolve(tmpl, req.Config, h.Key)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "config resolution failed: "+err.Error())
		return
	}

	// 3. Generate ID and construct Spec
	serviceID := uuid.New().String()[:8] // simple 8-char ID for now
	spec := &types.ServiceSpec{
		ServiceID:   serviceID,
		Template:    tmpl,
		ResolvedEnv: resolved,
	}

	// 4. Deploy (runs asynchronously so we don't block the HTTP request if it takes long)
	go func() {
		if err := h.Manager.DeployService(spec); err != nil {
			h.Logger.Error("async deploy failed", "service", serviceID, "error", err)
		}
	}()

	writeJSON(w, http.StatusAccepted, DeployResponse{ServiceID: serviceID})
}

// LogsService handles GET /api/v1/services/{id}/logs.
func (h *ServiceHandler) LogsService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "id")
	// For M4, we assume docker runtime. In M5 we'd dynamically lookup which runtime owns it.
	rc, err := h.Manager.LogsService(r.Context(), serviceID, "docker")
	if err != nil {
		h.Logger.Error("failed to get logs", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to get logs")
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	buf := make([]byte, 1024)
	for {
		// Wait for context cancellation or stream closure
		select {
		case <-r.Context().Done():
			return
		default:
			n, err := rc.Read(buf)
			if n > 0 {
				_, _ = w.Write(buf[:n])
				flusher.Flush()
			}
			if err != nil {
				if err != io.EOF {
					h.Logger.Error("log stream error", "error", err)
				}
				return
			}
			time.Sleep(100 * time.Millisecond) // prevent tight loop
		}
	}
}
