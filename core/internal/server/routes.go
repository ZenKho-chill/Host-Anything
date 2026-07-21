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

// RegisterRoutes sets up all API endpoints on the given router.
func RegisterRoutes(r chi.Router, opts Options) {
	// Public Routes
	r.Get("/api/v1/health", api.HealthHandler(opts.Version, opts.EnabledRuntimes))
	r.Post("/api/v1/auth/login", api.AuthHandler(opts.Config, opts.Logger))

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware(opts.Config.Auth.JWTSecret))

		// Templates
		r.Get("/api/v1/templates", api.TemplateListHandler(opts.Registry, opts.Logger))
		r.Get("/api/v1/templates/{name}", api.TemplateGetHandler(opts.Registry, opts.Logger))

		// Services
		svcHandler := &api.ServiceHandler{
			Manager:  opts.Manager,
			Registry: opts.Registry,
			Logger:   opts.Logger,
			Key:      opts.MasterKey,
		}
		r.Get("/api/v1/services", svcHandler.ListServices)
		r.Post("/api/v1/services", svcHandler.DeployService)
		r.Get("/api/v1/services/{id}/logs", svcHandler.LogsService)
	})
}
