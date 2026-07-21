# SPEC-001: Template Schema

## Status
Approved

## Overview
The Template Schema defines the definitive configuration standard for all services managed by Host Anything. Based on TOML, this schema describes metadata, hardware requirements, configuration variables, runtime definitions, networking, storage, and health checks needed to safely deploy and manage a service.

## Motivation
To provide a unified, human-readable, and machine-verifiable format for describing how applications should be run, regardless of the underlying runtime (Docker, Podman, K8s, Host mode). A well-defined schema ensures templates from the community can be safely ingested and deployed.

## Scope
- Complete TOML schema definition for templates.
- Validation rules for fields.
- Definitions for meta, requirements, config variables, runtime, volumes, network, healthcheck, and update sections.

## Out of Scope
- How the template is stored locally or retrieved from the marketplace (see SPEC-040).
- Translation logic into specific runtime payloads (e.g., generating a docker-compose.yml or k8s manifest).

## Specification

A Host Anything template is a TOML document containing several key sections.

### Data Schemas

#### 1. `[meta]`
Contains descriptive information about the template itself.
- `name` (string, required): Unique identifier/slug for the service (e.g., "nginx").
- `version` (string, required): SemVer of the template.
- `description` (string, required): A short explanation of the service.
- `author` (string, required): Maintainer name or handle.
- `license` (string, required): SPDX license identifier.
- `tags` (array of strings, optional): Used for marketplace categorization.
- `homepage` (string, optional): Project upstream URL.

#### 2. `[requirements]`
Defines minimum system resources.
- `min_memory` (string, optional): E.g., "512MB", "1GB".
- `min_cpu` (float, optional): E.g., 0.5 (half a core).
- `ports` (array of integers, optional): Required host ports.
- `disk_space` (string, optional): E.g., "10GB".

#### 3. `[[config]]` (Array of Tables)
Defines user-configurable variables injected into the runtime environment.
- `name` (string, required): Environment variable name (e.g., "DB_PASSWORD").
- `type` (string, required): "string", "int", "boolean", "secret", "enum".
- `default` (any, optional): Default value.
- `required` (boolean, required): Whether the user must provide it.
- `description` (string, optional): Help text for the UI.
- `validation_regex` (string, optional): Regex to validate the input.
- `options` (array of strings, optional): Valid values if type is "enum".

#### 4. `[runtime]`
Configures the execution environment.
- `preferred` (string, optional): "docker", "podman", "k8s", "host".
- `supported` (array of strings, required): List of supported runtimes.
- `image` (string, required): The container image (e.g., "nginx:latest").
- `command` (array of strings, optional): Custom entrypoint/command.

#### 5. `[[volumes]]`
Maps persistent data storage.
- `name` (string, required): Logical name for the volume.
- `mount_path` (string, required): Path inside the container/service.
- `description` (string, optional): What this volume stores.

#### 6. `[[network]]`
Exposed services.
- `internal_port` (int, required): Port the app listens on.
- `external_port` (int, optional): Default host port.
- `protocol` (string, optional): "tcp", "udp", "http". Defaults to "tcp".

#### 7. `[healthcheck]`
Determines if the service is operational.
- `command` (string, required): Command to run inside the container (e.g., "curl -f http://localhost/ || exit 1").
- `interval` (string, optional): E.g., "30s".
- `timeout` (string, optional): E.g., "10s".
- `retries` (int, optional): E.g., 3.

#### 8. `[update]`
Defines how the service should handle updates.
- `strategy` (string, required): "recreate" (stop then start) or "rolling".

## Error Handling
The template engine will return structured validation errors including the section, line number, and a human-readable explanation if a TOML file fails schema validation.

## Security
- `secret` config types must never be logged or exported in plain text.
- Validation regexes will use Go's `regexp` which is safe against ReDoS attacks.

## Testing Strategy
- Unit tests validating correct TOML files.
- Unit tests ensuring missing required fields trigger specific errors.
- Unit tests testing `validation_regex` against edge case inputs.

## Full Example TOML
```toml
[meta]
name = "redis"
version = "1.0.0"
description = "In-memory data structure store"
author = "Host Anything Community"
license = "MIT"
tags = ["database", "cache"]
homepage = "https://redis.io"

[requirements]
min_memory = "256MB"
min_cpu = 0.5
disk_space = "1GB"

[[config]]
name = "REDIS_PASSWORD"
type = "secret"
required = true
description = "Password for Redis auth"
validation_regex = "^.{8,}$"

[[config]]
name = "MAX_MEMORY"
type = "string"
default = "256mb"
required = false
description = "Maximum memory usage"

[runtime]
preferred = "docker"
supported = ["docker", "podman", "k8s"]
image = "redis:7-alpine"
command = ["redis-server", "--requirepass", "${REDIS_PASSWORD}"]

[[volumes]]
name = "redis-data"
mount_path = "/data"
description = "Persistent database storage"

[[network]]
internal_port = 6379
external_port = 6379
protocol = "tcp"

[healthcheck]
command = "redis-cli ping | grep PONG"
interval = "30s"
timeout = "5s"
retries = 3

[update]
strategy = "recreate"
```
