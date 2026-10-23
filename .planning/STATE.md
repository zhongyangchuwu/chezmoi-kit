# State

## Project Reference

See `.planning/PROJECT.md` (updated 2026-10-23).

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** Phase 3.3 workspace source Git handoff is locally verified; obtain review and CI before merging. Phase 4 management and diagnostics is next.

## Current Position

- Phase: Phase 3.3 — Workspace Source Git Handoff
- Plan: `phase-3.3-workspace-source-git-01`
- Status: verified locally, pending PR
- Active Artifact: `.planning/phases/phase-3.3-workspace-source-git/VERIFICATION.md`
- Branch: `feat/workspace-source-git`
- Last activity: 2026-10-23 — implemented lazygit workspace handoff, nonfatal source Git summary, documentation, and isolated integrations.
- Progress: Phases 1, 2, 2.1, 2.2, 3.1, and 3.2 are complete and merged; Phase 3.3 is locally complete; Phase 4 is planned.

## Accumulated Context

### Product Decisions

- `cm ui [path...]` is an inspection-first workbench. Browsing and preview remain non-mutating; explicit `e`, `E`, `m`, and `g` are terminal handoffs.
- `e` permits managed files/symlinks with a source map and runs `chezmoi edit --apply=false --watch=false <absolute-target>`.
- `E` permits only managed regular files with a source map, invokes `VISUAL` then `EDITOR` directly with the destination path, rejects symlinks, and does not modify source/apply chezmoi state.
- `m` permits only managed regular files with a source map and delegates destination/source/rendered-target reconciliation to configured chezmoi merge. It may change source or destination but never triggers apply.
- `g` opens lazygit in the configured source repository. Workspace inventory stores the source Git porcelain change count; Git status errors are visible but do not prevent file inventory.
- Every launched handoff yields terminal ownership, re-inventories original scopes after success/cancellation/nonzero return, preserves viable workspace context, and invalidates previews/reveals/search/scroll state.
- Bare `cm` remains status; `cm sync` retains Phase 1 review/fingerprint/postflight safety behavior.

### Blockers/Concerns

- No current blocker.
- Explicit reveal can expose rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is intentionally bounded at 10,000 entries and 2,000 queries; users must narrow large scopes.
- Merge tools and lazygit can change source membership or paths. Workspace refresh relies on chezmoi inventory authority and existing selection fallback; cm does not claim a synchronized result.

## Verification Contract

- Isolated integrations prove source/destination editor safety, configured merge argument roles, lazygit source working directory, and source Git summary inventory.
- TUI coverage includes all handoff eligibility/guidance, terminal ownership, source Git status, refresh/restoration, stale-preview response rejection, tool error refresh, and refresh failure.
- Full Go, race, vet, module, LSP, and whitespace gates passed locally.

## Session Continuity

- Last session: 2026-10-23
- Completed: Phase 3.3 implementation and local verification on `feat/workspace-source-git`.
- Next Action: create PR, obtain review and CI, squash merge, capture delivery evidence, then plan Phase 4 management and diagnostics.
- Resume file: `.planning/phases/phase-3.3-workspace-source-git/VERIFICATION.md`
