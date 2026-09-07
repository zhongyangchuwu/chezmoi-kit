# Verification: Directory Browser

## Claims Checked

1. Direct-child Current rows align files and directories and no longer render `▾/▸` disclosure glyphs.
2. Missing or unknown ancestors become navigable directories with aggregate state.
3. `h`/left and `l`/right/Enter navigate paths in Current focus and restore parent selections; Preview focus retains horizontal scrolling.
4. One explicit directory scope opens there; multiple scopes open at their nearest common directory.
5. Filters retain matching ancestors; global search finds inventory nodes, and Enter locates and loads the selected result.
6. Directory previews are fresh local summaries, include virtual direct children, normalize descendant paths consistently with the browser index, and do not call chezmoi content preview.
7. Projection indexes avoid a full inventory scan for every directory node.
8. Wide, medium, narrow, help, semantic-color, and `NO_COLOR=1` layouts remain usable.
9. Workspace remains read-only.

## Automated Evidence

- `go test ./internal/tui` passed virtual/unknown ancestors, single/multiple scopes, key navigation, cursor restoration, filtered Parent context, search locating and preview loading, fresh local summaries, aligned rows, and responsive panes.
- `go test ./...` passed.
- `go test -race ./...` passed.
- `go vet ./...`, `go mod tidy -diff`, `git diff --check`, and TUI LSP diagnostics passed.
- Review inspection confirmed one cached node/match index rebuilt in bounded $O(n \times d)$ path-depth work rather than the prior per-node full-inventory scan.

## Manual Evidence

- Wide PTY: Parent/Current/Preview showed the explicit workspace scope; directories used trailing `/`, aligned `C:d`/`C:f` rows, local summaries, `l` enter, `h` return, and `?` help.
- Post-review PTY: `/dirty` plus Enter located `dirty.conf` and loaded its authoritative diff; returning with `h` showed `view: directory summary`, aggregate `D:d workspace/`, 10 direct entries, and 12 matching descendants.
- `NO_COLOR=1` PTY: textual state/type/attribute markers, legend, pane labels, and responsive footer remained readable.
- 50-column PTY: only Current rendered, no Parent or Preview appeared, and Tab remained the preview route.
- Source/destination hash after all exits matched the known baseline: `f4c92711ff0f5bf1e2e50969f52bd5d083854777d24697d841ec8a91f81bb804`.

## Remote Evidence

- PR #5 was squash-merged: https://github.com/zhongyangchuwu/chezmoi-kit/pull/5
- Merged commit: `6fed2ca feat(tui): improve workspace navigation and guidance`
- Merged-main CI passed: https://github.com/zhongyangchuwu/chezmoi-kit/actions/runs/34134102852

## Result

passed
