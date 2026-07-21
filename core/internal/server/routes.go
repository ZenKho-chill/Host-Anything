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
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/api"
)

// RegisterRoutes sets up all API endpoints on the given router.
func RegisterRoutes(r chi.Router, opts Options) {
	// Public Routes
	r.Get("/api/v1/health", api.HealthHandler(opts.Version, opts.EnabledRuntimes))
	r.Post("/api/v1/auth/login", api.AuthHandler(opts.DB, opts.Config, opts.Logger))

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware(opts.Config.Auth.JWTSecret, opts.DB))

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

		// Marketplace
		mktHandler := api.NewMarketplaceHandler(opts.Config.Paths.TemplateDir, opts.Logger)
		r.Get("/api/v1/marketplace/search", mktHandler.SearchTemplates)
		r.Get("/api/v1/marketplace/preview/{owner}/{repo}", mktHandler.PreviewTemplate)
		r.Post("/api/v1/marketplace/install", mktHandler.InstallTemplate)

		// Enterprise Core APIs
		// Users
		userHandler := api.NewUserHandler(opts.DB, opts.Logger)
		r.Get("/api/v1/users", userHandler.ListUsers)
		r.Post("/api/v1/users", userHandler.CreateUser)

		// Roles
		roleHandler := api.NewRoleHandler(opts.DB, opts.Logger)
		r.Get("/api/v1/roles", roleHandler.ListRoles)

		// Schedules
		scheduleHandler := api.NewScheduleHandler(opts.DB, opts.Logger)
		r.Get("/api/v1/schedules", scheduleHandler.ListSchedules)
		r.Post("/api/v1/schedules", scheduleHandler.CreateSchedule)
		r.Delete("/api/v1/schedules/{id}", scheduleHandler.DeleteSchedule)

		// Settings
		settingsHandler := api.NewSettingsHandler(opts.DB, opts.Logger)
		r.Get("/api/v1/settings", settingsHandler.GetSettings)
		r.Post("/api/v1/settings", settingsHandler.UpdateSettings)

		// File Manager
		fileHandler := api.NewFileHandler(opts.Config.Paths.DataDir, opts.Logger)
		r.Get("/api/v1/files/*", fileHandler.ListFiles)
	})

	// Serve static frontend files from web/dist (if present)
	workDir, _ := os.Getwd()
	// Try to find the dist directory relative to the current working directory
	// In production (when distributed as a .deb), files will likely be in /usr/share/hostanything/web
	// For local dev, they are in ../web/dist or ./web/dist
	distDirs := []string{
		filepath.Join(workDir, "web", "dist"),
		filepath.Join(workDir, "..", "web", "dist"),
		"/usr/share/hostanything/web",
	}

	var staticDir string
	for _, dir := range distDirs {
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			staticDir = dir
			break
		}
	}

	if staticDir != "" {
		fs := http.FileServer(http.Dir(staticDir))

		// Fallback route for SPA (React Router)
		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			path := filepath.Join(staticDir, req.URL.Path)
			if stat, err := os.Stat(path); os.IsNotExist(err) || stat.IsDir() {
				// If file doesn't exist or is a directory, serve index.html for SPA routing
				http.ServeFile(w, req, filepath.Join(staticDir, "index.html"))
				return
			}
			fs.ServeHTTP(w, req)
		})
	}
}
