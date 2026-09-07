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
7. **High — each directory node rescanned the full inventory during projection.** At the 10,000-entry cap this produced avoidable quadratic work and Parent rendering repeated it. Replaced it with one cached node/match index rebuilt in path-depth-linear work per projection.
8. **High — accepting a global file-search result could leave its preview unloaded.** Search selection changed without loading during typing, then Enter compared the result to itself. Enter now explicitly loads the located file or directory preview; regression coverage observes the preview service call.
9. **Medium — multiple explicit scopes started at destination root.** Startup now computes their nearest common directory while retaining file-scope parent behavior.
10. **Medium — directory summaries could be stale, undercount virtual children, or iterate a shared slice asynchronously.** Directory summaries now bypass persistent cache reuse, count direct projected children, compute before command dispatch, and use an accurate directory-summary label.
11. **Medium — actual unknown directory ancestors and filtered Parent context could become non-navigable or invisible.** Any inventory node with descendants becomes a directory with aggregate state; Parent always retains and centers the current directory.
12. **Medium — Copilot found raw directory-summary prefix matching.** Non-canonical relative paths could be indexed correctly but omitted from descendant counts. Summary matching now normalizes the directory and every candidate through the same `cleanWorkspaceRelative` path used by the browser; regression coverage includes `./`, duplicate separator, and trailing separator input.

## Intentional Non-Features

- No files are opened, edited, selected, copied, deleted, or mutated.
- No tabs, mouse navigation, configurable layout/theme, recursive scanning, or external handoff was added.

## Remaining Risk

- Node and filter indexes rebuild in $O(n \times d)$ work, where $n$ is the bounded inventory and $d$ is path depth. This removes the reviewed quadratic path; profile only if future inventory limits grow materially.

## Merge Assessment

- No Blocker, High, Medium, or Low findings remain open. The change is ready for full verification, PR CI, and squash merge.
