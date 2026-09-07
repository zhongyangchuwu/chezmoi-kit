# Context: Phase 2 Unified File Workspace

## Goal

Add an explicit persistent `cm ui [path...]` workbench that can browse and fully inspect clean/dirty managed entries, scoped unmanaged candidates, and ignored entries while preserving the existing `cm` and `cm sync` contracts.

## Constraints

- Phase 1 authoritative reconciliation is the baseline; do not weaken its diff, fingerprint, postflight, script, encryption, or type-safety behavior.
- This phase is read-only in `cm ui`. Existing mutations remain in focused `cm sync` and direct commands until later phases add contextual management.
- Chezmoi owns managed/unmanaged/ignored classification, source mapping, target types, templates, encryption, target rendering, and ignore rules.
- Unmanaged discovery runs only for explicit path scopes. No default HOME traversal and no automatic add.
- Source/destination/target content is bounded. Encrypted source and template/encrypted rendered target content require a second explicit reveal action.
- Do not build an editor, merge tool, Git UI, filesystem watcher, ignore-rule interpreter, or source filename attribute parser.
- Reuse one TUI model and rendering path for focused sync and the general workbench; mode-specific behavior must stay explicit.

## Decisions

- Add `WorkspaceService`, `WorkspaceSnapshot`, `WorkspaceEntry`, and `WorkspacePreview` to the app boundary.
- Add `cm ui [path...]`; bare `cm` stays status and `cm sync` stays focused mutation.
- Build inventory from `chezmoi managed --path-style=all --format=json`, type-filtered managed queries, strict status, `ignored`, and scoped recursive `unmanaged` queries.
- Use an all-expanded tree with collapsible directory rows plus a flat view; both are projections of one sorted inventory.
- Keep selected identity by absolute target path across tree/filter/search rebuilds; fall back to the nearest previous index, then the first entry.
- Keep authoritative diff as preview 1. Preview 2/3/4 are destination, rendered target, and source state.
- `/` searches file paths when file list has focus and preview text when preview has focus. `n`/`N` traverse preview matches.
- `z` toggles full-screen preview; `h`/`l` scroll horizontally in preview focus; `[`/`]` jump diff hunks.

## Open Questions

- None blocking. More advanced tree loading, inline hunk editing, tool handoffs, and mutating file lifecycle actions remain later phases.

## Verification Expectations

- Unit tests cover inventory merging, scope normalization, projections, selection retention/fallback, collapse, filter/search, preview switching, horizontal scrolling, preview search, hunk jumps, full-screen layout, and sensitive reveal state.
- Real chezmoi integration covers clean/dirty files, templates, encrypted state, symlink, directory, script, ignored source entry, and nested unmanaged discovery under an explicit scope.
- Actual PTY smoke covers `cm ui` at normal and narrow sizes, clean-file persistence, filters/search, long-line tail, hunk jump, content views, and quitting without mutation.
- Existing `cm`, `cm status`, `cm diff`, and `cm sync` tests and real behavior remain green.
