# Verification: Directory Browser

## Claims Checked

1. Direct-child Current rows align files and directories and no longer render `▾/▸` disclosure glyphs.
2. Missing ancestors become navigable virtual directories with aggregate state.
3. `h`/left and `l`/right/Enter navigate paths in Current focus and restore parent selections; preview focus retains horizontal scrolling.
4. A single explicit directory scope opens at that directory.
5. Filters retain matching ancestors; global search finds inventory nodes and Enter locates the selected result.
6. Directory previews are local summaries and do not call chezmoi content preview.
7. Wide, medium, narrow, help, semantic-color, and `NO_COLOR=1` layouts remain usable.
8. Workspace remains read-only.

## Automated Evidence

- `go test ./internal/tui` passed, including virtual ancestor, scope, key navigation, cursor restoration, filtering/search locating, local directory summary, aligned rows, and responsive panes.
- `go test ./...` passed.
- `go test -race ./...` passed.
- `go vet ./...`, `go mod tidy -diff`, `git diff --check`, and TUI LSP diagnostics passed.

## Manual Evidence

- Wide PTY: Parent/Current/Preview showed the explicit workspace scope; directories used trailing `/`, aligned `C:d`/`C:f` rows, local summary preview, `l` enter, `h` return, and `?` help.
- `NO_COLOR=1` PTY: textual state/type/attribute markers, legend, pane labels, and responsive footer remained readable.
- 50-column PTY: only Current rendered, no Parent or Preview appeared, and Tab remained the preview route.
- Source/destination hash after all exits matched the known baseline: `f4c92711ff0f5bf1e2e50969f52bd5d083854777d24697d841ec8a91f81bb804`.

## Result

passed
