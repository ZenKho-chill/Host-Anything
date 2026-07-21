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

package host

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/host-anything/hostanything/pkg/types"
)

// Adapter implements types.RuntimeAdapter for raw Host processes using systemd.
type Adapter struct{}

// NewAdapter creates a new Host Adapter.
func NewAdapter() (*Adapter, error) {
	// Verify systemd is available (systemctl exists)
	if _, err := exec.LookPath("systemctl"); err != nil {
		return nil, fmt.Errorf("systemctl not found in PATH: %w", err)
	}
	if _, err := exec.LookPath("systemd-run"); err != nil {
		return nil, fmt.Errorf("systemd-run not found in PATH: %w", err)
	}
	return &Adapter{}, nil
}

func (a *Adapter) Deploy(ctx context.Context, spec *types.ServiceSpec) error {
	unitName := "ha-" + spec.ServiceID + ".service"

	// Ensure any previous transient unit is cleared
	_ = a.Remove(ctx, spec.ServiceID)

	args := []string{
		"--unit=ha-" + spec.ServiceID,
		"--property=Restart=always",
		"--remain-after-exit",
	}

	// Environment variables
	for k, v := range spec.ResolvedEnv {
		args = append(args, fmt.Sprintf("--setenv=%s=%s", k, v))
	}

	// The 'Image' field for host adapter acts as the absolute path to the binary or script.
	// For example, Image: "/usr/local/bin/my-service"
	binaryPath := spec.Template.Runtime.Image
	if binaryPath == "" {
		return fmt.Errorf("host adapter requires 'image' field to specify the binary path")
	}
	args = append(args, binaryPath)

	// Command arguments
	if len(spec.Template.Runtime.Command) > 0 {
		args = append(args, spec.Template.Runtime.Command...)
	}

	cmd := exec.CommandContext(ctx, "systemd-run", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("systemd-run failed: %w, output: %s", err, string(out))
	}

	// Wait for it to start
	if err := a.Start(ctx, spec.ServiceID); err != nil {
		return fmt.Errorf("failed to start unit %s: %w", unitName, err)
	}

	return nil
}

func (a *Adapter) Stop(ctx context.Context, serviceID string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "stop", "ha-"+serviceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("systemctl stop: %w", err)
	}
	return nil
}

func (a *Adapter) Start(ctx context.Context, serviceID string) error {
	cmd := exec.CommandContext(ctx, "systemctl", "start", "ha-"+serviceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("systemctl start: %w", err)
	}
	return nil
}

func (a *Adapter) Remove(ctx context.Context, serviceID string) error {
	unitName := "ha-" + serviceID

	_ = a.Stop(ctx, serviceID)
	
	// Reset failed state to clear transient units
	cmd := exec.CommandContext(ctx, "systemctl", "reset-failed", unitName)
	_ = cmd.Run()
	
	return nil
}

func (a *Adapter) Status(ctx context.Context, serviceID string) (*types.ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", "ha-"+serviceID)
	out, _ := cmd.CombinedOutput()
	
	statusStr := strings.TrimSpace(string(out))
	state := types.ServiceStateError

	// is-active returns non-zero if not active, so we just check the string output
	if statusStr == "active" {
		state = types.ServiceStateRunning
	} else if statusStr == "inactive" || statusStr == "failed" {
		state = types.ServiceStateStopped
	}

	// Host mode doesn't do port mapping natively (it shares host network)
	// so we return an empty map.
	return &types.ServiceStatus{
		State:        state,
		Message:      statusStr,
		PortMappings: make(map[int]int),
		RuntimeID:    "ha-" + serviceID,
	}, nil
}

func (a *Adapter) Logs(ctx context.Context, serviceID string) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "journalctl", "-u", "ha-"+serviceID, "-f", "--output=cat")
	
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		_ = cmd.Wait()
		pw.Close()
	}()

	return pr, nil
}
