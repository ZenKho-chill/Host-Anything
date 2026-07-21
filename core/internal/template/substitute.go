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
	"os"

	"github.com/host-anything/hostanything/pkg/types"
)

// Substitute interpolates variables in the format ${VAR} into the given template's
// Command array, using the provided map of resolved variables.
// Unknown variables are left unchanged (or expanded to empty string based on os.Expand behavior).
func Substitute(t *types.Template, resolved map[string]string) []string {
	if len(t.Runtime.Command) == 0 {
		return nil
	}

	mapper := func(varName string) string {
		if val, ok := resolved[varName]; ok {
			return val
		}
		return ""
	}

	result := make([]string, len(t.Runtime.Command))
	for i, arg := range t.Runtime.Command {
		result[i] = os.Expand(arg, mapper)
	}

	return result
}
