# Milestone 2: Template Engine

## Timeline
**Target:** Q4 2026 (~2 months)

## Goal
Implement the core Template Engine described in SPEC-001. This allows Host Anything to understand service descriptions, validate user input, and securely manage variables.

## Key Deliverables
1. **TOML Parser & Validator:** Go structs and parsers to read `template.toml` files and strictly validate them against SPEC-001.
2. **Local Registry:** File-system based registry to index and store templates in `/var/lib/hostanything/templates`.
3. **Variable Substitution Engine:** Logic to merge user config (SPEC-003) with template definitions, ensuring required fields are present and regex validations pass.
4. **Secret Handling Core:** Cryptographic utilities to encrypt `secret` variables before saving to disk.

## Success Criteria
- Engine successfully parses all sections of a valid complex TOML template.
- Engine correctly rejects invalid templates with precise, line-numbered error messages.
- Engine correctly flags user configuration that fails to meet the `validation_regex`.
- Secret variables are verified to be encrypted at rest in the data directory.

## Out of Scope
- Actually deploying the parsed configuration to Docker.
- Fetching templates from GitHub.
