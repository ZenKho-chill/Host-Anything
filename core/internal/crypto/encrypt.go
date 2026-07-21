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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryptedPrefix is the versioned prefix that marks an encrypted value.
// All values produced by Encrypt start with this string.
const encryptedPrefix = "enc:v1:aes256gcm:"

// Encrypt encrypts plaintext using AES-256-GCM with the given 32-byte key.
// A fresh random 12-byte nonce is generated for every call, making
// the output non-deterministic even for identical inputs.
//
// Output format: enc:v1:aes256gcm:<base64(nonce || ciphertext || tag)>
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != KeySize {
		return "", fmt.Errorf("crypto.Encrypt: key must be %d bytes, got %d", KeySize, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto.Encrypt: create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto.Encrypt: create GCM wrapper: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto.Encrypt: generate nonce: %w", err)
	}

	// Seal appends ciphertext+tag to nonce, giving us nonce||ciphertext||tag.
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}
