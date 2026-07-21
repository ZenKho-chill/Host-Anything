# SPEC-030: Authentication

## Status
Approved

## Overview
This specification covers how authentication is handled in the Host Anything platform. It details JWT token management, session handling, rate limiting, and seamless integration with fail2ban to protect the host against brute-force attacks.

## Motivation
Host Anything controls fundamental system components and application deployments. Unauthorized access would lead to full system compromise. Standardizing auth and integrating it natively with OS-level protections (fail2ban) provides robust security.

## Scope
- JWT Authentication flow.
- Login and Refresh endpoints.
- Rate limiting at the application level.
- Fail2ban integration (logs, filters, jails).

## Out of Scope
- OAuth2 / SSO integration (future milestone).
- Role-Based Access Control (RBAC) - currently single admin user assumption.

## Specification

### 1. API Endpoints

#### `POST /api/v1/auth/login`
- **Request Body**: `{"username": "admin", "password": "..."}`
- **Success Response (`200 OK`)**:
  ```json
  {
    "token": "eyJhbG...",
    "expires_in": 3600,
    "refresh_token": "..."
  }
  ```
- **Failure Response (`401 Unauthorized`)**: Invalid credentials.

#### `POST /api/v1/auth/refresh`
- Swaps a valid refresh token for a new access token.

#### `POST /api/v1/auth/logout`
- Invalidates the current refresh token and clears the session.

### 2. JWT Configuration
- **Algorithm**: HS256
- **Secret**: Auto-generated high-entropy string stored in `/etc/hostanything/hostanything.toml`.
- **Access Token TTL**: 1 hour.
- **Refresh Token TTL**: 7 days.

### 3. Rate Limiting
- The `/api/v1/auth/login` endpoint is limited to 5 requests per minute per IP address internally via Go middleware.

### 4. Fail2ban Integration

To provide OS-level IP banning, Host Anything writes specific authentication events to a dedicated log file.

#### Log Format
Host Anything will append to `/var/log/hostanything/auth.log` in the following format:
`YYYY-MM-DD HH:MM:SS [AUTH] FAILED LOGIN from <IP_ADDRESS> for user <USERNAME>`

#### Fail2ban Filter (`/etc/fail2ban/filter.d/hostanything.conf`)
```ini
[Definition]
failregex = ^.* \[AUTH\] FAILED LOGIN from <HOST> for user .*$
ignoreregex =
```

#### Fail2ban Jail (`/etc/fail2ban/jail.d/hostanything.conf`)
```ini
[hostanything]
enabled = true
port = http,https,8080
filter = hostanything
logpath = /var/log/hostanything/auth.log
maxretry = 5
findtime = 600
bantime = 3600
```

## Error Handling
- Rate limit violations return `429 Too Many Requests`.
- Token expiration returns `401 Unauthorized` with a specific `TOKEN_EXPIRED` error code to signal the frontend to attempt a refresh.

## Security
- Passwords are hashed using bcrypt with cost=12.
- Session invalidation relies on a server-side deny-list for refresh tokens.
- Fail2ban mitigates distributed brute-force attacks at the firewall level (iptables/nftables).

## Testing Strategy
- Unit tests verifying JWT signing and verification.
- E2E tests simulating failed logins and verifying exact log output format matches the fail2ban regex.
- Integration tests ensuring rate limiter blocks excessive requests independently of fail2ban.
