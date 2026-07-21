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

var nameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func validateMeta(meta types.TemplateMeta) error {
	if meta.Name == "" {
		return fmt.Errorf("missing required field 'name'")
	}
	if !nameRegex.MatchString(meta.Name) {
		return fmt.Errorf("name %q must contain only lowercase letters, numbers, and hyphens", meta.Name)
	}
	if meta.Version == "" {
		return fmt.Errorf("missing required field 'version'")
	}
	if meta.Description == "" {
		return fmt.Errorf("missing required field 'description'")
	}
	if meta.Author == "" {
		return fmt.Errorf("missing required field 'author'")
	}
	if meta.License == "" {
		return fmt.Errorf("missing required field 'license'")
	}
	return nil
}
