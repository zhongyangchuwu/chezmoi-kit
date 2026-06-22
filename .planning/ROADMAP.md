# Roadmap: cm

## Overview

`v0.2.0` planning is currently in discussion. The next scope is not committed yet: the active question is how `cm sync` should support chezmoi template-backed targets by showing destination, rendered target, and source template state, and by explaining chezmoi `edit`/`merge` behavior before any implementation plan is written.

## Version Policy

- Patch releases (`v0.1.Z`) may add optional commands or flags, improve diagnostics, fix display defects, and update non-runtime project hygiene when defaults and the user mental model stay stable.
- Minor releases (`v0.Y.0`) are reserved for behavior, safety-model, status-model, output-contract, or command-semantics changes that users must consciously absorb.
- Major release planning (`v1.0.0`) is deferred until the public command contracts and safety model are stable enough to commit to long-term compatibility.

## Phases

- [ ] Phase 1: Template Operation Model Discussion

## Phase Details

### Phase 1: Template Operation Model Discussion

- Goal: Decide the TUI-first operation model for template-backed targets before implementation scope is committed.
- Depends on: `v0.1.1` stable sync TUI, report rendering, and doctor diagnostics.
- Requirements: TEMPLATE-DISCUSS-01, TEMPLATE-DISCUSS-02, TEMPLATE-DISCUSS-03.
- Success Criteria:
  - The discussion distinguishes destination file, rendered target, source state, and source template content for template-backed targets.
  - The discussion records how `chezmoi merge` works for templates: destination, source, and target participate, and target-state failure can fall back to two-way merge.
  - The discussion decides whether template source should be displayed inline, in a detail/modal view, or through delegated external commands.
  - The discussion decides whether direct template mutation is limited to delegated `chezmoi edit`, `chezmoi merge`, and optional `chezmoi add --template` paths.
  - No implementation plan is written until unresolved TUI layout, keybinding, and CLI compatibility decisions are settled.
- Plans: TBD
  - [ ] pending-discussion-plan

## Progress

| Phase | Plans Complete | Status | Completed |
|---|---:|---|---|
| Phase 1: Template Operation Model Discussion | 0 | discussing | - |
