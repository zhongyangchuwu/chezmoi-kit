# Context: Phase 1 Template Operation Model Discussion

## Goal
Decide how `cm` should support chezmoi template-backed targets in the `cm sync` TUI before committing to implementation details.

## Constraints
- Do not implement a second chezmoi template engine in `cm`.
- Keep chezmoi authoritative for source naming, template evaluation, data loading, ignore rules, encryption handling, and mutation semantics.
- Keep `cm sync` as the main user workflow; CLI commands should exist for compatibility and scripts, not as a competing primary UX.
- Do not change planning from unresolved questions into committed requirements until decisions are made.
- Keep automatic git commit/push/pull, daemon/watch behavior, and persistent state databases out of scope.

## Decisions
- `cm` should support template-backed targets by making destination, rendered target, and source template states visible together before mutation.
- For rendered target content, delegate to chezmoi (`chezmoi cat`/current diff path); current `cm diff` already compares rendered target output against destination content.
- For source template editing, prefer `chezmoi edit <target>` over direct writes to `.tmpl` files so chezmoi handles source path rules, encrypted files, hardlink/editor behavior, and template syntax checks.
- For merging, treat `chezmoi merge <target>` as a terminal-bound chezmoi operation whose semantics must be explained in the TUI before launch.

## Open Questions
- What TUI layout best shows rendered diff and source template without overloading the current two-pane review flow?
- Should TUI expose source template text inline, in a modal/detail view, or by temporarily launching `chezmoi edit`/external viewer?
- Which key should select “add as template” without confusing it with existing `a` add?
- Should `cm sync` show source git diffs for template files, or only source path/status metadata?
- Should `cm template data`/`cm template exec` be added in the same phase, or delayed until after TUI template review is designed?
- Should `chezmoi edit --apply` ever be exposed, or should `cm` require returning to review before apply?

## Verification Expectations
- Discussion output must separate decided boundaries from open UX choices.
- Any later plan must name the exact states shown, actions offered, and mutation path for each action.
- Implementation should not start until TUI layout, keybindings, and CLI compatibility scope are approved.
