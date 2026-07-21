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

import "strings"

// IsEncrypted reports whether value is an encrypted string produced by [Encrypt].
// It detects the versioned "enc:v1:aes256gcm:" prefix without attempting decryption.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}
