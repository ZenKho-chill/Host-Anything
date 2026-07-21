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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/host-anything/hostanything/internal/api"
)

func TestHealthHandler_ReturnsHTTP200(t *testing.T) {
	handler := api.HealthHandler("0.1.0", []string{"docker"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHealthHandler_ContentTypeIsJSON(t *testing.T) {
	handler := api.HealthHandler("0.1.0", []string{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type=application/json, got %q", ct)
	}
}

func TestHealthHandler_BodyStatusIsUp(t *testing.T) {
	handler := api.HealthHandler("0.1.0", []string{"docker"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["status"] != "up" {
		t.Errorf("expected status=up, got %v", body["status"])
	}
}

func TestHealthHandler_BodyVersionMatches(t *testing.T) {
	handler := api.HealthHandler("1.2.3", nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["version"] != "1.2.3" {
		t.Errorf("expected version=1.2.3, got %v", body["version"])
	}
}

func TestHealthHandler_NilRuntimes_SerializesAsEmptyArray(t *testing.T) {
	handler := api.HealthHandler("0.1.0", nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	runtimes, ok := body["runtimes"].([]interface{})
	if !ok {
		t.Fatalf("expected runtimes to be a JSON array, got %T", body["runtimes"])
	}
	if len(runtimes) != 0 {
		t.Errorf("expected empty runtimes array, got %v", runtimes)
	}
}

func TestHealthHandler_RuntimesArePresentInBody(t *testing.T) {
	handler := api.HealthHandler("0.1.0", []string{"docker", "host"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	runtimes, ok := body["runtimes"].([]interface{})
	if !ok {
		t.Fatalf("expected runtimes to be a JSON array, got %T", body["runtimes"])
	}
	if len(runtimes) != 2 {
		t.Errorf("expected 2 runtimes, got %d", len(runtimes))
	}
}

func TestHealthHandler_IsReusable_SameResponseEachCall(t *testing.T) {
	handler := api.HealthHandler("0.1.0", []string{"docker"})

	body1 := callHandler(t, handler)
	body2 := callHandler(t, handler)

	if body1 != body2 {
		t.Errorf("expected identical responses on repeated calls, got:\n%s\n%s", body1, body2)
	}
}

// callHandler is a test helper that invokes a handler and returns the raw body string.
func callHandler(t *testing.T, h http.HandlerFunc) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec.Body.String()
}
