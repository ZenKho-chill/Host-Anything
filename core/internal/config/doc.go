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

// Package config handles loading, validation, and default-application of the
// hostanything system configuration (hostanything.toml).
//
// Configuration is defined in SPEC-003. The system config uses TOML format
// and is typically located at /etc/hostanything/hostanything.toml on Debian.
//
// Usage:
//
//	cfg, err := config.Load("/etc/hostanything/hostanything.toml")
//	if err != nil {
//	    log.Fatalf("failed to load config: %v", err)
//	}
package config
