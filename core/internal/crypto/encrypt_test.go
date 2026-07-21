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

package crypto_test

import (
	"strings"
	"testing"

	"github.com/host-anything/hostanything/internal/crypto"
)

func TestEncrypt_ProducesPrefix(t *testing.T) {
	key, _ := crypto.GenerateKey()
	ct, err := crypto.Encrypt("hello", key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if !strings.HasPrefix(ct, "enc:v1:aes256gcm:") {
		t.Errorf("expected enc:v1:aes256gcm: prefix, got %q", ct)
	}
}

func TestEncrypt_NonDeterministic(t *testing.T) {
	key, _ := crypto.GenerateKey()
	ct1, _ := crypto.Encrypt("same input", key)
	ct2, _ := crypto.Encrypt("same input", key)
	if ct1 == ct2 {
		t.Error("expected different ciphertexts for same plaintext (random nonce)")
	}
}

func TestDecrypt_RoundTrip(t *testing.T) {
	key, _ := crypto.GenerateKey()
	plaintext := "super-secret-password-123!"

	ct, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	got, err := crypto.Decrypt(ct, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plaintext {
		t.Errorf("expected %q, got %q", plaintext, got)
	}
}

func TestDecrypt_WrongKey_Fails(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()

	ct, _ := crypto.Encrypt("secret", key1)
	_, err := crypto.Decrypt(ct, key2)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestDecrypt_BadPrefix_Fails(t *testing.T) {
	key, _ := crypto.GenerateKey()
	_, err := crypto.Decrypt("plaintext-value", key)
	if err == nil {
		t.Error("expected error for value without enc: prefix")
	}
}

func TestDecrypt_CorruptedCiphertext_Fails(t *testing.T) {
	key, _ := crypto.GenerateKey()
	ct, _ := crypto.Encrypt("secret", key)
	corrupted := ct[:len(ct)-4] + "XXXX"
	_, err := crypto.Decrypt(corrupted, key)
	if err == nil {
		t.Error("expected error for corrupted ciphertext")
	}
}

func TestIsEncrypted_True(t *testing.T) {
	key, _ := crypto.GenerateKey()
	ct, _ := crypto.Encrypt("value", key)
	if !crypto.IsEncrypted(ct) {
		t.Error("expected IsEncrypted=true for encrypted value")
	}
}

func TestIsEncrypted_False_Plaintext(t *testing.T) {
	if crypto.IsEncrypted("plaintext") {
		t.Error("expected IsEncrypted=false for plaintext")
	}
}

func TestIsEncrypted_False_Empty(t *testing.T) {
	if crypto.IsEncrypted("") {
		t.Error("expected IsEncrypted=false for empty string")
	}
}

func TestEncrypt_WrongKeySize_Fails(t *testing.T) {
	_, err := crypto.Encrypt("hello", []byte("tooshort"))
	if err == nil {
		t.Error("expected error for wrong key size")
	}
}

func TestDecrypt_WrongKeySize_Fails(t *testing.T) {
	_, err := crypto.Decrypt("enc:v1:aes256gcm:abc", []byte("tooshort"))
	if err == nil {
		t.Error("expected error for wrong key size")
	}
}
