# Development Setup Guide

## Prerequisites
- Go 1.22+
- Node.js 20+
- Docker/Podman
- Make

## Clone and Setup

```bash
git clone https://github.com/host-anything/hostanything.git
cd hostanything
```

## Running the App

To run both backend and frontend:
```bash
./scripts/dev.sh
```

## Testing and Linting
```bash
./scripts/lint.sh
go test ./...
```

## Building Packages Locally
To build a `.deb` file:
```bash
./scripts/build.sh
./scripts/package-deb.sh amd64
```

## Editor Setup
- **VS Code:** Install the Go and ESLint extensions.
- **GoLand:** Ensure you select Go 1.22 in settings and enable ESLint integration.
