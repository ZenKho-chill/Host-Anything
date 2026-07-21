# SPEC-040: GitHub Marketplace

## Status
Approved

## Overview
Host Anything uses GitHub as a decentralized, open marketplace for templates. Instead of hosting a proprietary registry server, the system leverages GitHub's API, topics, and repository structures to discover, fetch, and validate templates.

## Motivation
To bootstrap an ecosystem rapidly without the overhead of maintaining registry infrastructure. Developers are already familiar with GitHub, making it easy to create, fork, and share templates.

## Scope
- Repository naming and structure conventions.
- Template discovery mechanism using GitHub API.
- Trust levels (Official vs. Community).
- Template installation flow.

## Out of Scope
- Automated testing of community templates (this is left to community CI/CD).

## Specification

### 1. Repository Conventions
To be recognized as a Host Anything template, a repository MUST adhere to the following rules:

- **Naming**: `hostanything-template-{name}` (e.g., `hostanything-template-nginx`).
- **Topic**: Must be tagged with the GitHub topic `hostanything-template`.
- **Required Files**:
  - `template.toml` at the repository root (must comply with SPEC-001).
  - `README.md` containing usage instructions.
- **Optional Files**:
  - `screenshots/` directory containing PNG/JPG files for the UI to display.

### 2. Discovery & Search Flow
The web UI queries the Host Anything core API (`/marketplace/search`).
The core performs a GitHub Search API request:
`GET https://api.github.com/search/repositories?q=topic:hostanything-template+{query}`

Data extracted from the GitHub API response includes:
- Repo name, description, stars, author, update timestamp.

### 3. Trust Levels
To protect users from malicious templates, trust levels are introduced:
- **Official**: Repositories located under the `github.com/host-anything` organization. Displayed with a "Verified Shield" in the UI. Priority in search results.
- **Verified**: Repositories owned by recognized partners. Hardcoded whitelist in core.
- **Community**: Any other repository matching the conventions. Requires user confirmation ("Warning: Community Template") before installation.

### 4. Fetch and Install Flow
1. **Selection**: User clicks "Install" on a marketplace item in the UI.
2. **Download**: Core downloads the raw `template.toml` using GitHub Raw User Content API:
   `https://raw.githubusercontent.com/{owner}/{repo}/{branch}/template.toml`
3. **Validation**: Core parses and validates the TOML against the SPEC-001 schema.
4. **Local Storage**: If valid, the template is saved to `/var/lib/hostanything/templates/{owner}_{repo}_{version}.toml`.
5. **Ready**: The UI is notified, and the user is redirected to the service configuration deployment screen.

## Error Handling
- GitHub API Rate Limits: Core must handle HTTP 403 Rate Limit Exceeded gracefully, returning a `503 Service Unavailable` with details to the UI. Authenticated GitHub tokens can be added to system config to increase limits.
- Invalid Templates: If a downloaded template fails validation, the installation aborts, and the UI displays the specific TOML parsing errors.

## Security
- Templates do not contain executable code for the Host Anything core; they only describe configs.
- The `command` and `image` directives in templates are executed in isolated runtimes (Docker/Podman), limiting host exposure.
- Explicit warnings are shown before installing Community templates.

## Testing Strategy
- API mocking of GitHub Search and Raw Content APIs to test discovery and download logic.
- Integration test feeding a structurally invalid `template.toml` to ensure the install flow safely aborts.
