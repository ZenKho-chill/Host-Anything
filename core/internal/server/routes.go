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

package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/api"
)

// RegisterRoutes sets up all HTTP endpoints for the application.
func RegisterRoutes(r chi.Router, opts Options) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", api.HealthHandler(opts.Version, opts.EnabledRuntimes))

		r.Route("/templates", func(r chi.Router) {
			r.Get("/", api.TemplateListHandler(opts.Registry, opts.Logger))
			r.Get("/{name}", api.TemplateGetHandler(opts.Registry, opts.Logger))
		})
	})
}
