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

func validateRuntime(rt types.RuntimeConfig) error {
	if len(rt.Supported) == 0 {
		return fmt.Errorf("missing required field 'supported'")
	}
	if rt.Image == "" {
		return fmt.Errorf("missing required field 'image'")
	}

	if rt.Preferred != "" {
		supported := false
		for _, s := range rt.Supported {
			if s == rt.Preferred {
				supported = true
				break
			}
		}
		if !supported {
			return fmt.Errorf("preferred runtime %q is not in the supported list %v", rt.Preferred, rt.Supported)
		}
	}

	return nil
}
