# Milestone 1: Foundation

## Timeline
**Target:** Q3 2026 (~2 months)

## Goal
Establish the project infrastructure, foundational Go module setup, configuration management, and build pipelines. Ensure the application can be cleanly packaged for Debian-based systems.

## Key Deliverables
1. **Go Module Initialization:** Repository setup with proper folder structure (`cmd`, `internal`, `pkg`).
2. **System Config Parser:** Implementation of SPEC-003 system config loading (`hostanything.toml`).
3. **Logging System:** Structured JSON logging setup (e.g., using `slog` or `zap`).
4. **Basic REST API:** Core server running with a `/health` endpoint to verify connectivity.
5. **CI/CD Pipeline:** GitHub Actions for linting, testing, and building binaries.
6. **Debian Packaging:** A functioning script to generate a `.deb` package containing the binary, systemd service file, and default configurations.

## Success Criteria
- Running the generated `.deb` file on a clean Debian/Ubuntu VM successfully installs `hostanything`.
- The systemd service starts without errors and persists through reboots.
- The `/api/v1/health` endpoint returns `200 OK` with JSON payload.
- CI pipeline achieves >80% test coverage on configuration parsing logic.

## Out of Scope
- Template parsing or validation.
- Any container runtime integration.
- Frontend web UI.
- Authentication mechanisms.
