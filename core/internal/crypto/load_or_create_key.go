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

package crypto

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadOrCreateKey loads the master encryption key from the given file path.
// If the file does not exist, a new 32-byte key is generated, written to the
// path with permissions 0600, and returned.
//
// The containing directory is created (mode 0700) if it does not exist.
//
// WARNING: The key file must never be committed to version control or backed
// up in the same location as the encrypted secrets it protects.
func LoadOrCreateKey(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		if len(data) != KeySize {
			return nil, fmt.Errorf(
				"crypto.LoadOrCreateKey: key file %q has unexpected length %d (want %d)",
				path, len(data), KeySize,
			)
		}
		return data, nil
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("crypto.LoadOrCreateKey: read key file: %w", err)
	}

	// File doesn't exist — generate and persist a fresh key.
	key, err := GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("crypto.LoadOrCreateKey: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("crypto.LoadOrCreateKey: create key directory %q: %w", dir, err)
	}

	if err := os.WriteFile(path, key, 0o600); err != nil {
		return nil, fmt.Errorf("crypto.LoadOrCreateKey: write key file %q: %w", path, err)
	}

	return key, nil
}
