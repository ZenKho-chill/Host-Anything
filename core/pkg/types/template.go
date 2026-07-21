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

// Template represents a fully parsed Host Anything service template.
// Templates are described in TOML format per SPEC-001 and are stored
// in the local registry under {template_dir}/{name}/{version}/template.toml.
type Template struct {
	// Meta holds descriptive metadata about the template.
	Meta TemplateMeta `toml:"meta" json:"meta"`
	// Requirements describes minimum system resources needed to run the service.
	Requirements TemplateRequirements `toml:"requirements" json:"requirements"`
	// Config defines user-configurable environment variables (zero or more).
	Config []ConfigVar `toml:"config" json:"config"`
	// Runtime specifies the execution environment.
	Runtime RuntimeConfig `toml:"runtime" json:"runtime"`
	// Volumes maps persistent data directories.
	Volumes []VolumeConfig `toml:"volumes" json:"volumes"`
	// Network lists ports exposed by the service.
	Network []NetworkConfig `toml:"network" json:"network"`
	// Healthcheck defines how to probe service health.
	Healthcheck HealthcheckConfig `toml:"healthcheck" json:"healthcheck"`
	// Update defines how configuration changes are applied.
	Update UpdateConfig `toml:"update" json:"update"`
}

// TemplateMeta holds descriptive metadata about a service template.
// All fields marked required must be non-empty per SPEC-001.
type TemplateMeta struct {
	// Name is the unique slug identifier (e.g. "redis", "nginx").
	Name string `toml:"name" json:"name"`
	// Version is a SemVer string (e.g. "1.0.0").
	Version string `toml:"version" json:"version"`
	// Description is a short, human-readable explanation of the service.
	Description string `toml:"description" json:"description"`
	// Author is the maintainer name or handle.
	Author string `toml:"author" json:"author"`
	// License is an SPDX license identifier (e.g. "MIT", "Apache-2.0").
	License string `toml:"license" json:"license"`
	// Tags are optional marketplace categorization labels.
	Tags []string `toml:"tags" json:"tags,omitempty"`
	// Homepage is the upstream project URL (optional).
	Homepage string `toml:"homepage" json:"homepage,omitempty"`
}

// TemplateRequirements describes minimum system resources the host must provide.
// All fields are optional.
type TemplateRequirements struct {
	// MinMemory is the minimum RAM required (e.g. "512MB", "1GB").
	MinMemory string `toml:"min_memory" json:"min_memory,omitempty"`
	// MinCPU is the minimum CPU cores required (e.g. 0.5 = half a core).
	MinCPU float64 `toml:"min_cpu" json:"min_cpu,omitempty"`
	// Ports lists host ports that must be available before deployment.
	Ports []int `toml:"ports" json:"ports,omitempty"`
	// DiskSpace is the minimum disk space required (e.g. "10GB").
	DiskSpace string `toml:"disk_space" json:"disk_space,omitempty"`
}

// ConfigVar defines a single user-configurable environment variable.
// The Default field accepts string, int, or boolean values as TOML supports.
// All values are converted to strings when injected into the runtime environment.
type ConfigVar struct {
	// Name is the environment variable name (e.g. "DB_PASSWORD").
	Name string `toml:"name" json:"name"`
	// Type is one of: string, int, boolean, secret, enum.
	Type string `toml:"type" json:"type"`
	// Default is the fallback value when the user does not provide one.
	// Use interface{} to support string, int, and boolean TOML values.
	Default interface{} `toml:"default" json:"default,omitempty"`
	// Required indicates the user must supply a value (unless Default is set).
	Required bool `toml:"required" json:"required"`
	// Description is UI help text explaining what this variable does.
	Description string `toml:"description" json:"description,omitempty"`
	// ValidationRegex is a Go regexp pattern the value must match (optional).
	ValidationRegex string `toml:"validation_regex" json:"validation_regex,omitempty"`
	// Options lists valid values when Type == "enum".
	Options []string `toml:"options" json:"options,omitempty"`
}

// RuntimeConfig specifies how and where the service runs.
type RuntimeConfig struct {
	// Preferred is the suggested runtime (e.g. "docker"). Optional.
	Preferred string `toml:"preferred" json:"preferred,omitempty"`
	// Supported is the list of runtimes that can run this service. Required.
	Supported []string `toml:"supported" json:"supported"`
	// Image is the container image reference (e.g. "redis:7-alpine"). Required.
	Image string `toml:"image" json:"image"`
	// Command overrides the container entrypoint. Optional.
	Command []string `toml:"command" json:"command,omitempty"`
}

// VolumeConfig defines a persistent storage mount for the service.
type VolumeConfig struct {
	// Name is the logical volume identifier.
	Name string `toml:"name" json:"name"`
	// MountPath is the path inside the container/service where data is stored.
	MountPath string `toml:"mount_path" json:"mount_path"`
	// Description explains what this volume stores.
	Description string `toml:"description" json:"description,omitempty"`
}

// NetworkConfig defines a network port exposed by the service.
type NetworkConfig struct {
	// InternalPort is the port the application listens on inside the container.
	InternalPort int `toml:"internal_port" json:"internal_port"`
	// ExternalPort is the suggested host-side port. Zero means no suggestion.
	ExternalPort int `toml:"external_port" json:"external_port,omitempty"`
	// Protocol is one of: tcp, udp, http. Defaults to "tcp" if empty.
	Protocol string `toml:"protocol" json:"protocol,omitempty"`
}

// HealthcheckConfig defines how the system determines if the service is healthy.
type HealthcheckConfig struct {
	// Command is the shell command to run (e.g. "curl -f http://localhost/ || exit 1").
	Command string `toml:"command" json:"command,omitempty"`
	// Interval between checks (e.g. "30s"). Defaults to "30s".
	Interval string `toml:"interval" json:"interval,omitempty"`
	// Timeout after which a check is considered failed (e.g. "10s").
	Timeout string `toml:"timeout" json:"timeout,omitempty"`
	// Retries before the service is considered unhealthy.
	Retries int `toml:"retries" json:"retries,omitempty"`
}

// UpdateConfig defines how the service handles configuration changes.
type UpdateConfig struct {
	// Strategy is one of: "recreate" (stop + start) or "rolling".
	Strategy string `toml:"strategy" json:"strategy,omitempty"`
}

// TemplateSummary is a lightweight representation of a template
// used for listing without loading full details.
type TemplateSummary struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}
