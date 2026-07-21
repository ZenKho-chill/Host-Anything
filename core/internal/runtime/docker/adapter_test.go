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

//go:build integration
// +build integration

package docker_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/host-anything/hostanything/internal/runtime/docker"
	"github.com/host-anything/hostanything/pkg/types"
)

func TestDockerAdapter_Lifecycle(t *testing.T) {
	// Note: testcontainers-go is great for standing up dependencies (like DBs),
	// but here we want to test our actual Docker Adapter's ability to talk to the
	// Docker daemon. We will use the adapter itself to deploy a tiny alpine container
	// and verify the lifecycle.

	adapter, err := docker.NewAdapter()
	if err != nil {
		t.Skipf("skipping test: failed to create docker adapter (docker not running or no permissions): %v", err)
	}

	ctx := context.Background()
	// Ping daemon to check if it's actually alive
	if _, pingErr := adapter.Status(ctx, "nonexistent"); pingErr != nil && strings.Contains(pingErr.Error(), "connect: ") {
		t.Skipf("skipping test: docker daemon unreachable: %v", pingErr)
	}
	serviceID := "test-alpine"

	spec := &types.ServiceSpec{
		ServiceID: serviceID,
		Template: &types.Template{
			Runtime: types.RuntimeConfig{
				Image:   "alpine:latest",
				Command: []string{"sh", "-c", "echo \"hello $MYVAR\" && sleep 3600"},
			},
		},
		ResolvedEnv: map[string]string{
			"MYVAR": "hostanything",
		},
	}

	// 1. Deploy
	if err := adapter.Deploy(ctx, spec); err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}
	defer adapter.Remove(context.Background(), serviceID)

	// Give it a second to start
	time.Sleep(1 * time.Second)

	// 2. Status
	status, err := adapter.Status(ctx, serviceID)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if status.State != types.ServiceStateRunning {
		t.Errorf("Expected state %s, got %s", types.ServiceStateRunning, status.State)
	}

	// 3. Logs
	logs, err := adapter.Logs(ctx, serviceID)
	if err != nil {
		t.Fatalf("Logs failed: %v", err)
	}
	logBytes, _ := io.ReadAll(logs)
	logs.Close()
	if !strings.Contains(string(logBytes), "hello hostanything") {
		t.Errorf("Expected logs to contain env output, got: %q", string(logBytes))
	}

	// 4. Stop
	if err := adapter.Stop(ctx, serviceID); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	status, _ = adapter.Status(ctx, serviceID)
	if status.State != types.ServiceStateStopped {
		t.Errorf("Expected state %s after stop, got %s", types.ServiceStateStopped, status.State)
	}
}
