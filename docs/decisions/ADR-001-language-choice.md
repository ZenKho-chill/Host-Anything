# ADR 001: Language Choice - Go

## Status
Accepted

## Context
Host Anything needs to be a highly reliable, low-footprint background daemon capable of managing complex state, communicating with local sockets (Docker/Podman), and serving a web API concurrently. The language choice must support rapid development while maintaining strict safety guarantees and easy distribution for end-users (primarily Debian/Ubuntu home server admins).

Options considered:
- **Python**: Great for scripting, but distributing as a single standalone executable without dependency hell is challenging. Concurrency model (GIL) is not ideal for managing multiple simultaneous deployment streams.
- **Rust**: Excellent performance and memory safety, but compilation times and the steep learning curve could slow down feature velocity.
- **Node.js**: Easy API building, but requires a heavy runtime installation on the host system.
- **Go**: Designed for network services, excellent concurrency (goroutines), compiles to a single static binary.

## Decision
We will use **Go** for the core backend system.

## Consequences

**Pros:**
- **Single Static Binary**: The entire backend, including the embedded frontend, can be shipped as a single executable. This makes `.deb` packaging trivial.
- **Concurrency**: Goroutines make handling multiple API requests, long-running Docker API streams (logs), and background health checks straightforward.
- **Standard Library**: Go's robust `net/http` means we don't need heavy third-party frameworks to serve the API.
- **Ecosystem**: Excellent existing client libraries for Docker (`docker/client`) and Kubernetes (`client-go`).

**Cons:**
- **Error Handling**: Verbose `if err != nil` boilerplate. We must adopt a strict pattern of wrapping errors with context (`fmt.Errorf("doing x: %w", err)`) to maintain traceability.
- **Generics**: While available in 1.18+, they can lead to unreadable code if overused. We will restrict generics to strictly necessary utility functions.
