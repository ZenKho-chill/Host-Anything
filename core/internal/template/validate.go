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

	"github.com/host-anything/hostanything/pkg/types"
)

// Validate ensures a parsed template strictly adheres to SPEC-001.
// It orchestrates sub-validators for each section.
func Validate(t *types.Template) error {
	if err := validateMeta(t.Meta); err != nil {
		return fmt.Errorf("invalid [meta] section: %w", err)
	}

	if err := validateConfigVars(t.Config); err != nil {
		return fmt.Errorf("invalid [[config]] section: %w", err)
	}

	if err := validateRuntime(t.Runtime); err != nil {
		return fmt.Errorf("invalid [runtime] section: %w", err)
	}

	if err := validateNetwork(t.Network); err != nil {
		return fmt.Errorf("invalid [[network]] section: %w", err)
	}

	return nil
}
