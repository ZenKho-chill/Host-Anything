# Core System Design

## Overview
The Core Daemon (`hostanything`) is a long-running background service written in Go. It acts as the central brain of the platform, managing all interactions between the user interface, the template marketplace, and the local container/process runtimes.

## Core Daemon Responsibilities
1. **Serve REST API**: Expose a secure, versioned API for the Web UI and potential CLI clients.
2. **Manage Templates**: Fetch, cache, validate, and parse TOML templates from local disk or the GitHub marketplace.
3. **Manage Service Configs**: Apply user-provided overrides (environment variables, port mappings, volumes, resource constraints) to template baselines.
4. **Handle Authentication**: Authenticate users, issue JWTs, and log security events (failed logins) for Fail2ban.
5. **Orchestrate Adapters**: Translate abstract service states into concrete deployments via the Runtime Adapters.

## Package Structure

The Go codebase follows standard Go project layout conventions:

```text
hostanything/
├── cmd/
│   └── hostanything/       # Main application entrypoint
│       └── main.go         # Bootstraps the daemon, initializes dependencies
├── internal/
│   ├── api/                # REST API handlers, routing, middleware (JWT)
│   ├── auth/               # User authentication, password hashing, JWT generation
│   ├── config/             # System configuration (daemon settings)
│   ├── manager/            # Core business logic, configuration application
│   ├── state/              # SQLite database interactions and models
│   └── template/           # TOML template parsing, validation, GitHub fetching
├── pkg/
│   └── runtime/            # Runtime adapter interfaces and implementations
│       ├── docker/
│       ├── podman/
│       ├── k8s/
│       └── host/
└── ui/                     # Embedded React static files (post-build)
```

## Configuration Management Design
Configuration management is decoupled from runtime execution. 
1. A **Template** defines the default, static requirements of a service.
2. A **Service Configuration** represents the user's specific instantiation of that template.

When a user updates a service (e.g., changes an environment variable):
1. The Core API receives the PATCH request.
2. The `manager` package updates the service record in the state database.
3. The `manager` fetches the underlying Template.
4. It merges the User Config over the Template Defaults.
5. It invokes `adapter.ApplyConfig()` to enforce the new state on the runtime environment.

## State Persistence Approach
**Decision: Embedded SQLite**

While flat files (JSON/YAML) are simpler for basic configs, managing relationships (Users -> Services -> Backups) and concurrent access quickly becomes error-prone. 
SQLite provides:
- ACID compliance preventing state corruption during unexpected shutdowns.
- Relational integrity.
- Single-file simplicity (easy to backup).
- Excellent Go integration via `mattn/go-sqlite3` or pure Go alternatives like `modernc.org/sqlite`.

State is stored typically in `/var/lib/hostanything/state.db` on Debian systems.

## Startup Sequence
1. **Initialization**: Read daemon config (`/etc/hostanything/config.toml`), parse CLI flags.
2. **Storage Setup**: Connect to SQLite database, run auto-migrations.
3. **Security Setup**: Initialize auth module, ensure default admin exists or prompt for creation if fresh install. Prepare Fail2ban log file (`/var/log/hostanything/auth.log`).
4. **Adapter Initialization**: Probe the host system to determine available runtimes (e.g., check for Docker socket, Podman socket, kubectl). Initialize applicable adapters.
5. **State Reconciliation**: Load all configured services from SQLite. Query active adapters for current status. Re-apply configurations if auto-start is enabled and services are missing.
6. **Server Start**: Bind to configured host/port (default: 127.0.0.1:8080) and start serving the REST API and embedded Web UI.
