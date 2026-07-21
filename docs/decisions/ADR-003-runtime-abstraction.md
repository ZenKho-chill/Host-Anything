# ADR 003: Runtime Abstraction via Adapter Pattern

## Status
Accepted

## Context
Host Anything aims to manage services. The current landscape of service hosting is fragmented: some users use Docker, some prefer daemonless Podman, enterprise users use Kubernetes, and some simple services run best as bare metal processes. If the core logic is tightly coupled to Docker, supporting Podman or K8s later will require a massive rewrite.

## Decision
We will implement a strict **Adapter Pattern** for all runtime executions. The core daemon will only interact with a generic `RuntimeAdapter` interface.

## Rationale
By forcing the core configuration manager to translate user intent into a generic `ServiceSpec` struct, we completely isolate the business logic (auth, API, state persistence, template resolution) from the operational logic (container API calls, systemd management). 

## Consequences

**Pros:**
- **Extensibility**: Adding a new runtime (e.g., LXC or Nomad) requires only writing a new struct that satisfies the `RuntimeAdapter` interface.
- **Testability**: We can easily create a `MockAdapter` for unit testing the core manager without needing a real Docker daemon running in CI.
- **Safety**: The core system is prevented from accidentally relying on Docker-specific quirks.

**Trade-offs:**
- **Lowest Common Denominator**: We can only expose features that can be reasonably implemented across *most* adapters. Highly specific orchestration features (like Kubernetes specific tolerations or Docker Swarm overlay networks) cannot be easily surfaced in the generic template format.
- **Development Overhead**: Requires maintaining multiple distinct adapter implementations and writing extensive integration tests for each environment.
