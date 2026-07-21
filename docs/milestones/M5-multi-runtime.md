# Milestone 5: Multi-Runtime

## Timeline
**Target:** Q3 2027 (~3 months)

## Goal
Expand Host Anything beyond Docker. Implement adapters for Podman (daemonless containers), Kubernetes, and Host processes to provide ultimate deployment flexibility.

## Key Deliverables
1. **Runtime Auto-Detection:** Core engine can detect which runtimes are available on the host machine.
2. **Podman Adapter:** Implementation using Podman's REST API or Go bindings.
3. **Kubernetes Adapter:** Implementation using client-go to create Deployments/Services.
4. **Host Adapter:** Implementation utilizing `os/exec` and systemd to run raw binaries directly on the host.
5. **Runtime Migration:** Ability to export state and move a service from Docker to Podman (or vice-versa).

## Success Criteria
- Seamless deployment of a service to Podman without modifying the underlying TOML template.
- K8s adapter successfully spins up a Pod and exposes it via a Service based on the template network config.
- Web UI allows users to select the desired runtime during deployment if multiple are supported by the template.

## Out of Scope
- K8s cluster management (Host Anything assumes the kubeconfig is pre-provisioned).
- Advanced K8s constructs like Ingress or PersistentVolumeClaims (simple NodePort/HostPath utilized initially).
