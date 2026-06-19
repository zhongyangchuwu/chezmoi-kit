# Requirements: cm

**Defined:** 2026-06-19
**Core Value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.

## Current Requirements

### Runtime Safety

- [x] **SAFE-01**: `cm sync` must cleanly preserve terminal state when interrupted by normal TUI quit or Ctrl+C.
- [x] **SAFE-02**: `cm sync` must not imply that an already-confirmed execution can be quit or cancelled unless subprocess cancellation is actually implemented.
- [x] **SAFE-03**: Source repository git status parsing must surface malformed non-empty output instead of silently hiding it.

### Command Surface

- [x] **CLI-01**: Every exposed command intended for v0.1.0 must have command-level test coverage for wiring to the service boundary.
- [x] **CLI-02**: `cm edit <target>` must remain wired to `chezmoi edit` and shell completion must use managed file names.

### Release Hygiene

- [x] **REL-01**: Go module files must be tidy and reproducible under the release checks.
- [x] **REL-02**: Repository must include release metadata required by public users: license, changelog, ignore rules, and versioned build commands.
- [x] **REL-03**: GitHub CI must enforce the local release gate on pushes and pull requests.
- [x] **REL-04**: Tag-based GitHub release workflow must build versioned artifacts for v0.1.0.

### Documentation

- [x] **DOC-01**: Public documentation must describe current requirements, installation, usage, commands, and release workflow accurately.
- [x] **DOC-02**: Stale historical planning artifacts under `docs/superpowers/` must be removed from release documentation.
- [x] **DOC-03**: Current architecture documentation must match the implemented package boundaries and behavior.

### Verification

- [ ] **VER-01**: Release readiness must be proven by observed local checks before tagging.
- [ ] **VER-02**: Manual smoke scenarios for status, diff, sync, edit, git, completion, and version must be documented or completed before v0.1.0.

## Deferred

- `cm doctor` prerequisite checker.
- Context-aware subprocess cancellation across `process.Runner`, `chezmoi.Client`, and `reconcile.ReviewService`.
- GoReleaser signing, notarization, Homebrew, Scoop, Winget, Docker, and package manager publishing.
- Issue templates, PR templates, and security policy files.
- Richer status interpretation beyond the current simplified `! differs from chezmoi` model.

## Out of Scope

| Feature | Reason |
|---|---|
| Replacing chezmoi | `cm` delegates to chezmoi and avoids owning source manipulation directly. |
| Automatic git commit/push/pull | Current value is review and explicit reconciliation, not source repository automation. |
| Daemon/watch mode | Background mutation increases surprise and is outside v0.1.0. |
| Persistent state database | Chezmoi state remains authoritative. |
| Full subprocess cancellation in v0.1.0 | Requires broader API changes; conservative executing semantics are safer for first release. |

## Traceability

| Requirement | Phase | Status |
|---|---|---|
| SAFE-01 | Phase 1 | Complete |
| SAFE-02 | Phase 1 | Complete |
| SAFE-03 | Phase 1 | Complete |
| CLI-01 | Phase 1 | Complete |
| CLI-02 | Phase 1 | Complete |
| REL-01 | Phase 1 | Complete |
| REL-02 | Phase 2 | Complete |
| DOC-01 | Phase 2 | Complete |
| DOC-02 | Phase 2 | Complete |
| DOC-03 | Phase 2 | Complete |
| REL-03 | Phase 3 | Complete |
| REL-04 | Phase 3 | Complete |
| VER-01 | Phase 4 | Pending |
| VER-02 | Phase 4 | Pending |

## Coverage Summary

- Total current requirements: 14
- Mapped count: 14
- Unmapped count: 0
