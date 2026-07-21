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

package types

import (
	"context"
	"io"
)

// ServiceSpec bundles the template definition and the user's specific deployment variables.
type ServiceSpec struct {
	// ServiceID is the unique identifier for this deployment (e.g. uuid).
	ServiceID string
	// Template is the parsed service definition.
	Template *Template
	// ResolvedEnv maps the final environment variables to inject (including decrypted secrets).
	ResolvedEnv map[string]string
}

// RuntimeAdapter abstracts the underlying execution environment (Docker, Podman, Host, etc.)
// from the core Host Anything lifecycle manager.
type RuntimeAdapter interface {
	// Deploy creates and starts a new service based on the spec.
	// If it already exists, it may return an error or recreate it.
	Deploy(ctx context.Context, spec *ServiceSpec) error

	// Stop gracefully terminates the service without removing its data or definition.
	Stop(ctx context.Context, serviceID string) error

	// Start resumes a previously stopped service.
	Start(ctx context.Context, serviceID string) error

	// Remove permanently deletes the service and optionally its volumes.
	Remove(ctx context.Context, serviceID string) error

	// Status fetches real-time health and port mapping information.
	Status(ctx context.Context, serviceID string) (*ServiceStatus, error)

	// Logs retrieves the standard output and error streams for the service.
	// The caller is responsible for closing the reader.
	Logs(ctx context.Context, serviceID string) (io.ReadCloser, error)
}
