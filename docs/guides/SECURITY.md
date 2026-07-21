# Security Policy

## Supported Versions
We support the latest minor release of Host Anything.

## Reporting a Vulnerability
Please do **not** open a public issue. Email security@host-anything.dev with details.
We will acknowledge receipt within 48 hours and provide a timeline for a fix.

## Auth Security Model
- Web UI uses secure session cookies and JWTs.
- Passwords are bcrypt hashed.
- Fail2ban integration is enabled by default to block brute-force attempts.

## Network Security Recommendations
- Do not expose the Host Anything management port to the public internet without a reverse proxy / TLS.
- Use Docker networks for inter-service communication.
