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

package config

// Default configuration values used when fields are absent from the config file.
// These are exported so that tests and callers can reference the canonical defaults.
const (
	// DefaultConfigPath is the system-wide config file location on Debian.
	DefaultConfigPath = "/etc/hostanything/hostanything.toml"

	// DefaultAPIPort is the default TCP port the REST API listens on.
	DefaultAPIPort = 8080

	// DefaultBindAddress is the default address the server binds to.
	// Defaults to loopback for security — expose via reverse proxy if needed.
	DefaultBindAddress = "127.0.0.1"

	// DefaultLogLevel is the default logging verbosity.
	DefaultLogLevel = "info"

	// DefaultSessionTimeout is the default JWT token lifetime.
	DefaultSessionTimeout = "24h"

	// DefaultDataDir is the default base directory for per-service data.
	DefaultDataDir = "/var/lib/hostanything/data"

	// DefaultTemplateDir is the default directory for installed service templates.
	DefaultTemplateDir = "/var/lib/hostanything/templates"
)
