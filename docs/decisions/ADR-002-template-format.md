# ADR 002: Template Format - TOML

## Status
Accepted

## Context
We need a human-readable, declarative format to define "Templates" (the blueprint for a service including required env vars, ports, and image details). The target audience includes both developers and hobbyist home-lab users.

Options considered:
- **JSON**: Too strict, lacks comments, difficult to read for large configurations.
- **YAML**: Extremely popular in the DevOps space (Kubernetes, Docker Compose). However, significant whitespace can cause notoriously difficult-to-debug errors for beginners. It also has complex edge cases (e.g., the "Norway problem" where `NO` evaluates to boolean false).
- **HCL (HashiCorp Configuration Language)**: Very powerful, but primarily associated with Terraform. Might be overkill and presents a learning curve.
- **TOML (Tom's Obvious, Minimal Language)**: Explicit, strongly typed, supports comments, and relies on clear block structures rather than indentation.

## Decision
We will use **TOML** as the standard format for Host Anything templates.

## Rationale
TOML strikes the best balance between readability and strictness. Its lack of reliance on significant whitespace makes it much friendlier for copy-pasting and manual editing by novice users. It maps cleanly to Go structs using standard libraries. It enforces a flat, readable structure that aligns well with our relatively simple template schema requirements.

## Consequences
- We avoid indentation-based bugs commonly seen in YAML.
- We must provide excellent documentation and examples, as users might be more accustomed to Docker Compose's YAML format.
- We will use `pelletier/go-toml/v2` in the Go backend for fast, reliable parsing.
