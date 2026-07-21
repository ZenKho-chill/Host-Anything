# Web UI Architecture

## Overview
The Web UI is the primary control plane for Host Anything. It is built as a Single Page Application (SPA) using React, TypeScript, and Vite. It is designed to be lightweight, responsive, and ultimately embedded into the Go binary for seamless distribution.

## Tech Stack Setup
- **Framework**: React 18+
- **Build Tool**: Vite (chosen for fast cold starts and HMR)
- **Language**: TypeScript (strict mode enabled for robustness)
- **Routing**: React Router DOM
- **Styling**: Vanilla CSS (Premium aesthetics with Dark Mode, Glassmorphism, Micro-animations)
- **Data Fetching**: Native Fetch API

## Authentication Flow
Security is paramount since Host Anything controls server infrastructure.

1. **Login**: User submits credentials to `/api/v1/auth/login`.
2. **JWT**: On success, the Core API returns a JSON response containing a JWT signed with HS256.
3. **Session**: The UI stores the JWT in `localStorage` and sends it as a `Bearer` token in the `Authorization` header for subsequent API calls.
4. **Fail2ban Integration**: If the login fails, the Go backend logs the failure (including IP address) as structured JSON. The server administrator configures Fail2ban to monitor this log and block IPs with excessive failed attempts. The UI gracefully handles 401 Unauthorized errors by redirecting to the login page.

## Key Views

1. **Dashboard**: High-level overview. Shows system resource usage (CPU/Mem of host), total active services, recent alerts, and quick actions.
2. **Services List**: A tabular or grid view of all installed services. Displays status (Running, Stopped, Error), assigned ports, and uptime.
3. **Service Detail**: Deep dive into a specific service.
   - **Metrics Tab**: Real-time CPU/Mem usage of the container/process.
   - **Logs Tab**: Terminal-like view streaming logs via WebSocket.
   - **Config Tab**: Form to edit environment variables, ports, and volumes (triggers `ApplyConfig`).
4. **Template Browser (Marketplace)**: Connects to the GitHub-backed marketplace. Users can search, filter by category, and read template READMEs before installation.
5. **Settings**: Daemon-level configuration. Manage user accounts, default runtime preferences, backup schedules, and network settings.

## API Communication
- **REST**: Standard CRUD operations (fetching templates, updating configs, starting/stopping services) use standard RESTful HTTP calls.
- **Server-Sent Events (SSE)**: Real-time log streaming is implemented via HTTP chunked transfer and standard `text/event-stream`.
