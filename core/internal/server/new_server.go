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
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

// HTTP server timeout constants. These are tuned for a local management API
// that is not expected to handle long-lived streaming connections in M1.
const (
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 5 * time.Second
)

// Options holds all dependencies needed to construct the HTTP server.
// All fields are required unless otherwise noted.
type Options struct {
	// Config is the parsed, validated system configuration.
	Config *types.SystemConfig

	// Logger is the structured logger for request logging.
	Logger *slog.Logger

	// Version is the binary version string embedded at build time.
	Version string

	// EnabledRuntimes is the list of runtime adapter names currently active.
	// Used to populate the /health response.
	EnabledRuntimes []string

	// Registry is the template registry for managing service templates.
	Registry *template.Registry
}

// NewServer constructs a fully configured [net/http.Server] with all routes
// and middleware registered. It does not start listening; call ListenAndServe
// on the returned server when ready.
//
// Returns an error if required options are missing.
func NewServer(opts Options) (*http.Server, error) {
	if opts.Config == nil {
		return nil, fmt.Errorf("server.NewServer: opts.Config must not be nil")
	}
	if opts.Logger == nil {
		return nil, fmt.Errorf("server.NewServer: opts.Logger must not be nil")
	}
	if opts.Registry == nil {
		return nil, fmt.Errorf("server.NewServer: opts.Registry must not be nil")
	}

	r := chi.NewRouter()

	// Global middleware — order matters.
	r.Use(middleware.RequestID)  // assign unique request ID
	r.Use(middleware.RealIP)     // respect X-Forwarded-For behind reverse proxy
	r.Use(requestLogger(opts.Logger)) // structured JSON request logging
	r.Use(middleware.Recoverer)  // recover from panics, log, return 500

	RegisterRoutes(r, opts)

	addr := fmt.Sprintf("%s:%d", opts.Config.Server.BindAddress, opts.Config.Server.APIPort)

	return &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}, nil
}

// requestLogger returns a chi middleware that logs each completed HTTP request
// using the provided slog.Logger in JSON format.
func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("http request served",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes_written", ww.BytesWritten(),
				"request_id", middleware.GetReqID(r.Context()),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}
