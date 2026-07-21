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

// Package crypto provides AES-256-GCM encryption utilities for the
// hostanything daemon. It encrypts "secret" type config variables
// before they are written to disk, per the secret handling policy in SPEC-003.
//
// Encrypted values use a versioned format:
//
//	enc:v1:aes256gcm:<base64(nonce || ciphertext)>
//
// The version prefix ("v1") allows future algorithm migration without
// invalidating stored ciphertext. Use [IsEncrypted] to detect the prefix.
//
// Key management: a 32-byte master key is generated once on first run and
// stored at {data_dir}/master.key (mode 0600). Use [LoadOrCreateKey] in
// the application startup sequence.
package crypto
