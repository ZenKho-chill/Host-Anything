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

	"github.com/host-anything/hostanything/internal/crypto"
	"github.com/host-anything/hostanything/pkg/types"
)

// Resolve merges user variables with template defaults. It also validates
// the user input and encrypts any secrets if an encryption key is provided.
// If cryptoKey is nil, secrets are left in plain text (e.g. for runtime injection).
// If cryptoKey is provided, secrets are encrypted (e.g. for saving to store).
func Resolve(t *types.Template, userVars map[string]string, cryptoKey []byte) (map[string]string, error) {
	if err := ValidateUserInput(t, userVars); err != nil {
		return nil, fmt.Errorf("template.Resolve: %w", err)
	}

	resolved := make(map[string]string)

	for _, v := range t.Config {
		val, provided := userVars[v.Name]

		if !provided {
			if v.Default != nil {
				val = fmt.Sprintf("%v", v.Default)
			} else {
				val = ""
			}
		}

		if val != "" && v.Type == ConfigTypeSecret && cryptoKey != nil {
			if !crypto.IsEncrypted(val) {
				encrypted, err := crypto.Encrypt(val, cryptoKey)
				if err != nil {
					return nil, fmt.Errorf("template.Resolve: encrypt %q: %w", v.Name, err)
				}
				val = encrypted
			}
		}

		resolved[v.Name] = val
	}

	return resolved, nil
}
