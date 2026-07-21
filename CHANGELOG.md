# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **M1 (Foundation):** Strict TOML configuration loader, JSON structured logging, and graceful REST API server.
- **M2 (Template Engine):** TOML template parsing (SPEC-001) and variable resolution with automatic AES-256-GCM encryption for secrets.
- **M3 (Docker Runtime):** Service state machine (SPEC-002) and Docker adapter (Deploy, Stop, Logs, Status) with automatic port mapping and volume binding.
- **M4 (Web UI):** Premium React+TS frontend using Vanilla CSS (Glassmorphism/Dark Mode).
- **M4 (Authentication):** Bearer JWT login system with Fail2Ban-compatible security logging.
- **M4 (API):** Service deployment and live log streaming (SSE) endpoints.

## [0.1.0] - 2026-07-21
### Added
- Project bootstrap: `.gitignore`, `.editorconfig`, `LICENSE`, `README.md`, `CHANGELOG.md`, `Makefile`.
