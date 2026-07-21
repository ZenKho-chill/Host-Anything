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

package api

import (
	"encoding/json"
	"net/http"
)

// healthResponse is the JSON body returned by GET /api/v1/health.
// Schema matches SPEC-004 §1.
type healthResponse struct {
	Status   string   `json:"status"`
	Version  string   `json:"version"`
	Runtimes []string `json:"runtimes"`
}

// HealthHandler returns an [http.HandlerFunc] for the GET /api/v1/health endpoint.
// It responds with 200 OK and a JSON body containing the daemon version and the
// list of currently enabled runtime adapters.
//
// This endpoint is intentionally unauthenticated per SPEC-004 — it is used
// for liveness probes and connectivity checks.
//
// The response is pre-encoded at construction time since it is static for the
// lifetime of the server.
func HealthHandler(version string, runtimes []string) http.HandlerFunc {
	resp := healthResponse{
		Status:   "up",
		Version:  version,
		Runtimes: runtimes,
	}
	// Ensure the JSON array is always [] rather than null when empty.
	if resp.Runtimes == nil {
		resp.Runtimes = []string{}
	}

	body, err := json.Marshal(resp)
	if err != nil {
		// json.Marshal on a static struct with known-good types cannot fail.
		panic("api.HealthHandler: failed to pre-encode health response: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}
