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

package template

import (
	"fmt"
	"regexp"

	"github.com/host-anything/hostanything/pkg/types"
)

func validateConfigVars(vars []types.ConfigVar) error {
	seenNames := make(map[string]bool)

	for i, v := range vars {
		if v.Name == "" {
			return fmt.Errorf("variable at index %d missing required field 'name'", i)
		}

		if seenNames[v.Name] {
			return fmt.Errorf("duplicate variable name %q", v.Name)
		}
		seenNames[v.Name] = true

		if v.Type == "" {
			return fmt.Errorf("variable %q missing required field 'type'", v.Name)
		}

		switch v.Type {
		case ConfigTypeString, ConfigTypeInt, ConfigTypeBoolean, ConfigTypeSecret, ConfigTypeEnum:
			// valid
		default:
			return fmt.Errorf("variable %q has invalid type %q", v.Name, v.Type)
		}

		if v.Type == ConfigTypeEnum && len(v.Options) == 0 {
			return fmt.Errorf("variable %q of type enum must provide 'options'", v.Name)
		}

		if v.ValidationRegex != "" {
			if _, err := regexp.Compile(v.ValidationRegex); err != nil {
				return fmt.Errorf("variable %q has invalid validation_regex: %w", v.Name, err)
			}
		}
	}

	return nil
}
