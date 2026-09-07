# State

## Project Reference

See `.planning/PROJECT.md` (updated 2026-09-07).

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** Directory browsing is complete; prepare the separate Phase 3 tool-handoff execution plan when its product direction is approved.

## Current Position

- Phase: Phase 2.2 — Directory Browser
- Plan: `phase-2.2-directory-browser-01`
- Status: complete
- Active Artifact: `.planning/phases/phase-2.2-directory-browser/VERIFICATION.md`
- Branch: `main`
- Last activity: 2026-09-07 — squash-merged workspace usability and directory browsing through PR #5 as `6fed2ca`; merged-main CI passed.
- Progress: Phases 1, 2, 2.1, and 2.2 are complete and merged; Phases 3/4 remain planned.

## Accumulated Context

### Product Decisions

- `cm ui [path...]` is an explicit persistent, read-only workbench in this phase.
- Bare `cm` remains status; `cm sync` remains focused mutation with Phase 1 safety behavior.
- One TUI model/rendering path serves both modes; mode-specific mutation and clean-entry behavior remain explicit.
- Inventory delegates to chezmoi managed/unmanaged/ignored/status and type filters. No ignore or source-name parser.
- Unmanaged discovery requires explicit scopes and recursively queries chezmoi for returned directories.
- File previews remain authoritative and bounded; virtual directory ancestors exist only for local navigation and summary previews.
- Workspace uses a direct-child Parent/Current/Preview browser. Current focus uses `h`/`l` for parent/child navigation; Preview focus retains horizontal `h`/`l` scrolling.
- Workspace color augments retained status/type/attribute text; `?` opens a keyboard-isolated help overlay and `NO_COLOR=1` preserves textual navigation.

### Blockers/Concerns

- No current blocker.
- Explicit reveal can expose rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is intentionally bounded at 10,000 entries and 2,000 queries; users must narrow large scopes.
- `ignored` reports entries but not ignore-rule provenance.

## Verification Contract

- Real fixture passed: clean, dirty, template, encrypted, symlink, directory, script, ignored, nested unmanaged, and secret-backed template reveal.
- Actual PTY passed: Parent/Current/Preview navigation, search-to-preview loading, directory summaries, normal/narrow/no-color layouts, full preview navigation, reveal, and no-mutation exit.
- Full Go, race, vet, module, LSP, and whitespace gates passed locally; PR and merged-main GitHub CI passed.

## Session Continuity

- Last session: 2026-09-07
- Completed: Phases 2.1/2.2 reviewed, squash-merged through PR #5, and verified on merged `main` at `6fed2ca`.
- Next Action: create the Phase 3 editor/merge/lazygit handoff execution plan before implementation.
- Resume file: `.planning/ROADMAP.md`
