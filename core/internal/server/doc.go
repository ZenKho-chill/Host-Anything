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

// Package server wires together the HTTP router, middleware stack, and API handlers
// into a runnable [net/http.Server].
//
// Responsibilities:
//   - Server construction with timeouts ([NewServer])
//   - Route registration ([RegisterRoutes])
//   - Structured request logging middleware
//
// No business logic lives here. All handler logic is delegated to [internal/api].
// Authentication middleware (M4) will be registered here per SPEC-030.
package server
