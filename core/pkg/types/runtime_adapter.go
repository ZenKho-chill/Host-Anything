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

// ServiceSpec is the normalized, runtime-agnostic description of a service.
// The core system derives it from a template and user-provided configuration,
// then passes it to a RuntimeAdapter. The adapter translates it into
// runtime-specific API calls (Docker SDK, Podman API, systemd unit, etc.).
//
// Secret values in Env must be decrypted in memory before being placed here;
// they are never stored in a ServiceSpec at rest.
type ServiceSpec struct {
	// ServiceID is the unique identifier of the service instance (UUID).
	ServiceID string

	// ServiceName is the human-readable display name.
	ServiceName string

	// Image is the container image reference (e.g. "nginx:stable-alpine").
	// Empty for host-mode services.
	Image string

	// Command overrides the container entrypoint. Nil means use the image default.
	Command []string

	// Env is the complete set of environment variables to inject into the service.
	// Secret values are decrypted in memory by the core before being placed here.
	Env map[string]string

	// Ports maps internal (container) port to external (host) port.
	// Example: {6379: 6379} exposes Redis on host port 6379.
	Ports map[int]int

	// Volumes maps volume name to host mount path.
	// Example: {"data": "/var/lib/hostanything/data/myredis/data"}
	Volumes map[string]string

	// ResourceLimits constrains the service's compute resource consumption.
	ResourceLimits ResourceLimits
}

// ResourceLimits describes compute resource constraints for a service.
// A zero value for any field means the resource is unconstrained.
type ResourceLimits struct {
	// MemoryMB is the maximum memory in megabytes. Zero means unlimited.
	MemoryMB int64

	// CPUMillicores is the CPU limit in millicores (1000m = 1 vCPU). Zero means unlimited.
	CPUMillicores int64
}

// LogOptions configures the log stream returned by [RuntimeAdapter.Logs].
type LogOptions struct {
	// Follow streams logs in real-time when true. Blocks until context is cancelled.
	Follow bool

	// TailLines is the number of recent lines to return before streaming.
	// Zero means return all historical logs.
	TailLines int
}

// RuntimeAdapter is the central interface for all runtime implementations.
// Every runtime (Docker, Podman, Kubernetes, host-mode) must satisfy this interface.
// The core system interacts exclusively through RuntimeAdapter — never with
// runtime-specific APIs directly. See ADR-003 for the rationale.
//
// All method implementations must be safe for concurrent use.
// All errors should wrap a sentinel from [github.com/host-anything/hostanything/pkg/errors]
// using fmt.Errorf("adapterName.MethodName: %w", err).
type RuntimeAdapter interface {
	// Name returns the unique, lowercase identifier of this runtime.
	// Examples: "docker", "podman", "k8s", "host".
	Name() string

	// Deploy creates and starts a new service from the given ServiceSpec.
	// Returns [errors.ErrConflict] if a service with the same ServiceID already exists.
	Deploy(ctx context.Context, spec ServiceSpec) error

	// Start resumes a previously stopped service.
	// Returns [errors.ErrNotFound] if no service with the given ID exists.
	Start(ctx context.Context, serviceID string) error

	// Stop gracefully halts a running service without removing it.
	// Returns [errors.ErrNotFound] if no service with the given ID exists.
	Stop(ctx context.Context, serviceID string) error

	// Status returns the current runtime status of a service.
	// Returns [errors.ErrNotFound] if no service with the given ID exists.
	Status(ctx context.Context, serviceID string) (ServiceStatus, error)

	// Logs returns a stream of the service's log output.
	// The caller is responsible for closing the returned [io.ReadCloser].
	Logs(ctx context.Context, serviceID string, opts LogOptions) (io.ReadCloser, error)

	// ApplyConfig updates the live configuration of a service.
	// The adapter determines whether an in-place reload or full recreation is required,
	// based on the template's update strategy (SPEC-001).
	ApplyConfig(ctx context.Context, serviceID string, spec ServiceSpec) error

	// Remove permanently deletes a service and optionally its associated volumes.
	// Returns [errors.ErrNotFound] if no service with the given ID exists.
	Remove(ctx context.Context, serviceID string, removeVolumes bool) error
}
