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

package api_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/host-anything/hostanything/internal/api"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

func TestTemplateGetHandler(t *testing.T) {
	tmpDir := t.TempDir()
	reg, _ := template.NewRegistry(tmpDir)

	// Add a dummy template with a secret config to test masking
	d := filepath.Join(tmpDir, "dummy", "1.0.0")
	os.MkdirAll(d, 0o755)
	toml := `[meta]
name = "dummy"
version = "1.0.0"
description = "desc"
author = "author"
license = "MIT"
[[config]]
name = "SEC"
type = "secret"
required = false
default = "supersecret"
[runtime]
supported = ["docker"]
image = "dummy"
`
	os.WriteFile(filepath.Join(d, "template.toml"), []byte(toml), 0o644)

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	handler := api.TemplateGetHandler(reg, logger)

	// Create chi router to populate URL param "name"
	r := chi.NewRouter()
	r.Get("/api/v1/templates/{name}", handler)

	// Test Success
	req := httptest.NewRequest(http.MethodGet, "/api/v1/templates/dummy", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var res types.Template
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res.Meta.Name != "dummy" {
		t.Errorf("expected template name dummy, got %s", res.Meta.Name)
	}
	if res.Config[0].Default != "***" {
		t.Errorf("expected secret default to be masked as ***, got %v", res.Config[0].Default)
	}

	// Test Not Found
	req = httptest.NewRequest(http.MethodGet, "/api/v1/templates/not-exist", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}
