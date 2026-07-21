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

package config_test

import (
	"testing"

	"github.com/host-anything/hostanything/internal/config"
	"github.com/host-anything/hostanything/pkg/errors"
	"github.com/host-anything/hostanything/pkg/types"
)

// validBaseConfig returns a fully-populated SystemConfig that passes Validate.
func validBaseConfig() *types.SystemConfig {
	return &types.SystemConfig{
		Server: types.ServerConfig{
			APIPort:     8080,
			BindAddress: "127.0.0.1",
			LogLevel:    "info",
		},
		Auth: types.AuthConfig{
			AdminUsername:  "admin",
			AdminPassword:  "admin",
			JWTSecret:      "supersecretkey",
			SessionTimeout: "24h",
		},
		Paths: types.PathsConfig{
			DataDir:     "/var/lib/hostanything/data",
			TemplateDir: "/var/lib/hostanything/templates",
		},
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	if err := config.Validate(validBaseConfig()); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}

func TestValidate_AllLogLevels_AreAccepted(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Server.LogLevel = level
			if err := config.Validate(cfg); err != nil {
				t.Errorf("expected no error for log_level=%q, got: %v", level, err)
			}
		})
	}
}

func TestValidate_InvalidPort_Zero(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Server.APIPort = 0
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for port 0, got: %v", err)
	}
}

func TestValidate_InvalidPort_Negative(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Server.APIPort = -1
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for negative port, got: %v", err)
	}
}

func TestValidate_InvalidPort_TooHigh(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Server.APIPort = 99999
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for port 99999, got: %v", err)
	}
}

func TestValidate_ValidPort_Boundaries(t *testing.T) {
	for _, port := range []int{1, 80, 443, 8080, 65535} {
		t.Run("port", func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Server.APIPort = port
			if err := config.Validate(cfg); err != nil {
				t.Errorf("expected no error for port %d, got: %v", port, err)
			}
		})
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Server.LogLevel = "verbose"
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for log_level='verbose', got: %v", err)
	}
}

func TestValidate_EmptyBindAddress(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Server.BindAddress = ""
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for empty bind_address, got: %v", err)
	}
}

func TestValidate_EmptyJWTSecret(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Auth.JWTSecret = ""
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for empty jwt_secret, got: %v", err)
	}
}

func TestValidate_EmptyDataDir(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Paths.DataDir = ""
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for empty data_dir, got: %v", err)
	}
}

func TestValidate_EmptyTemplateDir(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Paths.TemplateDir = ""
	if err := config.Validate(cfg); !errors.Is(err, errors.ErrValidation) {
		t.Errorf("expected ErrValidation for empty template_dir, got: %v", err)
	}
}
