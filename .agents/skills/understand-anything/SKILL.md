---
name: understand-anything
description: Skill for systematically understanding the Host Anything codebase.
---

# Understand Anything Skill

Use this skill when you need to gain context about the Host Anything project before making modifications or implementing new features. Follow these steps systematically:

1. **Start with the Architecture**: Read `docs/architecture/00-overview.md` to understand the high-level system design.
2. **Review Decisions**: Read relevant Architectural Decision Records (ADRs) in `docs/decisions/` to understand why certain design choices were made.
3. **Read Specifications**: Review spec files in `docs/specs/` for the component you are working on to grasp functional and technical requirements.
4. **Trace Execution Flow**: When analyzing code, trace the path from the `cmd/` entry point → down to `internal/` services → and review shared types in `pkg/`. 
5. **Check Coding Rules**: Always review the `.agents/AGENTS.md` file to ensure you understand the project's strict coding rules, naming conventions, and layer boundaries before making changes.
6. **Understand Runtime Adapters**: When working with or understanding a runtime adapter (Docker, Podman, K8s, Host), read the `RuntimeAdapter` interface definition in `pkg/types/` first to see the contract it must fulfill.
7. **Reference Milestones**: Check milestone documents to understand what features are planned versus what is already implemented.
