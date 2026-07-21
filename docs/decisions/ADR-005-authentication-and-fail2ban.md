# ADR 005: Authentication and Fail2ban Integration

## Status
Accepted

## Context
Host Anything manages underlying server infrastructure. Unauthorized access could lead to complete server compromise. We need a robust authentication system. Since this is a self-hosted tool, relying on external OAuth providers (Google, GitHub) is unacceptable, as the tool must work fully air-gapped or during internet outages. Furthermore, self-hosted endpoints are frequently targeted by brute-force bots.

## Decision
1. **Authentication Mechanism**: Local username/password authentication issuing stateless JWTs (JSON Web Tokens) stored in HttpOnly, Secure cookies.
2. **Brute Force Protection**: Integration with the industry-standard `fail2ban` system.

## Rationale
- **JWT + HttpOnly Cookies**: This prevents XSS attacks from stealing the session token, while avoiding the database overhead of maintaining stateful server-side sessions.
- **Fail2ban**: Rather than writing custom rate-limiting and IP banning logic inside the Go daemon (which is complex and error-prone), we leverage the host OS ecosystem. Host Anything will explicitly log authentication failures (with timestamps, target users, and source IPs) to a designated flat file (`/var/log/hostanything/auth.log`). 

## Consequences
- **Setup Requirement**: For maximum security, the `.deb` package installation script will automatically drop a fail2ban jail configuration into `/etc/fail2ban/jail.d/hostanything.conf` to monitor our log file.
- **Security Considerations**: We must be incredibly careful to sanitize usernames in the log output to prevent log injection attacks (e.g., a user putting a newline character in their username to spoof log entries).
- **Statelessness**: Since JWTs are stateless, immediate session revocation (e.g., kicking a user out) requires a token blocklist in the database, slightly negating the pure stateless benefit, but necessary for admin control.
