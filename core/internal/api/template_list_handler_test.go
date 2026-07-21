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

	"github.com/host-anything/hostanything/internal/api"
	"github.com/host-anything/hostanything/internal/template"
	"github.com/host-anything/hostanything/pkg/types"
)

func TestTemplateListHandler(t *testing.T) {
	tmpDir := t.TempDir()
	reg, _ := template.NewRegistry(tmpDir)

	// Add a dummy template
	d := filepath.Join(tmpDir, "dummy", "1.0.0")
	os.MkdirAll(d, 0o755)
	toml := `[meta]
name = "dummy"
version = "1.0.0"
description = "desc"
author = "author"
license = "MIT"
[runtime]
supported = ["docker"]
image = "dummy"
`
	os.WriteFile(filepath.Join(d, "template.toml"), []byte(toml), 0o644)

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	handler := api.TemplateListHandler(reg, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/templates", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var res []types.TemplateSummary
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res) != 1 || res[0].Name != "dummy" {
		t.Errorf("expected 1 template named dummy, got %v", res)
	}
}
