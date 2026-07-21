# Runtime Adapter System

## Overview
The Runtime Adapter System is the core mechanism that allows Host Anything to deploy services across vastly different execution environments (Docker, Podman, Kubernetes, Host) without tying the core business logic to any specific orchestrator.

## The `RuntimeAdapter` Interface
All interactions with execution environments occur through a strict Go interface defined in `pkg/runtime`.

```go
package runtime

import "context"

type Adapter interface {
    // Name returns the identifier for this adapter (e.g., "docker", "podman")
    Name() string
    
    // Deploy instantiates a new service based on the generalized spec
    Deploy(ctx context.Context, spec *ServiceSpec) error
    
    // Stop gracefully halts a running service
    Stop(ctx context.Context, serviceID string) error
    
    // Start resumes a stopped service
    Start(ctx context.Context, serviceID string) error
    
    // Status returns the current health and state of the service
    Status(ctx context.Context, serviceID string) (*ServiceStatus, error)
    
    // Logs streams stdout/stderr from the service
    Logs(ctx context.Context, serviceID string, options LogOptions) (LogStream, error)
    
    // ApplyConfig updates a running service with new configuration parameters
    ApplyConfig(ctx context.Context, serviceID string, newSpec *ServiceSpec) error
    
    // Remove permanently deletes the service and optionally its volumes
    Remove(ctx context.Context, serviceID string, purgeVolumes bool) error
}
```

## Adapter Selection
Adapters are selected based on a two-pass system:
1. **System Capability Detection**: On startup, Host Anything probes the host. If `/var/run/docker.sock` exists, the Docker adapter is enabled. If `podman` is in the PATH, Podman is enabled, etc.
2. **Template Declaration vs. User Preference**: A template may define a `recommended_runtime`. If the user doesn't specify one during deployment, the system attempts to use the recommended one. If unavailable, it falls back to the user's default configured runtime.

## Docker Adapter
- **Implementation**: Uses the official `github.com/docker/docker/client` Go SDK.
- **Mechanism**: Connects via local unix socket. Translates `ServiceSpec` into Docker API container creation calls.
- **Handling Updates**: `ApplyConfig` typically involves stopping the container, removing it, and recreating it with the new configuration, reattaching the existing persistent volumes.

## Podman Adapter
- **Implementation**: Communicates over the Podman REST API (requires podman system service to be running) or CLI wrapping.
- **Mechanism**: Behaves almost identically to Docker, but leverages rootless containers by default where applicable.
- **Benefits**: Better security model for home lab users who prefer not to run a root-level daemon.

## Kubernetes Adapter
- **Implementation**: Uses `k8s.io/client-go`.
- **Mechanism**: Translates a `ServiceSpec` into standard Kubernetes manifests (Deployment, Service, PersistentVolumeClaim) and applies them to the configured namespace.
- **Complexity**: `ApplyConfig` is handled natively by patching the K8s Deployments, allowing K8s to handle the rolling restart smoothly.

## Host Mode Adapter
- **Implementation**: Uses standard Go `os/exec` and systemd integration.
- **Mechanism**: Used for running raw binaries or scripts directly on the host OS. The adapter generates a transient `systemd` unit file for the service, allowing systemd to handle restarts, logging (via journald), and resource limits (via cgroups).
- **Security**: Heavily restricted. Requires explicit admin approval in the UI, as host mode breaks isolation guarantees.

## Compatibility Matrix

| Feature | Docker | Podman | Kubernetes | Host |
| :--- | :---: | :---: | :---: | :---: |
| Image Pulling | ✅ | ✅ | ✅ | ❌ |
| Port Mapping | ✅ | ✅ | ✅ (NodePort/LB) | ⚠️ (Manual) |
| Env Vars | ✅ | ✅ | ✅ | ✅ |
| Volume Mounts | ✅ | ✅ | ✅ (PVCs) | ✅ (Direct paths) |
| Resource Limits| ✅ | ✅ | ✅ | ✅ (via systemd) |
| Log Streaming | ✅ | ✅ | ✅ | ✅ (journalctl) |
