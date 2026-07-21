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

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/host-anything/hostanything/pkg/types"
)

// ApplyDefaults fills in zero-value fields in cfg with sensible defaults.
// It must be called before [Validate] so that optional fields are populated
// before they are checked for correctness.
//
// ApplyDefaults never overrides a value that has already been set by the caller.
func ApplyDefaults(cfg *types.SystemConfig) {
	if cfg.Server.APIPort == 0 {
		cfg.Server.APIPort = DefaultAPIPort
	}
	if cfg.Server.BindAddress == "" {
		cfg.Server.BindAddress = DefaultBindAddress
	}
	if cfg.Server.LogLevel == "" {
		cfg.Server.LogLevel = DefaultLogLevel
	}
	if cfg.Auth.SessionTimeout == "" {
		cfg.Auth.SessionTimeout = DefaultSessionTimeout
	}
	if cfg.Auth.AdminUsername == "" {
		cfg.Auth.AdminUsername = "admin"
	}
	if cfg.Auth.AdminPassword == "" {
		cfg.Auth.AdminPassword = "admin" // In a real app, this should force a change on first boot
	}
	if cfg.Auth.JWTSecret == "" {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		cfg.Auth.JWTSecret = hex.EncodeToString(b)
	}

	if cfg.Paths.DataDir == "" {
		cfg.Paths.DataDir = DefaultDataDir
	}
	if cfg.Paths.TemplateDir == "" {
		cfg.Paths.TemplateDir = DefaultTemplateDir
	}
}
