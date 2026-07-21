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
	"testing"

	"github.com/host-anything/hostanything/internal/crypto"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

func TestResolve_Defaults(t *testing.T) {
	tmpl := &types.Template{
		Config: []types.ConfigVar{
			{Name: "VAR_A", Type: template.ConfigTypeString, Required: true},
			{Name: "VAR_B", Type: template.ConfigTypeInt, Default: 42},
		},
	}

	userVars := map[string]string{
		"VAR_A": "user_value",
	}

	resolved, err := template.Resolve(tmpl, userVars, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved["VAR_A"] != "user_value" {
		t.Errorf("expected user_value, got %q", resolved["VAR_A"])
	}
	if resolved["VAR_B"] != "42" {
		t.Errorf("expected default 42, got %q", resolved["VAR_B"])
	}
}

func TestResolve_SecretEncryption(t *testing.T) {
	tmpl := &types.Template{
		Config: []types.ConfigVar{
			{Name: "MY_SECRET", Type: template.ConfigTypeSecret, Required: true},
		},
	}

	userVars := map[string]string{
		"MY_SECRET": "super-secret",
	}

	key, _ := crypto.GenerateKey()

	// Resolve with key (should encrypt)
	resolvedWithKey, err := template.Resolve(tmpl, userVars, key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !crypto.IsEncrypted(resolvedWithKey["MY_SECRET"]) {
		t.Errorf("expected encrypted value, got %q", resolvedWithKey["MY_SECRET"])
	}

	// Resolve without key (should leave plain)
	resolvedPlain, err := template.Resolve(tmpl, userVars, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if crypto.IsEncrypted(resolvedPlain["MY_SECRET"]) {
		t.Error("expected plain text when key is nil")
	}
	if resolvedPlain["MY_SECRET"] != "super-secret" {
		t.Errorf("expected original secret, got %q", resolvedPlain["MY_SECRET"])
	}
}
