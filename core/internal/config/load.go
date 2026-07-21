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
	"os"

	"github.com/BurntSushi/toml"
	"github.com/host-anything/hostanything/pkg/types"
)

// Load reads, parses, and validates the system configuration from the TOML file
// at the given path. It applies defaults for optional fields before validating.
//
// Unknown TOML fields are rejected to enforce strict configuration hygiene
// (per project rule #10). A misconfigured file on startup is a fatal condition.
//
// All errors are wrapped with context: "config.Load: ...".
func Load(path string) (*types.SystemConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If file is missing, proceed with empty config to apply defaults
			data = []byte{}
		} else {
			return nil, fmt.Errorf("config.Load: read file %q: %w", path, err)
		}
	}

	var cfg types.SystemConfig

	if len(data) > 0 {
		md, err := toml.Decode(string(data), &cfg)
		if err != nil {
			return nil, fmt.Errorf("config.Load: parse toml: %w", err)
		}

		// Reject unknown keys — configuration must be explicit and typo-free.
		if undecoded := md.Undecoded(); len(undecoded) > 0 {
			return nil, fmt.Errorf("config.Load: unknown fields in config file (check for typos): %v", undecoded)
		}
	}

	ApplyDefaults(&cfg)

	if err := Validate(&cfg); err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}

	return &cfg, nil
}
