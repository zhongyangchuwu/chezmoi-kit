# State

## Project Reference

See `.planning/PROJECT.md` (updated 2026-09-07).

**Core value:** Make personal chezmoi reconciliation explicit, reviewable, and low-surprise before any file mutation happens.
**Current focus:** Workspace usability is complete; prepare the separate Phase 3 tool-handoff execution plan when its product direction is approved.

## Current Position

- Phase: Phase 2.1 — Workspace Usability
- Plan: `phase-2.1-workspace-usability-01`
- Status: complete
- Active Artifact: `.planning/phases/phase-2.1-workspace-usability/VERIFICATION.md`
- Branch: `feat/tui-usability`
- Last activity: 2026-09-07 — added semantic workspace hierarchy, responsive help/feedback, no-color fallback, behavior coverage, PTY smoke, and local quality gates.
- Progress: Phases 1, 2, and 2.1 are complete; Phases 3/4 remain planned.

## Accumulated Context

### Product Decisions

- `cm ui [path...]` is an explicit persistent, read-only workbench in this phase.
- Bare `cm` remains status; `cm sync` remains focused mutation with Phase 1 safety behavior.
- One TUI model/rendering path serves both modes; mode-specific mutation and clean-entry behavior remain explicit.
- Inventory delegates to chezmoi managed/unmanaged/ignored/status and type filters. No ignore or source-name parser.
- Unmanaged discovery requires explicit scopes and recursively queries chezmoi for returned directories.
- Diff is default. Destination, target, and source views are labeled and bounded; known sensitive content needs explicit reveal.
- Tree/flat/filter/search/collapse are projections over one absolute-path identity set.
- Workspace color augments retained status/type/attribute text; `?` opens a keyboard-isolated help overlay and `NO_COLOR=1` preserves textual navigation.

### Blockers/Concerns

- No current blocker.
- Explicit reveal can expose rendered/decrypted content in terminal scrollback.
- Scoped unmanaged discovery is intentionally bounded at 10,000 entries and 2,000 queries; users must narrow large scopes.
- `ignored` reports entries but not ignore-rule provenance.

## Verification Contract

- Real fixture passed: clean, dirty, template, encrypted, symlink, directory, script, ignored, nested unmanaged, and secret-backed template reveal.
- Actual PTY passed: normal/narrow layouts, clean persistence, projections, full preview navigation, content views, reveal, and no-mutation exit.
- Full Go, race, vet, module, vulnerability, GoReleaser check/snapshot, LSP/YAML, and whitespace gates passed.

## Session Continuity

- Last session: 2026-09-07
- Completed: Phase 2 verified and captured.
- Next Action: create the Phase 3 editor/merge/lazygit handoff execution plan before implementation.
- Resume file: `.planning/ROADMAP.md`
