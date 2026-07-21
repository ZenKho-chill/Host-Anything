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
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/host-anything/hostanything/internal/crypto"
)

func TestLoadOrCreateKey_CreatesNewKey_IfNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "master.key")

	key1, err := crypto.LoadOrCreateKey(keyPath)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if len(key1) != crypto.KeySize {
		t.Errorf("expected key size %d, got %d", crypto.KeySize, len(key1))
	}

	// Verify file permissions (0600 expected)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(keyPath)
		if err != nil {
			t.Fatalf("stat key file: %v", err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Errorf("expected permissions 0600, got %04o", info.Mode().Perm())
		}
	}
}

func TestLoadOrCreateKey_LoadsExistingKey(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "master.key")

	// Generate and save the first key
	key1, err := crypto.LoadOrCreateKey(keyPath)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	// Load it again
	key2, err := crypto.LoadOrCreateKey(keyPath)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	// They must be identical
	if string(key1) != string(key2) {
		t.Error("LoadOrCreateKey did not return the existing key")
	}
}

func TestLoadOrCreateKey_CreatesParentDirs(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "nested", "deep", "master.key")

	_, err := crypto.LoadOrCreateKey(keyPath)
	if err != nil {
		t.Fatalf("failed to create key in nested dir: %v", err)
	}

	// Verify parent dir permissions (0700 expected)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Dir(keyPath))
		if err != nil {
			t.Fatalf("stat parent dir: %v", err)
		}
		if info.Mode().Perm() != 0o700 {
			t.Errorf("expected parent dir permissions 0700, got %04o", info.Mode().Perm())
		}
	}
}

func TestLoadOrCreateKey_WrongSize_Fails(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "master.key")

	// Write a bad key
	if err := os.WriteFile(keyPath, []byte("too-short"), 0o600); err != nil {
		t.Fatalf("write bad key: %v", err)
	}

	_, err := crypto.LoadOrCreateKey(keyPath)
	if err == nil {
		t.Error("expected error when loading key file with wrong size")
	}
}
