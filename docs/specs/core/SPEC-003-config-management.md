# SPEC-003: Config Management

## Status
Approved

## Overview
This specification details how Host Anything manages configurations at both the system level and the per-service level. It ensures that configurations are durable, secure, and can be applied to runtime services consistently.

## Motivation
Clear separation between global host settings and individual app settings is required. Furthermore, secret handling needs a structured approach so passwords and tokens are never accidentally exposed via UI or logs.

## Scope
- System-wide configuration (`hostanything.toml`).
- Per-service configuration (`data/services/{id}/config.toml`).
- Configuration validation and application to running services.
- Secret handling.

## Out of Scope
- Specifics of how the web UI renders the configuration forms.

## Specification

### 1. System Config (`/etc/hostanything/hostanything.toml`)
Manages global behavior of the Host Anything daemon.

```toml
[server]
api_port = 8080
bind_address = "127.0.0.1"
log_level = "info" # debug, info, warn, error

[auth]
jwt_secret = "auto-generated-on-install"
session_timeout = "24h"
fail2ban_enabled = true

[runtimes]
docker_enabled = true
podman_enabled = false
k8s_enabled = false
host_enabled = true

[paths]
data_dir = "/var/lib/hostanything/data"
template_dir = "/var/lib/hostanything/templates"
```

### 2. Service Config (`/var/lib/hostanything/data/services/{id}/config.toml`)
Manages the user-provided variables for a specific instantiated service.

```toml
[service]
id = "123e4567-e89b-12d3-a456-426614174000"
template_name = "redis"
template_version = "1.0.0"
runtime_used = "docker"
status = "RUNNING"

[env]
MAX_MEMORY = "512mb"
# Secrets are stored encrypted at rest
REDIS_PASSWORD = "enc:v1:aes256gcm:base64encodedciphertext"

[ports]
6379 = 6379
```

### Config Validation and Updates
When a user updates a service configuration:
1. Core parses the update payload.
2. Core cross-references the `template_name` and `template_version` to load the template schema (SPEC-001).
3. Core validates the provided values against the `type`, `required`, and `validation_regex` fields in the template.
4. If validation passes, the values are written to disk.
5. If the service is RUNNING, the `Update Config` transition (SPEC-002) is triggered. The service is restarted based on the `[update].strategy`.

### Secret Handling (Encryption at Rest)
- When a variable is defined as `type = "secret"` in the template, Host Anything intercepts the plaintext value at the API layer.
- The value is encrypted using AES-256-GCM.
- The key used is derived from an internal master key generated during initial installation (stored in a secure host file with 0600 permissions).
- When applying config to the runtime adapter, the core decrypts the value in memory and passes it securely to the adapter environment injection process.

## Error Handling
- Invalid system configuration on startup will cause a fatal exit with descriptive logs.
- Invalid service config updates via API will return `400 Bad Request` with field-level validation errors.

## Security
- System configuration files are restricted to root/admin (0600).
- Secrets are never returned in plaintext by any API endpoint. GET requests for config will return placeholders (e.g., `********`) for secret fields.

## Testing Strategy
- E2E tests verifying that secret variables are encrypted on disk.
- API tests ensuring secret values are masked on retrieval.
- Validation tests confirming regex checks correctly block invalid user input.
