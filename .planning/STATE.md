# State

## Project Reference

See `.planning/PROJECT.md` (updated 2026-10-23).

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** Phase 3.2 workspace destination edit and merge is locally verified; obtain review and CI before merging, then implement source Git/lazygit handoff.

## Current Position

- Phase: Phase 3.2 — Workspace Destination Edit and Merge
- Plan: `phase-3.2-workspace-destination-edit-merge-01`
- Status: verified locally, pending PR
- Active Artifact: `.planning/phases/phase-3.2-workspace-destination-edit-merge/VERIFICATION.md`
- Branch: `feat/workspace-editor-merge`
- Last activity: 2026-10-23 — implemented destination editor, three-way merge, generalized workspace handoff refresh, documentation, and real chezmoi integration coverage.
- Progress: Phases 1, 2, 2.1, 2.2, and 3.1 are complete and merged; Phase 3.2 is locally complete; source Git/lazygit and Phase 4 remain planned.

## Accumulated Context

### Product Decisions

- `cm ui [path...]` is an explicit persistent, inspection-first workbench. Browsing and preview remain non-mutating; explicit `e`, `E`, and `m` are the only current workspace handoffs.
- `e` permits managed files/symlinks with a source map, including templates and encrypted files. It runs `chezmoi edit --apply=false --watch=false <absolute-target>` and never auto-applies/watches destination changes.
- `E` permits only managed regular files with a source map. It invokes `VISUAL`, falling back to `EDITOR`, directly with the absolute destination path; it rejects symlinks and does not change source or apply chezmoi state.
- `m` permits only managed regular files with a source map, matching the sync merge matrix. It delegates destination/source/rendered-target reconciliation to configured `chezmoi merge`; merge may change source or destination but cm does not apply after return.
- Every launched handoff yields terminal ownership, re-inventories original scopes after success/cancellation/nonzero return, preserves viable workspace context, and invalidates previews/reveals/search/scroll state.
- Bare `cm` remains status; `cm sync` retains Phase 1 review/fingerprint/postflight safety behavior.

### Blockers/Concerns

- No current blocker.
- Explicit reveal can expose rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is intentionally bounded at 10,000 entries and 2,000 queries; users must narrow large scopes.
- Merge tools may write source or destination. Workspace documents the three roles before handoff and refreshes afterward; it does not claim synchronization.

## Verification Contract

- Real source edit integration verifies no destination auto-apply; real destination edit integration verifies source preservation; real configured merge integration verifies destination/source/target argument roles and no implicit apply.
- TUI coverage includes eligibility, contextual guidance, terminal handoff input lock, refresh/restoration, stale-preview response rejection, tool error refresh, and refresh failure.
- Full Go, race, vet, module, LSP, and whitespace gates passed locally.

## Session Continuity

- Last session: 2026-10-23
- Completed: Phase 3.2 implementation and local verification on `feat/workspace-editor-merge`.
- Next Action: create PR, obtain review and CI, squash merge, capture delivery evidence, then plan source Git/lazygit handoff.
- Resume file: `.planning/phases/phase-3.2-workspace-destination-edit-merge/VERIFICATION.md`
