# Requirements: cm

**Defined:** 2026-06-19
**Core Value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Current Requirements

### v0.1.1 — Safe Usability Polish

- [ ] OUT-FLAGS-01 — Users can explicitly select plain, ANSI, or Markdown report output for existing non-interactive report-backed commands without changing default output behavior.
- [ ] OUT-FLAGS-02 — Users can explicitly select color policy (`auto`, `always`, `never`) while preserving current TTY auto-detection and `NO_COLOR` defaults.
- [ ] DOCTOR-01 — `cm doctor` provides a read-only prerequisite and environment diagnostic for required and optional external tools without mutating local files or chezmoi source state.
- [ ] TUI-WIDTH-01 — `cm sync` TUI truncates file names, status text, and diff lines by display width rather than byte length, preserving valid UTF-8 and improving wide-character rendering.

## Deferred

- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and app sync services.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package-manager publishing.
- Issue templates, PR templates, and security policy files.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.

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
| OUT-FLAGS-01 | v0.1.1 Phase 1 — Output Controls | Planned |
| OUT-FLAGS-02 | v0.1.1 Phase 1 — Output Controls | Planned |
| DOCTOR-01 | v0.1.1 Phase 2 — Doctor Diagnostics | Planned |
| TUI-WIDTH-01 | v0.1.1 Phase 3 — TUI Display Width Polish | Planned |

## Coverage Summary

- Total current requirements: 4
- Mapped count: 4
- Unmapped count: 0
