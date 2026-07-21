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
	"strconv"
	"strings"

	"github.com/host-anything/hostanything/pkg/types"
)

// ValidateUserInput ensures that a map of user-provided values satisfies
// the constraints (required, type, regex, options) defined in the template.
func ValidateUserInput(t *types.Template, userVars map[string]string) error {
	for _, v := range t.Config {
		val, provided := userVars[v.Name]

		if !provided {
			if v.Required && v.Default == nil {
				return fmt.Errorf("missing required variable %q", v.Name)
			}
			continue // Valid, will use default later in Resolve
		}

		// Type validation
		switch v.Type {
		case ConfigTypeInt:
			if _, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("variable %q must be an integer, got %q", v.Name, val)
			}
		case ConfigTypeBoolean:
			lower := strings.ToLower(val)
			if lower != "true" && lower != "false" {
				return fmt.Errorf("variable %q must be a boolean ('true' or 'false'), got %q", v.Name, val)
			}
		case ConfigTypeEnum:
			valid := false
			for _, opt := range v.Options {
				if val == opt {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("variable %q must be one of %v, got %q", v.Name, v.Options, val)
			}
		}

		// Regex validation
		if v.ValidationRegex != "" {
			matched, err := regexp.MatchString(v.ValidationRegex, val)
			if err != nil {
				// Should have been caught by validateConfigVars, but handle just in case.
				return fmt.Errorf("invalid regex for variable %q: %w", v.Name, err)
			}
			if !matched {
				return fmt.Errorf("variable %q value does not match regex %q", v.Name, v.ValidationRegex)
			}
		}
	}

	// Optional: Check for unknown variables provided by user that aren't in template
	for k := range userVars {
		found := false
		for _, v := range t.Config {
			if v.Name == k {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown variable %q provided", k)
		}
	}

	return nil
}
