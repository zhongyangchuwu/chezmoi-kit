# State

## Project Reference

See `.planning/PROJECT.md` (updated 2026-09-10).

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** Phase 3.1 workspace source edit is implemented and locally verified; continue the remaining external-tool handoffs after its PR is merged.

## Current Position

- Phase: Phase 3.1 — Workspace Source Edit
- Plan: `phase-3.1-workspace-source-edit-01`
- Status: verified locally, pending PR
- Active Artifact: `.planning/phases/phase-3.1-workspace-source-edit/VERIFICATION.md`
- Branch: `feat/workspace-source-edit`
- Last activity: 2026-09-10 — implemented source-edit terminal handoff, authoritative refresh, stale-preview invalidation, documentation, and real-tool validation.
- Progress: Phases 1, 2, 2.1, and 2.2 are complete and merged; Phase 3.1 source edit is locally complete; remaining Phase 3 and Phase 4 work is planned.

## Accumulated Context

### Product Decisions

- `cm ui [path...]` is an explicit persistent, inspection-first workbench. Browsing and preview remain non-mutating; explicit `e` edits eligible chezmoi source only and forces no apply/watch.
- Bare `cm` remains status; `cm sync` remains focused mutation with Phase 1 safety behavior.
- One TUI model/rendering path serves both modes; mode-specific mutation and clean-entry behavior remain explicit.
- Inventory delegates to chezmoi managed/unmanaged/ignored/status and type filters. No ignore or source-name parser.
- Unmanaged discovery requires explicit scopes and recursively queries chezmoi for returned directories.
- File previews remain authoritative and bounded; virtual directory ancestors exist only for local navigation and summary previews.
- Workspace uses a direct-child Parent/Current/Preview browser. Current focus uses `h`/`l` for parent/child navigation; Preview focus retains horizontal `h`/`l` scrolling.
- Workspace color augments retained status/type/attribute text; `?` opens a keyboard-isolated help overlay and `NO_COLOR=1` preserves textual navigation.
- `e` permits managed regular files/symlinks with a source map, including templates and encrypted files. It runs `chezmoi edit --apply=false --watch=false <absolute-target>`, yields terminal ownership, re-inventories original scopes on every editor return, preserves viable navigation state, and invalidates previews/reveals/search/scroll state.

### Blockers/Concerns

- No current blocker.
- Explicit reveal can expose rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is intentionally bounded at 10,000 entries and 2,000 queries; users must narrow large scopes.
- `ignored` reports entries but not ignore-rule provenance.

## Verification Contract

- Real fixture passed: clean, dirty, template, encrypted, symlink, directory, script, ignored, nested unmanaged, and secret-backed template reveal.
- Source edit service integration passed for regular, template, and encrypted source; it confirmed source edits do not auto-apply destination. TUI coverage includes eligibility, handoff input locking, refresh/restoration, stale-preview response rejection, editor error, and refresh failure.
- Actual PTY passed: configured editor changed managed source through `e`, TUI returned/refreshed, and destination remained unchanged.
- Full Go, race, vet, module, LSP, and whitespace gates passed locally.

## Session Continuity

- Last session: 2026-09-10
- Completed: Phase 3.1 source edit implementation and local verification on `feat/workspace-source-edit`; branch is ready for review/PR.
- Next Action: create PR, obtain review and CI, then squash merge and capture delivery evidence.
- Resume file: `.planning/phases/phase-3.1-workspace-source-edit/VERIFICATION.md`
