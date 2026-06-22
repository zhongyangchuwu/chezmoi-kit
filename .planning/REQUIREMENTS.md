# Requirements: cm

**Defined:** 2026-06-19
**Core Value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Current Requirements

### v0.2.0 — Template Operation Model Discussion

- [ ] TEMPLATE-DISCUSS-01 — Decide how `cm sync` should show destination file content, rendered chezmoi target content, and source template content for template-backed targets before implementation scope is committed.
- [ ] TEMPLATE-DISCUSS-02 — Decide how `cm sync` should explain and launch chezmoi-backed source template operations such as `edit`, `merge`, and optional `add --template` without directly writing templates.
- [ ] TEMPLATE-DISCUSS-03 — Decide which CLI commands are compatibility/script surfaces for the TUI model, without making CLI the primary template workflow.

## Deferred

- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and app sync services.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.
- Issue templates, PR templates, and security policy files.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.
- Greedy `chezmoi add --autotemplate` exposure until `cm` has an explicit generated-template review workflow.

## Out of Scope

| Feature | Reason |
|---|---|
| Replacing chezmoi | `cm` delegates to chezmoi and avoids owning source manipulation directly. |
| Automatic git commit/push/pull | Current value is review and explicit reconciliation, not source repository automation. |
| Daemon/watch mode | Background mutation increases surprise and is outside current scope. |
| Persistent state database | Chezmoi state remains authoritative. |

## Traceability

| Requirement | Phase | Status |
|---|---|---|
| TEMPLATE-DISCUSS-01 | v0.2.0 Phase 1 — Template Operation Model Discussion | Discussing |
| TEMPLATE-DISCUSS-02 | v0.2.0 Phase 1 — Template Operation Model Discussion | Discussing |
| TEMPLATE-DISCUSS-03 | v0.2.0 Phase 1 — Template Operation Model Discussion | Discussing |

## Coverage Summary

- Total current requirements: 3
- Mapped count: 3
- Unmapped count: 0
