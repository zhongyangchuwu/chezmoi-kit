# Context: Directory Browser

## Goal

Replace the recursive workspace tree with a Yazi-inspired parent/current/preview browser that eliminates mixed disclosure alignment and makes directory navigation explicit.

## Constraints

- Preserve the read-only workspace boundary, chezmoi inventory authority, preview/reveal policy, and sync behavior.
- Do not add file operations, selection mode, tabs, mouse support, filesystem scanning beyond the existing inventory, or configurable layouts.
- Preserve textual state/type/attribute semantics and `NO_COLOR=1` fallback.
- Remove tree/flat/collapse behavior cleanly; do not retain parallel navigation models.

## Decisions

- Wide terminals render Parent, Current, and Preview panes. Medium terminals hide Parent; narrow terminals use Current/Preview focus switching.
- The Current pane shows direct children of one `currentDir`; directories use a trailing `/` and no disclosure glyph.
- `h`/left returns to parent only in file focus. `l`/right/enter enters a selected directory only in file focus. Preview focus retains horizontal `h/l` scroll.
- Virtual directories derived from inventory paths serve navigation only and use aggregate child state; actual managed directory entries retain their metadata.
- File-path search is global. Accepting a result locates the parent directory and selects it; Escape restores the previous directory/selection.
- State filters keep a directory visible when it contains matching descendants.

## Open Questions

- None blocking.

## Verification Expectations

- TUI tests cover virtual directory creation, initial scope, entering/leaving directories, cursor restoration, filtering reachability, global search locate/restore, directory summary, responsive pane visibility, and no disclosure alignment.
- Actual PTY smoke covers parent/current/preview, h/l navigation, directories with and without actual managed directory entries, search, filtering, narrow fallback, help, and read-only exit.
