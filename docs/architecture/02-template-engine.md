# Template Engine Design

## What is a Template?
A Template in Host Anything is a declarative, TOML-based definition of a deployable service. It dictates the base image or binary, required configuration points, default port mappings, and necessary volumes. Templates act as a blueprint that the system uses to scaffold a service instance.

## Template Schema Overview
Templates use TOML for its excellent readability and strict typing. 

```toml
name = "Nextcloud"
version = "1.0.0"
author = "Host Anything Community"
description = "A safe home for all your data."
recommended_runtime = "docker"

[image]
repository = "nextcloud"
tag = "27.0.2"

[config]
  [config.env]
  MYSQL_PASSWORD = { type = "string", required = true, secret = true, description = "Database password" }
  NEXTCLOUD_ADMIN_USER = { type = "string", default = "admin", required = true }

[network]
  [[network.ports]]
  internal = 80
  default_external = 8080
  protocol = "tcp"

[storage]
  [[storage.volumes]]
  internal_path = "/var/www/html"
  name = "nextcloud_data"
  required = true

[resources]
  min_memory_mb = 512
  recommended_memory_mb = 1024
```

## Template Resolution Pipeline
The pipeline ensures that templates are safe, valid, and properly structured before deployment.

1. **Fetch**: Templates are fetched from a local directory (`/var/lib/hostanything/templates`) or dynamically downloaded from GitHub repositories matching the marketplace convention.
2. **Validate**: The engine parses the TOML and validates it against a strict internal Go struct schema using the `validator` package. It checks for required fields, valid data types, and sensible resource limits.
3. **Render**: The engine merges the static template definition with user-provided configurations (from the database).
4. **Apply**: The resolved structural representation is passed to the Configuration Manager.

## Variable Substitution System
Templates often need dynamic data generated at deployment time or provided by the user. The engine supports a lightweight variable substitution syntax `{{ .VarName }}`.

- **System Variables**: `{{ .System.HostIP }}`, `{{ .System.ServiceName }}`
- **Generated Variables**: `{{ generate_password(16) }}` for auto-filling required secrets on first deploy.
- **User Inputs**: Values explicitly provided by the user in the UI, mapped to the `[config.env]` declarations.

## Applying Configuration to a Template
The core daemon translates the rendered template into a generalized `ServiceSpec` struct. 

- **Environment Variables**: Extracted directly from user input + generated secrets, passed to the adapter.
- **Ports**: Merges template `default_external` with user overrides. Checks host for port conflicts before proceeding.
- **Volumes**: Maps user-defined host paths or named volumes to the template's `internal_path`.
- **Resource Limits**: Applies user-defined CPU shares and Memory limits, ensuring they fall within the template's `min` and `max` boundaries (if defined).

This generalized `ServiceSpec` is what the Runtime Adapters consume, completely separating the TOML parsing logic from the deployment logic.
