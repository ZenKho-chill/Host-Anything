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
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

// Decrypt decrypts a value produced by [Encrypt].
// The value must begin with the "enc:v1:aes256gcm:" prefix.
// Returns an error if the format is unexpected or authentication fails.
func Decrypt(encrypted string, key []byte) (string, error) {
	if len(key) != KeySize {
		return "", fmt.Errorf("crypto.Decrypt: key must be %d bytes, got %d", KeySize, len(key))
	}

	if !strings.HasPrefix(encrypted, encryptedPrefix) {
		return "", fmt.Errorf("crypto.Decrypt: value does not have expected prefix %q", encryptedPrefix)
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encrypted, encryptedPrefix))
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: base64 decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: create GCM wrapper: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(payload) < nonceSize {
		return "", fmt.Errorf("crypto.Decrypt: payload is too short to contain a nonce")
	}

	nonce, ciphertext := payload[:nonceSize], payload[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: authentication failed (wrong key or corrupted data): %w", err)
	}

	return string(plaintext), nil
}
