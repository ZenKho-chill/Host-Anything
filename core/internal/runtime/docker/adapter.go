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

package docker

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/host-anything/hostanything/pkg/types"
)

// Adapter implements types.RuntimeAdapter for Docker.
type Adapter struct {
	cli *client.Client
}

// NewAdapter creates a new Docker Runtime Adapter, connecting to the daemon.
func NewAdapter() (*Adapter, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker.NewAdapter: %w", err)
	}
	return &Adapter{cli: cli}, nil
}

// Deploy pulls the image and creates+starts the container.
func (a *Adapter) Deploy(ctx context.Context, spec *types.ServiceSpec) error {
	imageName := spec.Template.Runtime.Image

	// Pull image (assuming public for now)
	out, err := a.cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("docker.Deploy: image pull %q: %w", imageName, err)
	}
	defer out.Close()
	io.Copy(io.Discard, out) // Wait for pull to complete

	// Convert env map to slice
	var env []string
	for k, v := range spec.ResolvedEnv {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Port bindings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, netCfg := range spec.Template.Network {
		proto := netCfg.Protocol
		if proto == "" || proto == "http" {
			proto = "tcp"
		}
		portStr := strconv.Itoa(netCfg.InternalPort)
		p, err := nat.NewPort(proto, portStr)
		if err != nil {
			return fmt.Errorf("docker.Deploy: invalid port %s/%s: %w", portStr, proto, err)
		}
		exposedPorts[p] = struct{}{}

		var binding []nat.PortBinding
		if netCfg.ExternalPort > 0 {
			binding = []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(netCfg.ExternalPort),
			}}
		} else {
			// Map to random ephemeral port
			binding = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "0"}}
		}
		portBindings[p] = binding
	}

	// Volumes
	var binds []string
	// For M3, we assume volumes map to named docker volumes using ServiceID prefix
	for _, vol := range spec.Template.Volumes {
		binds = append(binds, fmt.Sprintf("ha_%s_%s:%s", spec.ServiceID, vol.Name, vol.MountPath))
	}

	containerName := "ha-" + spec.ServiceID
	config := &container.Config{
		Image:        imageName,
		Env:          env,
		ExposedPorts: exposedPorts,
		Cmd:          spec.Template.Runtime.Command,
		Labels: map[string]string{
			"sh.hostanything.managed": "true",
			"sh.hostanything.service": spec.ServiceID,
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	// Create the container
	resp, err := a.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		// If container exists, maybe it's an update scenario, but deploy assumes fresh.
		return fmt.Errorf("docker.Deploy: create %q: %w", containerName, err)
	}

	// Start it
	if err := a.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("docker.Deploy: start %q: %w", containerName, err)
	}

	return nil
}

// Stop gracefully stops the container.
func (a *Adapter) Stop(ctx context.Context, serviceID string) error {
	containerName := "ha-" + serviceID
	if err := a.cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
		return fmt.Errorf("docker.Stop: %w", err)
	}
	return nil
}

// Start resumes a stopped container.
func (a *Adapter) Start(ctx context.Context, serviceID string) error {
	containerName := "ha-" + serviceID
	if err := a.cli.ContainerStart(ctx, containerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("docker.Start: %w", err)
	}
	return nil
}

// Remove forcefully deletes the container and its anonymous volumes.
func (a *Adapter) Remove(ctx context.Context, serviceID string) error {
	containerName := "ha-" + serviceID
	if err := a.cli.ContainerRemove(ctx, containerName, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return fmt.Errorf("docker.Remove: %w", err)
	}
	return nil
}

// Status inspects the container.
func (a *Adapter) Status(ctx context.Context, serviceID string) (*types.ServiceStatus, error) {
	containerName := "ha-" + serviceID
	info, err := a.cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return nil, fmt.Errorf("docker.Status: %w", err)
	}

	state := types.ServiceStateError
	if info.State.Running {
		state = types.ServiceStateRunning
	} else if info.State.Status == "exited" {
		if info.State.ExitCode == 0 {
			state = types.ServiceStateStopped
		}
	}

	mappings := make(map[int]int)
	for p, bindings := range info.NetworkSettings.Ports {
		if len(bindings) > 0 {
			internal := p.Int()
			external, _ := strconv.Atoi(bindings[0].HostPort)
			mappings[internal] = external
		}
	}

	return &types.ServiceStatus{
		State:        state,
		Message:      info.State.Error,
		PortMappings: mappings,
		RuntimeID:    info.ID,
	}, nil
}

// Logs returns the stdout/stderr streams.
func (a *Adapter) Logs(ctx context.Context, serviceID string) (io.ReadCloser, error) {
	containerName := "ha-" + serviceID
	rc, err := a.cli.ContainerLogs(ctx, containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("docker.Logs: %w", err)
	}
	return rc, nil
}
