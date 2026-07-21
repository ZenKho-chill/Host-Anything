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
	"fmt"

	"github.com/host-anything/hostanything/pkg/errors"
	"github.com/host-anything/hostanything/pkg/types"
)

// validLogLevels is the set of accepted log level strings.
var validLogLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// Validate checks that all required fields in cfg are present and within bounds.
// It returns [errors.ErrValidation] wrapped with field-level context on failure.
//
// Validate must be called after [ApplyDefaults] to ensure optional fields have
// been populated before they are inspected.
func Validate(cfg *types.SystemConfig) error {
	if cfg.Server.APIPort < 1 || cfg.Server.APIPort > 65535 {
		return fmt.Errorf(
			"config.Validate: server.api_port %d is out of range [1, 65535]: %w",
			cfg.Server.APIPort, errors.ErrValidation,
		)
	}
	if cfg.Server.BindAddress == "" {
		return fmt.Errorf("config.Validate: server.bind_address must not be empty: %w", errors.ErrValidation)
	}
	if !validLogLevels[cfg.Server.LogLevel] {
		return fmt.Errorf(
			"config.Validate: server.log_level %q is invalid, must be one of: debug, info, warn, error: %w",
			cfg.Server.LogLevel, errors.ErrValidation,
		)
	}
	if cfg.Auth.JWTSecret == "" {
		return fmt.Errorf("config.Validate: auth.jwt_secret must not be empty: %w", errors.ErrValidation)
	}
	if cfg.Auth.AdminUsername == "" {
		return fmt.Errorf("config.Validate: auth.admin_username must not be empty: %w", errors.ErrValidation)
	}
	if cfg.Auth.AdminPassword == "" {
		return fmt.Errorf("config.Validate: auth.admin_password must not be empty: %w", errors.ErrValidation)
	}
	if cfg.Paths.DataDir == "" {
		return fmt.Errorf("config.Validate: paths.data_dir must not be empty: %w", errors.ErrValidation)
	}
	if cfg.Paths.TemplateDir == "" {
		return fmt.Errorf("config.Validate: paths.template_dir must not be empty: %w", errors.ErrValidation)
	}
	return nil
}
