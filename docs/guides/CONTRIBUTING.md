# Contributing to Host Anything

Welcome! We are excited to have you contribute.

## Dev Environment Setup

1. Install Go 1.22+
2. Install Node 20+
3. Install Docker
4. Install golangci-lint
5. Use `make` (or run scripts in `scripts/`)

## Running Locally

Run `make dev` or `bash scripts/dev.sh` to start the backend and frontend dev servers.

## Coding Rules

See AGENTS.md for a full overview. We enforce:
- 1 function or 1 domain group per file
- No global mutable state
- Wrap errors with context

## Pull Requests

1. Fork the repo
2. Create a branch
3. Use Conventional Commits
4. Open PR

## Writing Templates

Check out `TEMPLATE-AUTHORING.md` for a step-by-step guide on creating new templates and submitting them to the marketplace.
