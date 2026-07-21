---
name: code-review-graph
description: Skill for reviewing code with a focus on architectural dependencies and graph analysis.
---

# Code Review Graph Skill

Use this skill when reviewing pull requests or evaluating existing code in the Host Anything project to ensure strict architectural compliance.

## Architectural Laws (Must NEVER be violated)
1. **No Import Cycles**: Packages must form a directed acyclic graph (DAG).
2. **Strict Layer Boundaries**:
   - `cmd/` importing `internal/` directly is **OK**.
   - `internal/` importing `cmd/` is **NEVER OK**.
   - `pkg/` importing `internal/` or `cmd/` is **NEVER OK** (pkg must be independent).
3. **No Global State**: Variables declared at the package level cannot be mutable. Dependencies must be injected.

## Verification Tasks

### 1. Adapter Pattern Compliance
- Check that all runtime engines (Docker, Podman, K8s, Host) implement the exact `RuntimeAdapter` interface defined in `pkg/types/`.
- Ensure no runtime-specific logic leaks into the core engine orchestration.

### 2. Error Wrapping Compliance
- Ensure every `return err` that bubbles up is wrapped contextually: `fmt.Errorf("package.FunctionName: %w", err)`.

### 3. Missing Tests
- Check that every new business logic file in `internal/` has a corresponding `_test.go` file in the SAME package.
- Ensure all public functions have basic test coverage.

## Review Checklist Format
Include this checklist in your review output:

- [ ] Dependency rules obeyed (`pkg` independent, `internal` independent of `cmd`).
- [ ] No global mutable state present.
- [ ] One function or domain group per file strictly adhered to.
- [ ] Function variants split into separate files.
- [ ] All returned errors are wrapped with context.
- [ ] `RuntimeAdapter` interface fully implemented for any new runtime logic.
- [ ] Corresponding tests exist in the same package.
- [ ] Exported symbols have godoc comments.
- [ ] No magic numbers; constants used appropriately.
