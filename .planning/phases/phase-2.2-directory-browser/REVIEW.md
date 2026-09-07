# Review: Directory Browser

## Scope Reviewed

- Directory projection/state model, virtual nodes, search/filter transitions, preview loading, responsive rendering, keyboard routing, behavior tests, and public documentation.

## Findings and Fixes

1. **High — directory disclosure consumed a column only on directories.** File state/type/name columns therefore misaligned with sibling directories. Replaced recursive disclosure rows with fixed-column direct-child rows and trailing directory slash.
2. **High — recursive expansion mixed selection and expansion affordances.** `>` selection and `▾/▸` disclosure were adjacent and ambiguous. Removed expand/collapse interaction and gave Current focus explicit parent/child navigation.
3. **Medium — path ancestors may not exist as managed directory entries.** Added virtual navigation nodes derived from inventory, not filesystem reconstruction.
4. **Medium — filtered descendants could become unreachable.** Direct-child filters retain every ancestor having a matching descendant.
5. **Medium — directory preview must not act as a file-content request.** Added local matching-child summary commands and asserted zero WorkspaceService preview calls.
6. **Medium — previous preview scroll width assumed two panes.** Updated scroll width calculation for wide three-pane, medium two-pane, narrow one-pane, and full-screen layouts.

## Intentional Non-Features

- No files are opened, edited, selected, copied, deleted, or mutated.
- No tabs, mouse navigation, configurable layout/theme, recursive scanning, or external handoff was added.

## Remaining Risk

- The directory browser synthesizes nodes at render/projection time. The existing inventory caps remain the boundary; profile large real inventories before adding richer directory metadata.
