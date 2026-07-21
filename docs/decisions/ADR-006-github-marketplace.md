# ADR 006: GitHub-Based Template Marketplace

## Status
Accepted

## Context
To make Host Anything accessible, we need a way for users to easily discover and install application templates (like Nextcloud, Plex, or a generic Node.js environment) without writing TOML files manually. Maintaining a centralized, proprietary database or registry API is costly and creates a single point of failure.

## Decision
We will use **GitHub** as the decentralized backend for the Template Marketplace.

## Rationale
- Templates will be stored in standard GitHub repositories.
- To publish a template, a user simply names their repository following a strict convention: `hostanything-template-<name>`.
- The Host Anything core daemon uses the public GitHub Search API to discover repositories matching this convention.
- Fetching a template involves downloading the raw `template.toml` from the main branch.

## Consequences

**Trust Model:**
Since anyone can publish a template, we must implement a trust model in the UI:
1. **Official Templates**: Hardcoded list of trusted authors/organizations (e.g., `hostanything-org`). Show with a "Verified" badge.
2. **Community Templates**: Displayed with warnings. Users must explicitly acknowledge they are running community-provided configurations.

**Pros:**
- Zero infrastructure cost for the marketplace.
- Leverages existing tools for version control, issue tracking, and collaboration (PRs) on templates.
- Highly resilient.

**Cons:**
- Subject to GitHub API rate limits. The Go daemon must aggressively cache marketplace search results locally.
- Requires internet access for discovery (though users can manually load TOML files into `/var/lib/hostanything/templates` for offline/air-gapped usage).
