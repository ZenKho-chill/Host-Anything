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

// RegisterRoutes registers all API routes onto the provided chi router.
// This function is kept separate from server construction to make the
// routing table easy to read and extend without touching server setup logic.
//
// Route layout follows SPEC-004. Authentication middleware will be added
// to protected route groups in M4 (SPEC-030).
func RegisterRoutes(r chi.Router, opts Options) {
	r.Route("/api/v1", func(r chi.Router) {
		// GET /api/v1/health — unauthenticated liveness probe (SPEC-004 §1).
		r.Get("/health", api.HealthHandler(opts.Version, opts.EnabledRuntimes))

		// TODO(M3): service lifecycle endpoints (SPEC-004 §2–7)
		// TODO(M4): template + marketplace endpoints (SPEC-004 §8–10)
		// TODO(M4): auth endpoints POST /auth/login (SPEC-030)
	})
}
