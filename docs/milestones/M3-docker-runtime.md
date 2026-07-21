# Milestone 3: Docker Runtime

**STATUS: COMPLETED**

## Timeline
**Target:** Q1 2027 (~3 months)

## Goal
Bridge the gap between Host Anything's internal service representation and actual process execution by fully implementing the Docker Runtime Adapter. Enable the full lifecycle (SPEC-002) for Docker containers.

## Key Deliverables
1. **RuntimeAdapter Interface:** Define the Go interface required for any runtime (`Deploy`, `Start`, `Stop`, `Remove`, `Status`, `Logs`, `ApplyConfig`).
2. **Docker Implementation:** Build the Docker adapter utilizing the official Docker Go SDK.
3. **Lifecycle Management:** Implement the state machine (SPEC-002) managing transitions (Deploying -> Running -> Stopping).
4. **Environment Injection:** Pass parsed, decrypted variables from the Template Engine directly into container environments.
5. **Volume & Port Mapping:** Translate template specifications into Docker host bindings.

## Success Criteria
- Given a valid template and user config, the system successfully pulls an image, creates a container, and starts it.
- Container health changes accurately reflect in the Host Anything state machine.
- Changing a configuration variable correctly stops, recreates, and restarts the container (if strategy is recreate).
- High coverage of integration tests utilizing `testcontainers-go`.

## Out of Scope
- Podman, K8s, or Host mode implementations.
- Multi-container templates (e.g., Docker Compose equivalents). Host Anything manages single service units; complex apps use multiple templates.
