# Milestone 4: Web UI

**STATUS: COMPLETED**

## Timeline
**Target:** Q2 2027 (~3 months)

## Goal
Provide a responsive, modern frontend dashboard for users to interact with Host Anything. Secure the platform with robust authentication and system-level brute force protection.

## Key Deliverables
1. **React+TS Application:** Setup Vite/React project within the repository.
2. **Authentication Flow:** Implement SPEC-030 (JWT login, session management).
3. **Fail2ban Integration:** Output auth logs in the correct format and provide install scripts for fail2ban jails.
4. **Service Dashboard:** View running services, statuses, resource usage, and live logs.
5. **Dynamic Configuration Forms:** Auto-generate HTML forms based on `template.toml` schemas (handling string, int, enum, and secret fields).
6. **Local Template Browser:** UI to select and deploy from locally installed templates.

## Success Criteria
- A user can log in securely via the web interface.
- 5 successive failed logins result in an OS-level ban via fail2ban (verified in E2E tests).
- A user can deploy a new service using a graphical form.
- Passwords entered in dynamic forms are masked and stored securely.
- Live streaming of container logs to the UI via WebSockets or Server-Sent Events (SSE).

## Out of Scope
- GitHub Marketplace search integration.
- Complex data visualizations (simple text/table metrics are sufficient).
