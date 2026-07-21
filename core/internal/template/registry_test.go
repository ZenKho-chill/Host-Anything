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

package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/host-anything/hostanything/internal/template"
)

func TestRegistry_ListAndGet(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup mock registry structure
	// redis/1.0.0/template.toml
	redisDir := filepath.Join(tmpDir, "redis", "1.0.0")
	if err := os.MkdirAll(redisDir, 0o755); err != nil {
		t.Fatal(err)
	}
	redisTOML := `
[meta]
name = "redis"
version = "1.0.0"
description = "mock redis"
author = "test"
license = "MIT"
[runtime]
supported = ["docker"]
image = "redis"
`
	if err := os.WriteFile(filepath.Join(redisDir, "template.toml"), []byte(redisTOML), 0o644); err != nil {
		t.Fatal(err)
	}

	reg, err := template.NewRegistry(tmpDir)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}

	// Test List
	summaries, err := reg.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Name != "redis" {
		t.Errorf("expected name redis, got %s", summaries[0].Name)
	}

	// Test Get specific version
	tmpl, err := reg.Get("redis", "1.0.0")
	if err != nil {
		t.Fatalf("Get 1.0.0 failed: %v", err)
	}
	if tmpl.Meta.Description != "mock redis" {
		t.Errorf("expected desc 'mock redis', got %s", tmpl.Meta.Description)
	}

	// Test Get latest (empty version)
	tmplLatest, err := reg.Get("redis", "")
	if err != nil {
		t.Fatalf("Get latest failed: %v", err)
	}
	if tmplLatest.Meta.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", tmplLatest.Meta.Version)
	}

	// Test Get not found
	_, err = reg.Get("not-exist", "1.0.0")
	if err == nil {
		t.Error("expected error for non-existent template")
	}
}

func TestRegistry_DirectoryTraversal_Blocked(t *testing.T) {
	tmpDir := t.TempDir()
	reg, _ := template.NewRegistry(tmpDir)

	_, err := reg.Get("../../../etc", "")
	if err == nil {
		t.Error("expected directory traversal to be blocked")
	}
}
