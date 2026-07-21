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
	"os"

	"github.com/host-anything/hostanything/pkg/types"
)

// ParseTOML reads a template TOML file from disk and parses it.
// See [ParseBytes] for details on validation.
func ParseTOML(path string) (*types.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template.ParseTOML: read file %q: %w", path, err)
	}
	t, err := ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("template.ParseTOML: file %q: %w", path, err)
	}
	return t, nil
}
