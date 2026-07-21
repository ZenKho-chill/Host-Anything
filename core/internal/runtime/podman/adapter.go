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

package podman

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/host-anything/hostanything/pkg/types"
)

// Adapter implements types.RuntimeAdapter for Podman via CLI.
type Adapter struct{}

// NewAdapter creates a new Podman Adapter.
func NewAdapter() (*Adapter, error) {
	// Verify podman is installed
	if _, err := exec.LookPath("podman"); err != nil {
		return nil, fmt.Errorf("podman not found in PATH: %w", err)
	}
	return &Adapter{}, nil
}

func (a *Adapter) Deploy(ctx context.Context, spec *types.ServiceSpec) error {
	containerName := "ha-" + spec.ServiceID

	// Ensure any old container is removed
	_ = a.Remove(ctx, spec.ServiceID)

	args := []string{"run", "-d", "--name", containerName, "--replace"}

	// Network
	for _, netCfg := range spec.Template.Network {
		if netCfg.ExternalPort > 0 {
			args = append(args, "-p", fmt.Sprintf("%d:%d", netCfg.ExternalPort, netCfg.InternalPort))
		}
	}

	// Env
	for k, v := range spec.ResolvedEnv {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	// Volumes
	for _, vol := range spec.Template.Volumes {
		args = append(args, "-v", fmt.Sprintf("ha_%s_%s:%s", spec.ServiceID, vol.Name, vol.MountPath))
	}

	// Labels
	args = append(args, "-l", "sh.hostanything.managed=true")
	args = append(args, "-l", "sh.hostanything.service="+spec.ServiceID)

	args = append(args, spec.Template.Runtime.Image)
	if len(spec.Template.Runtime.Command) > 0 {
		args = append(args, spec.Template.Runtime.Command...)
	}

	cmd := exec.CommandContext(ctx, "podman", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("podman deploy failed: %w, output: %s", err, string(out))
	}

	return nil
}

func (a *Adapter) Stop(ctx context.Context, serviceID string) error {
	cmd := exec.CommandContext(ctx, "podman", "stop", "ha-"+serviceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("podman stop: %w", err)
	}
	return nil
}

func (a *Adapter) Start(ctx context.Context, serviceID string) error {
	cmd := exec.CommandContext(ctx, "podman", "start", "ha-"+serviceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("podman start: %w", err)
	}
	return nil
}

func (a *Adapter) Remove(ctx context.Context, serviceID string) error {
	cmd := exec.CommandContext(ctx, "podman", "rm", "-f", "-v", "ha-"+serviceID)
	if err := cmd.Run(); err != nil {
		// Ignore not found errors during removal
		return nil
	}
	return nil
}

func (a *Adapter) Status(ctx context.Context, serviceID string) (*types.ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "podman", "inspect", "-f", "{{.State.Status}}", "ha-"+serviceID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("podman inspect: %w", err)
	}

	status := strings.TrimSpace(string(out))
	state := types.ServiceStateError
	if status == "running" {
		state = types.ServiceStateRunning
	} else if status == "exited" || status == "stopped" {
		state = types.ServiceStateStopped
	}

	return &types.ServiceStatus{
		State:        state,
		PortMappings: make(map[int]int),
		RuntimeID:    "ha-" + serviceID,
	}, nil
}

func (a *Adapter) Logs(ctx context.Context, serviceID string) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "podman", "logs", "-f", "ha-"+serviceID)

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
