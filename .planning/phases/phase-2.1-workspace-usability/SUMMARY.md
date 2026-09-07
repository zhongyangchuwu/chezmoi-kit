# Summary: Workspace Usability

## Completed Changes

- Added a semantic Lip Gloss palette for workspace states, types, attributes, selection, panes, preview feedback, and unified diff lines.
- Added color-coded `C/D/U/I/R/?`, directory/symlink type treatment, and `[T]/[E]` badges while retaining textual meaning without color.
- Added compact file-pane legend, selected-row hierarchy, responsive contextual footer, and preview loading/error/withheld/clean feedback styles.
- Added a workspace-only `?` help overlay with quick start, files/preview key groups, state/type legend, and keyboard isolation.
- Added `NO_COLOR=1` style fallback and public workspace usage/design/development documentation.

## Files Changed

- Runtime: `internal/tui/styles.go`, `internal/tui/workspace_help.go`, and existing workspace rendering, update, model, key, and diff renderer files.
- Tests: `internal/tui/workspace_test.go`.
- User docs: `docs/usage.md`, `docs/design.md`, and `docs/development.md`.

## Deviations

- The compact help layout prioritizes quick start and text legends at narrow sizes rather than adding paging or persistent onboarding state.

## Evidence

- `go test ./...`, `go test -race ./...`, `go vet ./...`, `go mod tidy -diff`, TUI LSP diagnostics, and `git diff --check` passed.
- Actual PTY exercised colored workspace hierarchy, `?` help open/close isolation, and `NO_COLOR=1` text fallback.
- Isolated source/destination archive hash was unchanged before and after workspace exits.

## Unresolved Risks

- Explicit reveal still intentionally exposes rendered/decrypted content in terminal scrollback.
- Very narrow terminals necessarily truncate help and footer prose, but preserve textual markers and help/quit routes.
