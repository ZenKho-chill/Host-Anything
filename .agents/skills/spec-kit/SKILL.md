---
name: spec-kit
description: Skill for creating and managing technical specifications for Host Anything.
---

# Spec Kit Skill

Use this skill to create and maintain technical specifications for the project.

## How to create a new SPEC file
- Specs should be stored in `docs/specs/`
- Naming convention: `SPEC-XXX-description.md` (e.g., `SPEC-001-template-marketplace.md`)
- Always base new specs on the provided template.

## Required Sections
Every spec MUST contain the following sections:
1. **Overview**: High-level summary of the specification.
2. **Motivation**: Why this specification is needed.
3. **Scope**: What is included in the implementation.
4. **Out of Scope**: Explicitly state what will NOT be addressed.
5. **Specification**: The core technical details, workflows, and logic.
6. **Data Schemas**: Definitions of TOML/JSON structures, database schemas, or payload formats.
7. **Error Handling**: How edge cases and failures are managed.
8. **Security Considerations**: Authentication, authorization, input validation, and fail2ban impacts.
9. **Testing Strategy**: How this feature will be unit and integration tested.
10. **Open Questions**: Unresolved issues or decisions to be made.

## Versioning Specs
- Use semantic versioning at the top of the spec file or date-based revisions in a changelog section within the document.
- Once a spec is implemented, any changes must be reflected as updates to the version and changelog.

## Template Link
Always start a new spec by copying the template: [spec-template.md](templates/spec-template.md)
