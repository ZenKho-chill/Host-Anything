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

// SystemConfig represents the root configuration structure for the
// hostanything daemon. It is parsed from hostanything.toml (SPEC-003).
//
// Unknown TOML fields are rejected at load time to enforce strict hygiene.
type SystemConfig struct {
	// Server controls HTTP server behavior.
	Server ServerConfig `toml:"server"`
	// Auth controls authentication settings.
	Auth AuthConfig `toml:"auth"`
	// Runtimes specifies which container/host runtimes are enabled.
	Runtimes RuntimesConfig `toml:"runtimes"`
	// Paths holds filesystem path overrides.
	Paths PathsConfig `toml:"paths"`
}

// ServerConfig controls the HTTP API server behavior.
type ServerConfig struct {
	// APIPort is the port the REST API listens on. Default: 8080.
	APIPort int `toml:"api_port"`
	// BindAddress is the address the server binds to. Default: "127.0.0.1".
	BindAddress string `toml:"bind_address"`
	// LogLevel controls log verbosity. One of: debug, info, warn, error.
	LogLevel string `toml:"log_level"`
}

// AuthConfig controls authentication and session settings.
type AuthConfig struct {
	// JWTSecret is the signing key for JWT tokens. Must be set; generated on install.
	JWTSecret string `toml:"jwt_secret"`
	// SessionTimeout is the duration before a JWT session expires. Default: "24h".
	SessionTimeout string `toml:"session_timeout"`
	// Fail2BanEnabled enables writing failed-login events to the auth log
	// so that fail2ban can monitor and block repeated failures.
	Fail2BanEnabled bool `toml:"fail2ban_enabled"`
}

// RuntimesConfig specifies which runtimes the daemon will use.
// At least one runtime must be enabled for the daemon to be useful.
type RuntimesConfig struct {
	// DockerEnabled enables the Docker runtime adapter.
	DockerEnabled bool `toml:"docker_enabled"`
	// PodmanEnabled enables the Podman runtime adapter.
	PodmanEnabled bool `toml:"podman_enabled"`
	// K8sEnabled enables the Kubernetes runtime adapter.
	K8sEnabled bool `toml:"k8s_enabled"`
	// HostEnabled enables the host (bare-metal process) runtime adapter.
	HostEnabled bool `toml:"host_enabled"`
}

// PathsConfig holds filesystem paths used by the daemon.
type PathsConfig struct {
	// DataDir is the base directory for per-service state and config files.
	// Default: /var/lib/hostanything/data
	DataDir string `toml:"data_dir"`
	// TemplateDir is the directory where installed service templates are stored.
	// Default: /var/lib/hostanything/templates
	TemplateDir string `toml:"template_dir"`
}
