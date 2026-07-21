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
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/marketplace"
)

// MarketplaceHandler handles all /api/v1/marketplace endpoints.
// It is a struct to allow injecting the marketplace client and installer.
type MarketplaceHandler struct {
	client    *marketplace.Client
	installer *marketplace.Installer
	logger    *slog.Logger
}

// NewMarketplaceHandler constructs a MarketplaceHandler with the given dependencies.
func NewMarketplaceHandler(templateDir string, logger *slog.Logger) *MarketplaceHandler {
	return &MarketplaceHandler{
		client:    marketplace.NewClient(),
		installer: marketplace.NewInstaller(templateDir),
		logger:    logger,
	}
}

// SearchTemplates handles GET /api/v1/marketplace/search?q=...
// It queries the GitHub Search API for templates matching the query.
func (h *MarketplaceHandler) SearchTemplates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	results, err := h.client.SearchTemplates(r.Context(), query)
	if err != nil {
		h.logger.Error("marketplace search failed", "query", query, "error", err)
		WriteError(w, http.StatusBadGateway, "failed to search GitHub marketplace")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"results": results,
		"count":   len(results),
	})
}

// PreviewTemplate handles GET /api/v1/marketplace/preview/{owner}/{repo}
// It fetches and parses the template.toml from GitHub without installing it.
func (h *MarketplaceHandler) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	tmpl, err := h.client.FetchTemplate(r.Context(), owner, repo)
	if err != nil {
		h.logger.Error("marketplace preview failed", "owner", owner, "repo", repo, "error", err)
		WriteError(w, http.StatusBadGateway, "failed to fetch template from GitHub")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"template":    tmpl,
		"is_official": owner == marketplace.OfficialOrg,
	})
}

// installRequest is the JSON body for the install endpoint.
type installRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

// InstallTemplate handles POST /api/v1/marketplace/install
// It downloads, validates, and registers the template locally.
func (h *MarketplaceHandler) InstallTemplate(w http.ResponseWriter, r *http.Request) {
	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Owner == "" || req.Repo == "" {
		WriteError(w, http.StatusBadRequest, "owner and repo are required")
		return
	}

	tmpl, err := h.installer.Install(r.Context(), req.Owner, req.Repo)
	if err != nil {
		h.logger.Error("marketplace install failed",
			"owner", req.Owner,
			"repo", req.Repo,
			"error", err,
		)
		WriteError(w, http.StatusBadGateway, "failed to install template: "+err.Error())
		return
	}

	h.logger.Info("template installed from marketplace",
		"name", tmpl.Meta.Name,
		"version", tmpl.Meta.Version,
		"owner", req.Owner,
	)

	WriteJSON(w, http.StatusOK, map[string]any{
		"message": "template installed successfully",
		"name":    tmpl.Meta.Name,
		"version": tmpl.Meta.Version,
	})
}
