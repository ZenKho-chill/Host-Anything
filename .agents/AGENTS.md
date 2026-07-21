# Host Anything - Project Coding Rules and AI Behavior Guidelines

## Project Context
**Host Anything** is a self-hosted platform for deploying and managing services effortlessly. The core system manages configuration of services (env vars, ports, volumes, resource limits) but does NOT interfere with service runtime logic. It relies on a GitHub-based template marketplace and supports Docker, Podman, Kubernetes, and host mode via an adapter pattern.

**Tech Stack**:
- **Core**: Go
- **Web UI**: React + TypeScript
- **Templates**: TOML
- **Distribution**: Debian (.deb packages)
- **Security**: Web UI with authentication + fail2ban integration

## Coding Rules (strictly enforced)
1. Each file must contain ONLY ONE function or ONE group of functions from the same domain.
2. If a function has multiple variants (e.g., parse TOML vs parse YAML), each variant MUST be in a separate file (`parse_toml.go`, `parse_yaml.go` — NOT `parse.go` with both).
3. No global mutable state — all dependencies passed via constructor or function parameters.
4. All errors must be wrapped with context using `fmt.Errorf("package.FunctionName: %w", err)`.
5. All exported symbols must have godoc comments.
6. Test files live in the same package (not in separate test/ folder).
7. Define interfaces before implementations — interfaces go in `pkg/types/`.
8. File names use `snake_case` for Go, `kebab-case` for TS/React.
9. No magic numbers — use named constants in a `constants.go` file within the package.
10. Configuration validation must be strict — reject unknown fields.

## Architecture Layers (must not be violated)
- `cmd/` → only wires dependencies, no business logic.
- `internal/` → business logic, not importable outside module.
- `pkg/` → shared types and utilities, no external deps.
- Runtime adapters must implement the `RuntimeAdapter` interface from `pkg/types`.
- API handlers must not contain business logic — delegate to `internal/` services.

## Documentation Rules
- Every new package must have a `doc.go` file with package documentation.
- All public-facing API changes require updating `docs/api/openapi.yaml`.
- All significant decisions require an ADR in `docs/decisions/`.
- Spec files in `docs/specs/` must be updated before implementation.

## Git Rules
- Conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`
- Branch naming: `feat/description`, `fix/description`, `docs/description`
- No direct commits to main.
- PR must reference an issue.
