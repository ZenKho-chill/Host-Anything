# Milestone 6: Template Marketplace

## Timeline
**Target:** Q4 2027 - Q2 2028

## Goal
Implement SPEC-040 to turn Host Anything into a vibrant ecosystem. Enable discovery, installation, and updates of community and official templates directly from GitHub.

## Key Deliverables
1. **GitHub Search Integration:** Core backend endpoints to query the GitHub API for `hostanything-template` topics.
2. **Marketplace UI:** A new section in the dashboard to browse, search, and view READMEs of remote templates.
3. **Install Flow:** Safely pull `template.toml`, validate it, and register it locally.
4. **Official Template Library:** Create and publish at least 10 high-quality official templates (e.g., PostgreSQL, Redis, Nginx, Node.js, Python apps) under the `host-anything` GitHub org.
5. **Update Management:** Notify users in the UI when a deployed template has a new version on GitHub.

## Success Criteria
- Users can search for a community template directly in the web UI.
- Clicking 'Install' successfully downloads and validates the TOML from raw.githubusercontent.com.
- Trust badges (Official, Community) are prominently displayed to prevent accidental malicious installs.
- Official library covers common web stack infrastructure.

## Out of Scope
- Hosting a standalone proprietary registry server.
- Paid/Monetized templates (marketplace remains open and free via GitHub).
