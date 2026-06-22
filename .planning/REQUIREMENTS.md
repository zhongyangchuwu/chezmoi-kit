# Requirements: cm

**Defined:** 2026-06-19
**Core Value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Current Requirements

None. The completed v0.1.0 requirements are archived in `.planning/archive/releases/v0.1.0/SUMMARY.md`.

## Deferred

- `cm doctor` prerequisite checker.
- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and app sync services.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.
- Issue templates, PR templates, and security policy files.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.
- User-facing `--output` and `--color` flags for internal report renderers.
- TUI file/diff layout display-width handling.

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
| None | - | - |

## Coverage Summary

- Total current requirements: 0
- Mapped count: 0
- Unmapped count: 0
